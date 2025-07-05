package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert user into MongoDB
	collection := client.Database("bowott").Collection("users")
	res, err := collection.InsertOne(r.Context(), user)
	if err != nil {
		log.Printf("Mongo InsertOne error: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User created with ID: %v", res.InsertedID)
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var foundUser User
	collection := client.Database("bowott").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&foundUser)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: foundUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("userId").(primitive.ObjectID)

	job := UploadJob{
		UserID:   userID,
		FileName: handler.Filename,
		Content:  buf.String(),
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Failed to create upload job", http.StatusInternalServerError)
		return
	}

	if rdb == nil {
		http.Error(w, "Redis is not available", http.StatusInternalServerError)
		return
	}

	if err := rdb.LPush(context.Background(), "upload-queue", jobJSON).Err(); err != nil {
		http.Error(w, "Failed to queue upload", http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("Failed to get queue count: %v", err)
		http.Error(w, "Failed to get queue count", http.StatusInternalServerError)
		return
	}

	print("Added file to queue")

	fmt.Fprintln(w, "File upload queued successfully")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	limit := int64(100)
	userID := r.Context().Value("userId").(primitive.ObjectID)

	collection := client.Database("bowott").Collection("documents")

	cursor, err := collection.Find(context.TODO(), bson.M{"userId": userID}, &options.FindOptions{Limit: &limit})
	if err != nil {
		http.Error(w, "Failed to retrieve documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var documents []UserDocument
	if err = cursor.All(context.TODO(), &documents); err != nil {
		http.Error(w, "Failed to decode documents", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(documents)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	print(query)
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	fileID := r.URL.Query().Get("fileId")
	if fileID == "" {
		http.Error(w, "Missing fileId", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, "Invalid fileId", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userId").(primitive.ObjectID)

	collection := client.Database("bowott").Collection("documents")
	var doc UserDocument
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "userId": userID}).Decode(&doc)
	if err != nil {
		http.Error(w, "Failed to retrieve document or unauthorized", http.StatusNotFound)
		return
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(doc.Content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var results []SearchResult
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			startIndex := strings.Index(strings.ToLower(line), strings.ToLower(query))
			print("startIndex %d", startIndex)
			endIndex := startIndex + len(query)
			results = append(results, SearchResult{
				Line:       line,
				LineNumber: i,
				StartIndex: startIndex,
				EndIndex:   endIndex,
			})
		}
	}

	json.NewEncoder(w).Encode(results)
}

func fileContentHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("fileId")
	if fileID == "" {
		http.Error(w, "Missing fileId", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, "Invalid fileId", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userId").(primitive.ObjectID)

	collection := client.Database("bowott").Collection("documents")
	var doc UserDocument
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "userId": userID}).Decode(&doc)
	if err != nil {
		http.Error(w, "Failed to retrieve document or unauthorized", http.StatusNotFound)
		return
	}

	fmt.Fprint(w, doc.Content)
}

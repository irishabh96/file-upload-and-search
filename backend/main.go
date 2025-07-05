package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func logRequest(r *http.Request) {
	log.Println("Request received:")
	log.Printf("Method: %s, URL: %s, Query: %s", r.Method, r.URL.Path, r.URL.RawQuery)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return
	}
	log.Printf("Body: %s", body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
}

func main() {
	initDB()
	initQueue()
	go worker()

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		registerHandler(w, r)
	})
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		loginHandler(w, r)
	})
	http.HandleFunc("/api/upload", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		uploadHandler(w, r)
	}))
	http.HandleFunc("/api/search", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		searchHandler(w, r)
	}))
	http.HandleFunc("/api/files", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		filesHandler(w, r)
	}))
	http.HandleFunc("/api/file", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fileContentHandler(w, r)
	}))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func worker() {
	if rdb == nil {
		return
	}
	for {
		result, err := rdb.BRPop(context.Background(), 0, "upload-queue").Result()
		if err != nil {
			fmt.Println("Error processing queue:", err)
			continue
		}

		var job UploadJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			fmt.Println("Error unmarshalling job:", err)
			continue
		}

		doc := UserDocument{
			UserID:  job.UserID,
			Name:    job.FileName,
			Content: job.Content,
		}

		collection := client.Database("bowott").Collection("documents")
		_, err = collection.InsertOne(context.TODO(), doc)
		if err != nil {
			fmt.Println("Error saving document from queue:", err)
		}
	}
}

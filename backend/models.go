package main

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type UserDocument struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"userId"`
	Name    string             `bson:"name"`
	Content string             `bson:"content"`
}

type Claims struct {
	UserID primitive.ObjectID `json:"userId"`
	jwt.StandardClaims
}

type UploadJob struct {
	UserID   primitive.ObjectID `json:"userId"`
	FileName string             `json:"fileName"`
	Content  string             `json:"content"`
}

type SearchResult struct {
	Line       string `json:"line"`
	LineNumber int    `json:"lineNumber"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
}

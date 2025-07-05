package main

import "os"

var (
	RedisURI = getEnv("REDIS_URI", "localhost:6379")
	MongoURL = getEnv("MONGO_URL", "mongodb://localhost:27017")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

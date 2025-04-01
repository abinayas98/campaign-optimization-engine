package data_adapters

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Redis address
	})
	// Ping to test if Redis is available
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Redis connection successful")
}

// Cache CPC and CVR Predictions for a Platform
func CachePrediction(platform string, predictedCPC, predictedCVR float64) {
	cacheKey := fmt.Sprintf("prediction:%s", platform)

	// Store both CPC and CVR as a single value (you can serialize it as JSON if needed)
	predictionValue := fmt.Sprintf("%f|%f", predictedCPC, predictedCVR)

	// Cache the value with a TTL (e.g., 1 hour)
	err := redisClient.Set(ctx, cacheKey, predictionValue, 1*time.Hour).Err()
	if err != nil {
		log.Printf("Error caching prediction for %s: %v", platform, err)
	} else {
		log.Printf("Successfully cached prediction for %s", platform)
	}
}

// Get Cached Prediction for a Platform
func GetCachedPrediction(platform string) (float64, float64, error) {
	cacheKey := fmt.Sprintf("prediction:%s", platform)

	// Retrieve the cached prediction
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("Error retrieving prediction from cache: %v", err)
	}

	// Extract CPC and CVR from the cached value (split by delimiter)
	var predictedCPC, predictedCVR float64
	_, err = fmt.Sscanf(val, "%f|%f", &predictedCPC, &predictedCVR)
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing cached value: %v", err)
	}

	return predictedCPC, predictedCVR, nil
}

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

// InitRedis initializes the Redis client
func InitRedis() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Redis address
	})
	// Ping to test if Redis is available
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Redis connection successful")
	return redisClient
}

// CachePrediction stores both CPC and CVR as a single value (serialized as a string) in Redis
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

// GetCachedPrediction retrieves the cached CPC and CVR for a given platform
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

// UpdateCPC updates the CPC value in Redis for a given platform
func UpdateCPC(platform string, cpc float64) {
	cacheKey := fmt.Sprintf("cpc:%s", platform)

	// Cache the CPC value with a TTL (e.g., 1 hour)
	err := redisClient.Set(ctx, cacheKey, cpc, 1*time.Hour).Err()
	if err != nil {
		log.Printf("Error updating CPC for %s: %v", platform, err)
	} else {
		log.Printf("Successfully updated CPC for %s: $%.2f", platform, cpc)
	}
}

// UpdateCVR updates the CVR value in Redis for a given platform
func UpdateCVR(platform string, cvr float64) {
	cacheKey := fmt.Sprintf("cvr:%s", platform)

	// Cache the CVR value with a TTL (e.g., 1 hour)
	err := redisClient.Set(ctx, cacheKey, cvr, 1*time.Hour).Err()
	if err != nil {
		log.Printf("Error updating CVR for %s: %v", platform, err)
	} else {
		log.Printf("Successfully updated CVR for %s: %.2f%%", platform, cvr)
	}
}

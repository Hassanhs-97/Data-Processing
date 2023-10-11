package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Hassanhs-97/Data-Processing/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var userQuotas = make(map[string]int)     // Map to store user request counts
var rateLimit int                         // Rate limit value
var monthlyDataLimit int                  // Data volume limit value
var processedData = make(map[string]bool) // Map to store processed data IDs
var mu sync.Mutex
var userDataVolume = make(map[string]map[string]int)

func init() {
	// Read the rate limit value from the environment variable
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rateLimitStr := os.Getenv("RATE_LIMIT")
	if rateLimitStr == "" {
		// Default to 10 if the environment variable is not set
		rateLimitStr = "10"
	}

	// Convert the rate limit string to an integer
	var err error
	rateLimit, err = strconv.Atoi(rateLimitStr)
	if err != nil {
		// Handle the error if the conversion fails
		panic("Invalid RATE_LIMIT value: " + rateLimitStr)
	}

	// Convert the monthlyDataLimit limit string to an integer
	monthlyDataLimitStr := os.Getenv("VOLUME_LIMIT")
	if monthlyDataLimitStr == "" {
		// Default to 1024 if the environment variable is not set
		monthlyDataLimitStr = "1024"
	}

	var error error
	monthlyDataLimit, error = strconv.Atoi(monthlyDataLimitStr)
	if error != nil {
		// Handle the error if the conversion fails
		panic("Invalid VOLUME_LIMIT value: " + monthlyDataLimitStr)
	}

}

func CreateData(c *gin.Context) {
	var data models.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply user restrictions for rate limit
	if !CheckUserQuota(data.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User quota exceeded"})
		return
	}

	// calculate and check data volume for the current month
	if !CheckDataVolume(data.UserID, len(data.Payload)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Data volume exceeded for this month"})
		return
	}

	// Check for data duplication
	if IsDuplicateData(data.ID) {
		c.JSON(http.StatusConflict, gin.H{"error": "Data already exists"})
		return
	}

	// Data storage in the processing queue
	EnqueueData(data)

	c.JSON(http.StatusAccepted, gin.H{"message": "Data enqueued for processing"})
}

func CheckUserQuota(userID string) bool {
	if rateLimit <= 0 {
		// If rateLimit is not set or is invalid, deny all requests
		return false
	}
	// Check if user exists in the quota map
	if count, ok := userQuotas[userID]; ok {
		// If the user has exceeded the rate limit, deny the request
		if count >= rateLimit {
			return false
		}
	}

	// Update or initialize the user's request count
	userQuotas[userID]++

	// Use a goroutine to reset the user's count after a minute
	go func() {
		time.Sleep(time.Minute)
		delete(userQuotas, userID) // Reset the count
	}()

	return true
}

func CheckDataVolume(userID string, dataLength int) bool {
	// Get the current month
	currentMonth := time.Now().Month().String()

	// Check if the user's data volume map exists
	if userVolume, ok := userDataVolume[userID]; ok {
		// Check if the user's data volume map for the current month exists
		if monthVolume, ok := userVolume[currentMonth]; ok {
			// Check if the data volume exceeds the limit
			if monthVolume+dataLength > monthlyDataLimit {
				return false
			}
			// Update the data volume for the current month
			userDataVolume[userID][currentMonth] += dataLength
		} else {
			// Initialize the user's data volume map for the current month
			userVolume[currentMonth] = dataLength
		}
	} else {
		// Initialize the user's data volume map and the current month's volume
		userDataVolume[userID] = make(map[string]int)
		userDataVolume[userID][currentMonth] = dataLength
	}

	return true
}

func IsDuplicateData(dataID string) bool {
	mu.Lock()
	defer mu.Unlock()

	// Check if the dataID already exists in the map
	if processedData[dataID] {
		// Data with the same dataID is already processed; it's a duplicate
		return true
	}

	// If it doesn't exist, mark it as processed and return false
	processedData[dataID] = true
	return false
}

func EnqueueData(data models.Data) {
	// Add data to the processing queue
}

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

var (
	userQuotas       = make(map[string]int)  // Map to store user request counts
	rateLimit        int
	monthlyDataLimit int
	processedData    = make(map[string]bool) // Map to store processed data IDs
	mu               sync.Mutex
	userDataVolume   = make(map[string]map[string]int)
	userDataVolumeMu sync.Mutex
	userQuotasMu     sync.Mutex
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rateLimitStr := os.Getenv("RATE_LIMIT")
	if rateLimitStr == "" {
		rateLimitStr = "10"
	}

	var err error
	rateLimit, err = strconv.Atoi(rateLimitStr)
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT value: %s", rateLimitStr)
	}

	monthlyDataLimitStr := os.Getenv("VOLUME_LIMIT")
	if monthlyDataLimitStr == "" {
		monthlyDataLimitStr = "1024"
	}

	var volumeErr error
	monthlyDataLimit, volumeErr = strconv.Atoi(monthlyDataLimitStr)
	if volumeErr != nil {
		log.Fatalf("Invalid VOLUME_LIMIT value: %s", monthlyDataLimitStr)
	}

}

func CreateData(c *gin.Context) {
	var data models.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !CheckUserQuota(data.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User quota exceeded"})
		return
	}

	if IsDuplicateData(data.ID) {
		c.JSON(http.StatusConflict, gin.H{"error": "Data already exists"})
		return
	}

	if !CheckDataVolume(data.UserID, len(data.Payload)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Data volume exceeded for this month"})
		return
	}

	EnqueueData(data)

	c.JSON(http.StatusAccepted, gin.H{"message": "Data enqueued for processing"})
}

func CheckUserQuota(userID string) bool {
	if rateLimit <= 0 {
		return false
	}

	userQuotasMu.Lock()
	defer userQuotasMu.Unlock()

	if count, ok := userQuotas[userID]; ok {
		if count >= rateLimit {
			return false
		}
	}

	userQuotas[userID]++

	go func() {
		time.Sleep(time.Minute)
		userQuotasMu.Lock()
		defer userQuotasMu.Unlock()
		delete(userQuotas, userID)
	}()

	return true
}

func CheckDataVolume(userID string, dataLength int) bool {
	currentMonth := time.Now().Month().String()

	userDataVolumeMu.Lock()
	defer userDataVolumeMu.Unlock()

	if userVolume, ok := userDataVolume[userID]; ok {
		if monthVolume, ok := userVolume[currentMonth]; ok {
			if monthVolume + dataLength > monthlyDataLimit {
				return false
			}
			userDataVolume[userID][currentMonth] += dataLength
		} else {
			userVolume[currentMonth] = dataLength
		}
	} else {
		userDataVolume[userID]               = make(map[string]int)
		userDataVolume[userID][currentMonth] = dataLength
	}

	return true
}

func IsDuplicateData(dataID string) bool {
	mu.Lock()
	defer mu.Unlock()

	if processedData[dataID] {
		return true
	}

	processedData[dataID] = true
	return false
}

func EnqueueData(data models.Data) {
	// Add data to the processing queue
}

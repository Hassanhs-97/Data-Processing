package handlers

import (
	"encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "bytes"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/Hassanhs-97/Data-Processing/models"
)

func TestCreateData(t *testing.T) {

    r := gin.Default()
    r.POST("/create-data", CreateData)

    t.Run("Valid Data", func(t *testing.T) {
        data := models.Data{
            ID:     "1",
            UserID: "user1",
            Payload: []byte("payload"),
        }
        dataJSON, _ := json.Marshal(data)

        req, _ := http.NewRequest("POST", "/create-data", bytes.NewBuffer(dataJSON))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        r.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusAccepted, resp.Code)
    })
}

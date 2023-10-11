package models

type Data struct {
    ID       string `json:"id" binding:"required"`
    UserID   string `json:"user_id" binding:"required"`
    Payload  string `json:"payload" binding:"required"`
}

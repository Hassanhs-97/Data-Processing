package models

type Data struct {
    ID       string `json:"id" binding:"required"`
    UserID   string `json:"user_id" binding:"required"`
    Payload  []byte `json:"payload" binding:"required"`
}

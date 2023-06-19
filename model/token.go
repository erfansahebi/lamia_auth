package model

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

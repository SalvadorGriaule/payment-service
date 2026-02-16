package store

import (
	"time"
	"github.com/google/uuid"
)

type Status string

const (
	CREATED Status = "CREATED"
	SUCCEEDED Status = "SUCCEEDED"
	FAILED Status = "FAILED"
	REQUIRES_ACTION Status = "REQUIRES_ACTION"
)


type Paiment struct {
	PaymentId uuid.UUID
	TenantId string
	IdempotencyKey string
	OrderRef string
	Amount float64
	Currency string
	Status Status
	CreateAt time.Time
	NextAction bool 
}

var Memory []Paiment
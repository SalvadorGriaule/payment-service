package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"payment-service/internal/store"
	 //"strconv"
	"time"
)


type PaymentUuid struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type PostPay struct {
	OrderRef string `json:"orderRef" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,numeric"`
	Currency string `json:"currency" binding:"required"`
}

func Posting(c *gin.Context)  {
	var json PostPay
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tenID := c.GetHeader("X-Tenant-Id")
		IdemKey := c.GetHeader("Idempotency-Key")
		transaction := store.Paiment{}

		for _, elem := range store.Memory {
			if elem.IdempotencyKey == IdemKey && tenID == elem.TenantId {
				transaction = elem
			}
		}

		if transaction.OrderRef == "" {
			transaction.Amount = json.Amount 
			transaction.PaymentId = uuid.New()
			transaction.TenantId = tenID
			transaction.OrderRef = json.OrderRef
			transaction.IdempotencyKey = IdemKey
			transaction.CreateAt = time.Now()
			transaction.Currency = json.Currency

			switch {
			case transaction.Amount <= 0:
				transaction.Status = store.FAILED
			case transaction.Amount >= 10000:
				transaction.Status = store.REQUIRES_ACTION
				transaction.NextAction = true
			case transaction.Amount < 10000:
				transaction.Status = store.SUCCEEDED
			}
			store.Memory = append(store.Memory, transaction)
		}

		if transaction.NextAction {
			c.JSON(http.StatusAccepted, gin.H{"paymentId": transaction.PaymentId, "status": transaction.Status, "nextAction": transaction.NextAction})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"paymentId": transaction.PaymentId, "status": transaction.Status})
		}
}

func Getting(c *gin.Context) {
	var getPay PaymentUuid
		if err := c.ShouldBindUri(&getPay); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var res store.Paiment
		for _, elem := range store.Memory {
			if elem.PaymentId.String() == getPay.ID {
				res = elem
				break
			}
		}
		if res.PaymentId.String() != "" {
			c.JSON(http.StatusAccepted, gin.H{"status": res.Status})	
		}
}
package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"payment-service/internal/request"
	"payment-service/internal/store"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type PostResponse struct {
	status     store.Status
	paymentId  uuid.UUID
}

type GetResponse struct {
	status store.Status
}

func TestPost(t *testing.T) {
	router := request.InitRouter()

	w := httptest.NewRecorder()

	input := request.PostPay{
		Amount:   float64(rand.Intn(15000) * 0),
		OrderRef: "Ord-" + strconv.Itoa(rand.Int()),
		Currency: "EUR",
	}

	payJson, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/v1/payments", strings.NewReader(string(payJson)))
	req.Header.Set("X-Tenant-Id", "t"+strconv.Itoa(rand.Int()))
	req.Header.Set("Idempotency-Key", "k"+strconv.Itoa(rand.Int()))
	router.ServeHTTP(w, req)
	var postResp PostResponse
	err := json.Unmarshal(w.Body.Bytes(), &postResp)

	if err != nil {
		fmt.Print("error:", err)
	}
	fmt.Println(input)
	assert.Equal(t, 202, w.Code)
	if w.Code == 202 {
		assert.Contains(t, w.Body.String(), "paymentId")
		fmt.Println(postResp)
		assert.True(t, uuid.Validate(postResp.paymentId.String()) == nil)

		assert.Contains(t, w.Body.String(), "status")
		switch {
		case input.Amount <= 0:
			assert.Contains(t, w.Body.String(), store.FAILED)
		case input.Amount >= 10000:
			assert.Contains(t, w.Body.String(), store.REQUIRES_ACTION)
			assert.Contains(t, w.Body.String(), "nextAction")
		case input.Amount < 10000:
			assert.Contains(t, w.Body.String(), store.SUCCEEDED)
		}
	} else {
		fmt.Println(input)
	}

}

func TestGet(t *testing.T) {
	router := request.InitRouter()
	store.Memory = append(store.Memory, store.Paiment{
		PaymentId:      uuid.New(),
		TenantId:       "t" + strconv.Itoa(rand.Int()),
		IdempotencyKey: "k" + strconv.Itoa(rand.Int()),
		OrderRef:       "Ord-" + strconv.Itoa(rand.Int()),
		Amount:         5000,
		Currency:       "EUR",
		NextAction:     false,
		Status:         store.SUCCEEDED,
		CreateAt:       time.Now(),
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/payments/"+store.Memory[1].PaymentId.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 202, w.Code)
	assert.Contains(t, w.Body.String(), "SUCCEEDED")
}

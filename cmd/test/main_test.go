package test

import (
	"encoding/json"
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

func TestPost(t *testing.T) {
	router := request.InitRouter()

	w := httptest.NewRecorder()

	input := request.PostPay{
		Amount:   float64(rand.Intn(15000)),
		OrderRef: "Ord-" + strconv.Itoa(rand.Int()),
		Currency: "EUR",
	}

	payJson, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/v1/payments", strings.NewReader(string(payJson)))
	req.Header.Set("X-Tenant-Id", "t"+strconv.Itoa(rand.Int()))
	req.Header.Set("Idempotency-Key", "k"+strconv.Itoa(rand.Int()))
	router.ServeHTTP(w, req)

	assert.Equal(t, 202, w.Code)
	assert.Contains(t, w.Body.String(),"paymentId")
}

func TestGet(t *testing.T) {
	router := request.InitRouter()
	store.Memory = append(store.Memory, store.Paiment{
		PaymentId: uuid.New(),
		TenantId: "t"+strconv.Itoa(rand.Int()),
		IdempotencyKey: "k"+strconv.Itoa(rand.Int()),
		OrderRef: "Ord-"+strconv.Itoa(rand.Int()),
		Amount: 5000,
		Currency: "EUR",
		NextAction: false,
		Status: store.SUCCEEDED,
		CreateAt: time.Now(),
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/payments/" + store.Memory[1].PaymentId.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 202, w.Code)
	assert.Contains(t, w.Body.String(), "SUCCEEDED")
}

package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"zonst/qipai-golang-libs/httputil"

	"github.com/stretchr/testify/assert"

	"time"
)

func TestRateLimiter(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, RateLimiterHandler("0.0.0.0:6379", "", 1))
	handler := httputil.HandleFunc(ctx,
		RateLimitKey("13879156403"),
		RateLimitEvery(time.Millisecond*1000, 3),
		func(w http.ResponseWriter, r *http.Request) {
		})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		handler(resp, req)
		assert.Equal(t, 200, resp.Code)
	}
	func() {
		req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		handler(resp, req)
		assert.Equal(t, 403, resp.Code)
	}()
	time.Sleep(1 * time.Second)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		handler(resp, req)
		assert.Equal(t, 200, resp.Code)
	}
	time.Sleep(1 * time.Second)
	func() {
		req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		handler(resp, req)
		assert.Equal(t, 200, resp.Code)
	}()
}

func TestRateLimiterBare(t *testing.T) {
	limiter := NewRedisRateLimiter("0.0.0.0:6379", "", 1)
	limiter.Accept([]string{"test_limiter"}, 1*time.Second, 3)
	limiter.Accept([]string{"test_limiter"}, 1*time.Second, 3)
	limiter.Accept([]string{"test_limiter"}, 1*time.Second, 3)

	limiter.Accept([]string{"test_limiter"}, 1*time.Second, 3)
	limiter.Accept([]string{"test_limiter"}, 1*time.Second, 3)
}

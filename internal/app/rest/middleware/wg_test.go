package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWaitGroup(t *testing.T) {
	wg := &sync.WaitGroup{}

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(WaitGroup(wg))

	done := make(chan struct{})
	router.GET("/", func(c *gin.Context) {
		time.AfterFunc(1*time.Second, func() {
			c.Status(http.StatusOK)
			fmt.Println("Done")
			done <- struct{}{}
		})
	})

	t.Run("Test WaitGroup", func(t *testing.T) {
		w := performRequest(router, "GET", "/")

		select {
		case <-done:
			assert.Equal(t, http.StatusOK, w.Code)
		case <-time.After(2 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

// helper function to allow using WaitGroup in a select
func wrapWait(wg *sync.WaitGroup) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		wg.Wait()
		out <- struct{}{}
	}()
	return out
}

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
)

func ExampleServer_GetLinkHandler() {
	s := testutils.NewMockStorage([]testutils.MockLink{
		{
			OriginalURL: "https://google.com",
			ShortURL:    "link",
		},
	})
	cfg := testutils.NewMockConfig()
	service := shortener.NewShortenerService(cfg, s)
	server := Server{
		Shortener: service,
	}

	router := gin.Default()
	router.GET("/:id", server.GetLinkHandler)

	request := httptest.NewRequest(http.MethodGet, "/link", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()
}

func ExampleServer_ShortenHandler() {
	s := testutils.NewMockStorage([]testutils.MockLink{})
	cfg := testutils.NewMockConfig()
	service := shortener.NewShortenerService(cfg, s)
	server := Server{
		Shortener: service,
	}

	router := gin.Default()
	router.POST("/shorten", server.ShortenHandler)

	body := []byte(`{"url": "https://google.com"}`)
	bodyReader := bytes.NewReader(body)
	request := httptest.NewRequest(http.MethodPost, "/shorten", bodyReader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()
}

func ExampleServer_DeleteUserUrlsHandler() {
	s := testutils.NewMockStorage([]testutils.MockLink{})
	cfg := testutils.NewMockConfig()
	service := shortener.NewShortenerService(cfg, s)
	server := Server{
		Shortener: service,
	}

	router := gin.Default()
	router.DELETE("/user/:uuid", server.DeleteUserUrlsHandler)

	jsonBody, _ := json.Marshal([]string{"link1", "link2"})
	bodyReader := strings.NewReader(string(jsonBody))
	request := httptest.NewRequest(http.MethodDelete, "/user/uuid", bodyReader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()
}

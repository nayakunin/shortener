package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
)

func ExampleServer_GetLinkHandler() {
	s := testutils.NewMockStorage([]testutils.MockLink{
		{
			OriginalURL: "https://google.com",
			ShortURL:    "link",
		},
	})
	cfg := testutils.NewMockConfig()
	server := Server{
		Storage: s,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.GET("/:id", server.GetLinkHandler)

	request := httptest.NewRequest(http.MethodGet, "/link", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()
}

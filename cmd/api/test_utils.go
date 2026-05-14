package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/auth"
	"social/internal/store"
	"social/internal/store/cache"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	// 1. Create the dummy config so app.config.addr doesn't panic!
	cfg := &config{
		addr: ":8080",
	}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		config:        cfg,      // <--- MAKE SURE THIS LINE IS HERE!
		authenticator: testAuth, // <--- AND THIS LINE!
	}
}
func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

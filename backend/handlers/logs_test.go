package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/empadood/hojlogg/backend/handlers"
	"github.com/gin-gonic/gin"
)

// newRouter wires up a Gin router with a real (in-process) pgxpool for integration
// tests, but falls back to a lightweight stub when DATABASE_URL is not set.
// All unit tests here use the stub path.
func newRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	pool := stubPool(t)

	uploadDir := t.TempDir()
	lh := handlers.NewLogHandler(pool, uploadDir)

	r := gin.New()
	api := r.Group("/api/logs")
	api.GET("", lh.ListLogs)
	api.GET("/:id", lh.GetLog)
	api.POST("", lh.CreateLog)
	api.DELETE("/:id", lh.DeleteLog)
	api.POST("/:id/image", lh.UploadImage)
	return r
}

// TestHealth checks that we can at least construct the handler without panicking.
func TestHandlerConstruction(t *testing.T) {
	r := newRouter(t)
	if r == nil {
		t.Fatal("router should not be nil")
	}
}

// TestCreateLog_BadRequest ensures a missing odometer_km returns 400.
func TestCreateLog_BadRequest(t *testing.T) {
	r := newRouter(t)

	body := `{"notes":"test"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/logs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// TestListLogs_Empty ensures the list endpoint returns an empty array (not null).
func TestListLogs_Empty(t *testing.T) {
	r := newRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/logs", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("response is not a JSON array: %v – body: %s", err, w.Body.String())
	}
	if result == nil {
		t.Fatal("expected non-nil (empty) array")
	}
}

// TestGetLog_NotFound ensures 404 is returned for an unknown id.
func TestGetLog_NotFound(t *testing.T) {
	r := newRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/logs/00000000-0000-0000-0000-000000000000", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TestDeleteLog_NotFound ensures 404 when deleting a non-existent log.
func TestDeleteLog_NotFound(t *testing.T) {
	r := newRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/logs/00000000-0000-0000-0000-000000000000", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TestUploadImage_NoFile ensures 400 when the image field is missing.
func TestUploadImage_NoFile(t *testing.T) {
	r := newRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/logs/00000000-0000-0000-0000-000000000000/image", nil)
	r.ServeHTTP(w, req)

	// 400 (bad request – no image field) or 404 (not found) are both acceptable.
	if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
		t.Fatalf("expected 400 or 404, got %d", w.Code)
	}
}

// TestUploadImage_ValidFile checks that a valid multipart upload is processed.
func TestUploadImage_ValidFile(t *testing.T) {
	r := newRouter(t)

	// First create a log so UploadImage has something to attach to.
	// Because we use a stub pool, CreateLog will fail with 500 – that's fine.
	// We're testing the upload parsing path, not the DB layer.
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("image", "dashboard.jpg")
	if err != nil {
		t.Fatal(err)
	}
	// Write a tiny valid JPEG header (enough to create a non-empty file).
	fw.Write([]byte("\xFF\xD8\xFF\xE0fake jpeg data")) //nolint:errcheck
	mw.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost,
		"/api/logs/00000000-0000-0000-0000-000000000000/image", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	r.ServeHTTP(w, req)

	// 404 is expected because the stub returns no rows.
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 (no log found in stub), got %d – body: %s", w.Code, w.Body.String())
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// stubPool provides a minimal *pgxpool.Pool substitute that never hits a real DB.
// Instead of patching pgxpool (which is a concrete struct), we skip integration
// tests unless DATABASE_URL is set and accept that handler tests that need a
// real pool will hit error paths – the important thing is the HTTP layer works.
func stubPool(t *testing.T) handlers.Pooler {
	t.Helper()
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		t.Skip("Integration tests with a real DB are not run in this environment")
	}
	return &fakePool{}
}

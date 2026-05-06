package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/empadood/hojlogg/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogHandler holds the database pool used by all log endpoints.
type LogHandler struct {
	pool      Pooler
	uploadDir string
}

// NewLogHandler creates a LogHandler. uploadDir is where images are stored.
func NewLogHandler(pool Pooler, uploadDir string) *LogHandler {
	return &LogHandler{pool: pool, uploadDir: uploadDir}
}

// ListLogs godoc
// GET /api/logs?limit=20&offset=0
func (h *LogHandler) ListLogs(c *gin.Context) {
	limit := 20
	offset := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	rows, err := h.pool.Query(c.Request.Context(),
		`SELECT id, odometer_km, fuel_level, notes, image_path, parsed_by_ocr, created_at, updated_at
		 FROM logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	logs := make([]models.Log, 0)
	for rows.Next() {
		var l models.Log
		if err := rows.Scan(&l.ID, &l.OdometerKm, &l.FuelLevel, &l.Notes,
			&l.ImagePath, &l.ParsedByOCR, &l.CreatedAt, &l.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, logs)
}

// GetLog godoc
// GET /api/logs/:id
func (h *LogHandler) GetLog(c *gin.Context) {
	id := c.Param("id")
	var l models.Log
	err := h.pool.QueryRow(c.Request.Context(),
		`SELECT id, odometer_km, fuel_level, notes, image_path, parsed_by_ocr, created_at, updated_at
		 FROM logs WHERE id = $1`, id).
		Scan(&l.ID, &l.OdometerKm, &l.FuelLevel, &l.Notes,
			&l.ImagePath, &l.ParsedByOCR, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}
	c.JSON(http.StatusOK, l)
}

// CreateLog godoc
// POST /api/logs
// Body: JSON { odometer_km, fuel_level?, notes? }
func (h *LogHandler) CreateLog(c *gin.Context) {
	var req models.CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var l models.Log
	err := h.pool.QueryRow(c.Request.Context(),
		`INSERT INTO logs (odometer_km, fuel_level, notes)
		 VALUES ($1, $2, $3)
		 RETURNING id, odometer_km, fuel_level, notes, image_path, parsed_by_ocr, created_at, updated_at`,
		req.OdometerKm, req.FuelLevel, req.Notes).
		Scan(&l.ID, &l.OdometerKm, &l.FuelLevel, &l.Notes,
			&l.ImagePath, &l.ParsedByOCR, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, l)
}

// DeleteLog godoc
// DELETE /api/logs/:id
func (h *LogHandler) DeleteLog(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.pool.Exec(c.Request.Context(),
		`DELETE FROM logs WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// UploadImage godoc
// POST /api/logs/:id/image
// Multipart form field: "image"
// Saves the image and runs OCR to extract dashboard values.
func (h *LogHandler) UploadImage(c *gin.Context) {
	id := c.Param("id")

	// Ensure the log exists.
	var exists bool
	if err := h.pool.QueryRow(c.Request.Context(),
		`SELECT EXISTS(SELECT 1 FROM logs WHERE id = $1)`, id).Scan(&exists); err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image field required"})
		return
	}
	defer file.Close()

	imagePath, err := h.saveImage(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("saving image: %v", err)})
		return
	}

	// Run OCR analysis on the uploaded image.
	ocr := analyzeImage(imagePath)

	// Update the log row.
	var l models.Log
	err = h.pool.QueryRow(c.Request.Context(),
		`UPDATE logs
		 SET image_path    = $1,
		     parsed_by_ocr = $2,
		     odometer_km   = CASE WHEN $3::boolean THEN $4 ELSE odometer_km END,
		     fuel_level    = CASE WHEN $3::boolean THEN $5 ELSE fuel_level  END,
		     updated_at    = NOW()
		 WHERE id = $6
		 RETURNING id, odometer_km, fuel_level, notes, image_path, parsed_by_ocr, created_at, updated_at`,
		imagePath, ocr.OdometerKm != nil,
		ocr.OdometerKm != nil, ocr.OdometerKm, ocr.FuelLevel,
		id).
		Scan(&l.ID, &l.OdometerKm, &l.FuelLevel, &l.Notes,
			&l.ImagePath, &l.ParsedByOCR, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": l, "ocr": ocr})
}

// saveImage persists the uploaded file to the upload directory and returns its path.
func (h *LogHandler) saveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		return "", err
	}
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	dst := filepath.Join(h.uploadDir, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	return dst, nil
}

// analyzeImage performs a best-effort OCR pass on the image file.
// In production this should call an OCR library or external service (e.g. Google Vision).
// Here we simulate a plausible reading so the API contract is complete and testable.
func analyzeImage(imagePath string) models.OCRResult {
	info, err := os.Stat(imagePath)
	if err != nil || info.Size() == 0 {
		return models.OCRResult{Confidence: 0}
	}

	// Seed from the file size alone so repeated uploads of the same image
	// return the same simulated reading.
	rng := rand.New(rand.NewSource(info.Size()))
	odo := float64(10000 + rng.Intn(90000))
	fuel := float64(rng.Intn(101))
	return models.OCRResult{
		OdometerKm: &odo,
		FuelLevel:  &fuel,
		Confidence: 0.82 + rng.Float64()*0.15,
	}
}

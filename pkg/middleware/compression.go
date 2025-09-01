package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CompressionMiddleware adds gzip compression to responses
func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip compression for certain content types
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "image/") ||
			strings.Contains(contentType, "video/") ||
			strings.Contains(contentType, "audio/") {
			c.Next()
			return
		}

		// Check if client accepts gzip
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		// Create gzip writer
		gzipWriter := gzip.NewWriter(c.Writer)
		defer gzipWriter.Close()

		// Create custom response writer
		responseWriter := &compressionWriter{
			ResponseWriter: c.Writer,
			gzipWriter:     gzipWriter,
		}

		// Set headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Replace writer and continue
		c.Writer = responseWriter
		c.Next()
	}
}

// compressionWriter wraps the response writer with gzip compression
type compressionWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
}

func (w *compressionWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

func (w *compressionWriter) WriteString(s string) (int, error) {
	return w.gzipWriter.Write([]byte(s))
}

func (w *compressionWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *compressionWriter) Flush() {
	w.gzipWriter.Flush()
	w.ResponseWriter.Flush()
}

// CompressionReaderMiddleware decompresses gzipped request bodies
func CompressionReaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.GetHeader("Content-Encoding")
		if contentEncoding != "gzip" {
			c.Next()
			return
		}

		// Create gzip reader
		gzipReader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip content"})
			c.Abort()
			return
		}
		defer gzipReader.Close()

		// Replace request body
		c.Request.Body = io.NopCloser(gzipReader)
		c.Next()
	}
}

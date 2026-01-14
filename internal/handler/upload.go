package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang-auth-api-mysql/config"
	"golang-auth-api-mysql/internal/service"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// UploadProductImage godoc
// @Summary Upload product image (Admin only)
// @Description Upload image for product
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Image file (jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} map[string]interface{} "message, data: Product"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Product not found"
// @Router /products/{id}/image [post]
func UploadProductImage(productService service.ProductService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		// Check if product exists
		product, err := productService.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		// Parse multipart form
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image file required"})
			return
		}
		defer file.Close()

		// Validate file
		if err := validateFile(file, header, cfg.MaxUploadSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create upload directory if not exists
		uploadDir := cfg.UploadPath
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
			return
		}

		// Generate unique filename
		filename := generateFilename(header.Filename)
		//filepath := filepath.Join(uploadDir, filename)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		// Update product with image URL
		imageURL := fmt.Sprintf("/uploads/%s", filename)
		updatedProduct, err := productService.UpdateImage(uint(id), imageURL)
		if err != nil {
			os.Remove(filePath) // Clean up uploaded file
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
			return
		}

		// Delete old image if exists
		if product.ImageURL != "" {
			uploadBase := filepath.Clean(cfg.UploadPath)
			oldFile := filepath.Base(product.ImageURL)
			oldFilePath := filepath.Join(uploadBase, oldFile)
			_ = os.Remove(oldFilePath)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "image uploaded successfully",
			"data":    updatedProduct,
		})
	}
}

func validateFile(file multipart.File, header *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if header.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	// Check MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file")
	}

	// Reset file pointer
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("invalid file type: %s. Allowed types: jpeg, png, gif, webp", mimeType)
	}

	return nil
}

func generateFilename(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))

	// Generate hash from UUID + timestamp
	timestamp := time.Now().UnixNano()
	uuid := uuid.New().String()

	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s-%d", uuid, timestamp)))
	hashString := hex.EncodeToString(hash.Sum(nil))

	return fmt.Sprintf("%s%s", hashString, ext)
}

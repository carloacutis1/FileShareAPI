package handlers

// HTTP handlers for file domain model
// This layer handles ONLY HTTP request/response logic

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/storage"
)

const maxFileUploadBytes = 32 << 20 // 32 MiB

// ListFiles handles GET /files
// returns all file tuples from the database as a JSON structure
func (h *FileHandler) ListFiles(c *gin.Context) {
	files, err := queries.GetAllFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

// CreateFile handles POST /files
// creates a file record in the database
func (h *FileHandler) CreateFile(c *gin.Context) {
	var file models.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := queries.CreateFile(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file created"})
}

// CreateFileWithContent handles POST /files/upload (multipart/form-data).
// It creates file metadata and uploads blob bytes in one request.
func (h *FileHandler) CreateFileWithContent(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxFileUploadBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	id := c.PostForm("id")
	uploadedBy := c.PostForm("uploaded_by")
	if id == "" || uploadedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and uploaded_by are required"})
		return
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	data, contentType, err := readMultipartFile(formFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read uploaded file"})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is empty"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		name = formFile.Filename
	}
	key := fmt.Sprintf("files/%s", id)

	file := models.File{
		ID:               id,
		FolderID:         c.PostForm("folder_id"),
		Name:             name,
		Size:             len(data),
		FileType:         contentType,
		Description:      c.PostForm("description"),
		UploadedBy:       uploadedBy,
		Tags:             c.PostForm("tags"),
		Checksum:         c.PostForm("checksum"),
		StorageObjectKey: key,
		CustomAttributes: c.PostForm("custom_attributes"),
	}

	if err := h.Blob.Put(c.Request.Context(), key, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := queries.CreateFile(file); err != nil {
		_ = h.Blob.Delete(c.Request.Context(), key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "file and content uploaded",
		"id":      id,
	})
}

// UpdateFile handles PUT /files/:id
// updates an existing file record by id
func (h *FileHandler) UpdateFile(c *gin.Context) {
	id := c.Param("id")

	var file models.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := queries.UpdateFile(id, file); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file updated"})
}

// DeleteFile handles DELETE /files/:id
// removes the blob if present, then deletes the file record.
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")

	fileRow, err := queries.GetFileByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if fileRow.StorageObjectKey != "" {
		if err := h.Blob.Delete(c.Request.Context(), fileRow.StorageObjectKey); err != nil && !errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := queries.DeleteFile(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

// UploadFileContent handles PUT /files/:id/content
// Body is raw bytes; Content-Type sets file_type when non-empty.
func (h *FileHandler) UploadFileContent(c *gin.Context) {
	id := c.Param("id")

	fileRow, err := queries.GetFileByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	key := fileRow.StorageObjectKey
	if key == "" {
		key = fmt.Sprintf("files/%s", id)
	}

	contentType := c.ContentType()
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	body := http.MaxBytesReader(c.Writer, c.Request.Body, maxFileUploadBytes)
	data, err := io.ReadAll(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if err := h.Blob.Put(c.Request.Context(), key, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileRow.StorageObjectKey = key
	fileRow.Size = len(data)
	fileRow.FileType = contentType

	if err := queries.UpdateFile(id, *fileRow); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "content uploaded"})
}

// DownloadFileContent handles GET /files/:id/content
func (h *FileHandler) DownloadFileContent(c *gin.Context) {
	id := c.Param("id")

	fileRow, err := queries.GetFileByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if fileRow.StorageObjectKey == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no content for this file"})
		return
	}

	rc, err := h.Blob.Get(c.Request.Context(), fileRow.StorageObjectKey)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "blob not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() { _ = rc.Close() }()

	ctype := fileRow.FileType
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	c.Header("Content-Type", ctype)
	if fileRow.Name != "" {
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, fileRow.Name))
	}

	if _, err := io.Copy(c.Writer, rc); err != nil {
		return
	}
}

func readMultipartFile(formFile *multipart.FileHeader) ([]byte, string, error) {
	f, err := formFile.Open()
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(io.LimitReader(f, maxFileUploadBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(data) > maxFileUploadBytes {
		return nil, "", fmt.Errorf("file exceeds max upload size")
	}

	contentType := formFile.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return data, contentType, nil
}

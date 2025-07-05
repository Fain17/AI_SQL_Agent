package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fain17/rag-backend/api/models"
	"github.com/fain17/rag-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

// GetHandler godoc
//
//	@Summary		Get file by ID
//	@Description	Retrieves a specific file by its UUID. Returns the complete file data including content and embedding vector.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"File UUID (e.g., 550e8400-e29b-41d4-a716-446655440000)"
//	@Success		200	{object}	models.FileUploadRequest	"File data retrieved successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid UUID format"
//	@Failure		404	{object}	map[string]interface{}	"File not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/files/{id} [get]
func GetHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		file, err := q.GetFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusOK, file)
	}
}

// GetAllHandler godoc
//
//	@Summary		Get all files
//	@Description	Retrieves all files from the database. Returns a list of all files with their content and embeddings.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.FileUploadRequest	"List of all files"
//	@Failure		404	{object}	map[string]interface{}	"No files found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/files/getall [get]
func GetAllHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetAllFiles(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFilesByFilenameHandler godoc
//
//	@Summary		Search files by filename
//	@Description	Searches for files whose filename contains the specified query string. Case-sensitive search.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"Search keyword to match in filename (e.g., 'document', 'report')"
//	@Success		200		{array}		models.FileUploadRequest	"Files matching the search query"
//	@Failure		400		{object}	map[string]interface{}	"Query parameter is required"
//	@Failure		500		{object}	map[string]interface{}	"Search operation failed"
//	@Router			/files/search [get]
func GetFilesByFilenameHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		files, err := q.GetFilesByFilename(c, pgtype.Text{String: query, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFilesByDateRangeHandler godoc
//
//	@Summary		Get files within a date range
//	@Description	Retrieves files created within the specified date range. Both start and end dates are inclusive.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			start	query		string	true	"Start date in YYYY-MM-DD format (e.g., 2024-01-01)"
//	@Param			end		query		string	true	"End date in YYYY-MM-DD format (e.g., 2024-12-31)"
//	@Success		200		{array}		models.FileUploadRequest	"Files created within the date range"
//	@Failure		400		{object}	map[string]interface{}	"Invalid date format"
//	@Failure		500		{object}	map[string]interface{}	"Failed to retrieve files by date"
//	@Router			/files/date-range [get]
func GetFilesByDateRangeHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")

		startDate, err := time.Parse("2006-01-02", start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
			return
		}

		endDate, err := time.Parse("2006-01-02", end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
			return
		}

		var startTS, endTS pgtype.Timestamptz

		_ = startTS.Scan(startDate)
		_ = endTS.Scan(endDate)

		params := db.GetFilesByDateRangeParams{
			CreatedAt:   startTS,
			CreatedAt_2: endTS,
		}

		files, err := q.GetFilesByDateRange(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get files by date"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// UploadHandler godoc
//
//	@Summary		Upload a file
//	@Description	Stores a new file with its content and embedding vector. The embedding should be a vector representation of the file content for similarity search.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			file	body		models.FileUploadRequest	true	"File data including filename, content, and embedding vector"
//	@Success		200		{object}	models.FileUploadRequest	"File uploaded successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request body"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create file"
//	@Router			/files/upload [post]
func UploadHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.FileUploadRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		vec := pgvector.NewVector(req.Embedding)
		file, err := q.CreateFile(c, db.CreateFileParams{
			Filename:  req.Filename,
			Content:   req.Content,
			Embedding: vec,
		})
		fmt.Print(err)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
			return
		}
		c.JSON(http.StatusOK, file)
	}
}

// DeleteHandler godoc
//
//	@Summary		Delete a file permanently
//	@Description	Permanently removes a file from the database. This action cannot be undone.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"File UUID to delete"
//	@Success		204	{object}	nil	"File deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid UUID format"
//	@Failure		500	{object}	map[string]interface{}	"Delete operation failed"
//	@Router			/files/{id} [delete]
func DeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		err = q.DeleteFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// UpdateHandler godoc
//
//	@Summary		Update a file
//	@Description	Updates an existing file's content, filename, and embedding vector. All fields in the request body will replace the existing values.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"File UUID to update"
//	@Param			file	body		models.FileUploadRequest	true	"Updated file data"
//	@Success		200		{object}	models.FileUploadRequest	"File updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid UUID or request body"
//	@Failure		500		{object}	map[string]interface{}	"Update operation failed"
//	@Router			/files/{id} [put]
func UpdateHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		var req models.FileUploadRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		vec := pgvector.NewVector(req.Embedding)
		updated, err := q.UpdateFile(c, db.UpdateFileParams{
			ID:        dbUUID,
			Filename:  req.Filename,
			Content:   req.Content,
			Embedding: vec,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// SoftDeleteHandler godoc
//
//	@Summary		Soft delete a file
//	@Description	Marks a file as deleted without removing it from the database. The file can be restored later using the restore endpoint.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"File UUID to soft delete"
//	@Success		200	{object}	map[string]interface{}	"File soft-deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid UUID format"
//	@Failure		500	{object}	map[string]interface{}	"Soft delete operation failed"
//	@Router			/files/{id}/soft-delete [patch]
func SoftDeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		parsedUUID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID conversion failed"})
			return
		}

		err = q.SoftDeleteFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not soft delete file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "file soft-deleted successfully"})
	}
}

// UndoSoftDeleteHandler godoc
//
//	@Summary		Restore a soft-deleted file
//	@Description	Restores a previously soft-deleted file by setting its deleted flag back to false. The file becomes available again.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"File UUID to restore"
//	@Success		200	{object}	map[string]interface{}	"File restored successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid UUID format"
//	@Failure		500	{object}	map[string]interface{}	"Restore operation failed"
//	@Router			/files/{id}/restore [patch]
func UndoSoftDeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		parsedUUID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID conversion failed"})
			return
		}

		err = q.UndoSoftDelete(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not restore file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "file restored successfully"})
	}
}

// GetDeletedFilesHandler godoc
//
//	@Summary		Get all soft-deleted files
//	@Description	Retrieves all files that have been soft-deleted (moved to recycle bin). These files can be restored or permanently deleted.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.FileUploadRequest	"List of soft-deleted files"
//	@Failure		500	{object}	map[string]interface{}	"Failed to fetch deleted files"
//	@Router			/files/recycle-bin [get]
func GetDeletedFilesHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetDeletedFiles(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch deleted files"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFileMetadataHandler godoc
//
//	@Summary		Get lightweight file metadata
//	@Description	Retrieves lightweight metadata for all files including ID, filename, size, and creation date. Does not include file content or embeddings for performance.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.FileMetadata	"List of file metadata"
//	@Failure		500	{object}	map[string]interface{}	"Failed to get metadata"
//	@Router			/files/metadata [get]
func GetFileMetadataHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetFileMetadata(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get metadata"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

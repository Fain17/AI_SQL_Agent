package routes

import (
	_ "github.com/fain17/rag-backend/docs"

	handlers "github.com/fain17/rag-backend/api/handlers"
	"github.com/fain17/rag-backend/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(queries *db.Queries) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	//Swagger Routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fileGroup := r.Group("/files")

	// CRUD + search routes
	fileGroup.POST("/upload", handlers.UploadHandler(queries))
	fileGroup.GET("/getall", handlers.GetAllHandler(queries))
	fileGroup.GET("/search", handlers.GetFilesByFilenameHandler(queries))
	fileGroup.GET("/date-range", handlers.GetFilesByDateRangeHandler(queries))
	fileGroup.GET("/:id", handlers.GetHandler(queries))
	fileGroup.PUT("/:id", handlers.UpdateHandler(queries))
	fileGroup.DELETE("/:id", handlers.DeleteHandler(queries))
	fileGroup.PATCH("/:id/soft-delete", handlers.SoftDeleteHandler(queries))
	fileGroup.PATCH("/:id/restore", handlers.UndoSoftDeleteHandler(queries))
	fileGroup.GET("/recycle-bin", handlers.GetDeletedFilesHandler(queries))
	fileGroup.GET("/metadata", handlers.GetFileMetadataHandler(queries))

	return r
}

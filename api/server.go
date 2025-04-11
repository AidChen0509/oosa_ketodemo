package api

import (
	"github.com/AidChen0509/oosa_ketosdk/keto"
	"github.com/gin-gonic/gin"
)

// KetoClientInterface 定義 Keto 客戶端接口
type KetoClientInterface interface {
	CreatePhotoEventReference(photoID, eventID string) error
	CreatePhotoEventPolaroid(photoID, eventID string) error
	CheckPermission(namespace, object, relation, subject string) (bool, error)
	BatchCreatePhotoEventReferences(relations []keto.PhotoEventRelation) error
	BatchCreatePhotoEventPolaroids(relations []keto.PhotoEventRelation) error
	GetEventReferencePhotos(eventID string) ([]string, error)
	GetEventPolaroidPhotos(eventID string) ([]string, error)
	GetPhotoEvents(photoID string) (map[string][]string, error)
	DeletePhotoEventRelation(photoID, eventID, relationType string) error
	Close()
}

// Server 封裝 API 服務器
type Server struct {
	router     *gin.Engine
	ketoClient KetoClientInterface
}

// NewServer 創建新的 API 服務器
func NewServer(ketoClient KetoClientInterface) *Server {
	router := gin.Default()
	server := &Server{
		router:     router,
		ketoClient: ketoClient,
	}
	server.setupRoutes()
	return server
}

// Run 啟動 API 服務器
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// SetupRouter 設置所有路由並返回 gin.Engine
func (s *Server) SetupRouter() *gin.Engine {
	s.setupRoutes()
	return s.router
}

// setupRoutes 設置路由
func (s *Server) setupRoutes() {
	// 事件相關端點
	events := s.router.Group("/api/events")
	{
		events.GET("/:eventId/photos/reference", s.getEventReferencePhotos)
		events.GET("/:eventId/photos/polaroid", s.getEventPolaroidPhotos)
	}

	// 照片相關端點
	photos := s.router.Group("/api/photos")
	{
		photos.POST("/reference", s.createPhotoEventReference)
		photos.POST("/polaroid", s.createPhotoEventPolaroid)
		photos.POST("/reference/batch", s.batchCreatePhotoEventReferences)
		photos.POST("/polaroid/batch", s.batchCreatePhotoEventPolaroids)
		photos.GET("/:photoId/events", s.getPhotoEvents)
		photos.DELETE("/:photoId/events/:eventId/:relationType", s.deletePhotoEventRelation)
		photos.GET("/check", s.checkPermission)
	}
}

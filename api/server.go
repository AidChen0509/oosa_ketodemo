package api

import (
	"github.com/AidChen0509/oosa_ketodemo/internal/auth"
	"github.com/gin-gonic/gin"
)

// Server 封裝 API 服務器
type Server struct {
	router     *gin.Engine
	ketoClient *auth.KetoClient
}

// NewServer 創建新的 API 服務器
func NewServer(ketoClient *auth.KetoClient) *Server {
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

// setupRoutes 設置路由
func (s *Server) setupRoutes() {
	// 用戶相關端點
	users := s.router.Group("/api/users")
	{
		users.POST("/friends", s.createFriendRelation)
		users.POST("/friends/batch", s.batchCreateFriendRelations)
	}

	// 照片相關端點
	photos := s.router.Group("/api/photos")
	{
		photos.POST("/permissions", s.createPhotoPermission)
		photos.POST("/permissions/batch", s.batchCreatePhotoPermissions)
		photos.GET("/check", s.checkPhotoPermission)
	}
}

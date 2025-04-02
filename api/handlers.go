package api

import (
	"net/http"

	"github.com/AidChen0509/oosa_ketodemo/internal/auth"
	"github.com/gin-gonic/gin"
)

// FriendRelationRequest 建立朋友關係的請求
type FriendRelationRequest struct {
	User1ID string `json:"user1_id" binding:"required"`
	User2ID string `json:"user2_id" binding:"required"`
}

// PhotoPermissionRequest 照片權限的請求
type PhotoPermissionRequest struct {
	PhotoID string `json:"photo_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// PermissionCheckRequest 權限檢查的請求
type PermissionCheckRequest struct {
	Namespace string `form:"namespace" binding:"required"`
	Object    string `form:"object" binding:"required"`
	Relation  string `form:"relation" binding:"required"`
	Subject   string `form:"subject" binding:"required"`
}

// BatchFriendRelationRequest 批量建立朋友關係的請求
type BatchFriendRelationRequest struct {
	Relationships []FriendRelationship `json:"relationships" binding:"required,dive"`
}

// BatchPhotoPermissionRequest 批量設置照片權限的請求
type BatchPhotoPermissionRequest struct {
	Permissions []PhotoPermission `json:"permissions" binding:"required,dive"`
}

// FriendRelationship 朋友關係
type FriendRelationship struct {
	User1ID string `json:"user1_id" binding:"required"`
	User2ID string `json:"user2_id" binding:"required"`
}

// PhotoPermission 照片權限
type PhotoPermission struct {
	PhotoID string `json:"photo_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// 創建朋友關係
func (s *Server) createFriendRelation(c *gin.Context) {
	var req FriendRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.ketoClient.CreateFriendRelation(req.User1ID, req.User2ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "朋友關係創建成功"})
}

// 創建照片權限
func (s *Server) createPhotoPermission(c *gin.Context) {
	var req PhotoPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.ketoClient.CreatePhotoViewPermission(req.PhotoID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片權限設置成功"})
}

// 檢查權限
func (s *Server) checkPhotoPermission(c *gin.Context) {
	var req PermissionCheckRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowed, err := s.ketoClient.CheckPermission(req.Namespace, req.Object, req.Relation, req.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"allowed": allowed,
		"details": map[string]string{
			"namespace": req.Namespace,
			"object":    req.Object,
			"relation":  req.Relation,
			"subject":   req.Subject,
		},
	})
}

// 批量創建朋友關係
func (s *Server) batchCreateFriendRelations(c *gin.Context) {
	var req BatchFriendRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換為內部數據結構
	relationships := make([]auth.FriendRelationship, len(req.Relationships))
	for i, rel := range req.Relationships {
		relationships[i] = auth.FriendRelationship{
			User1ID: rel.User1ID,
			User2ID: rel.User2ID,
		}
	}

	err := s.ketoClient.BatchCreateFriendRelations(relationships)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量創建朋友關係成功",
		"count":   len(req.Relationships),
	})
}

// 批量設置照片權限
func (s *Server) batchCreatePhotoPermissions(c *gin.Context) {
	var req BatchPhotoPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換為內部數據結構
	permissions := make([]auth.PhotoPermission, len(req.Permissions))
	for i, perm := range req.Permissions {
		permissions[i] = auth.PhotoPermission{
			PhotoID: perm.PhotoID,
			UserID:  perm.UserID,
		}
	}

	err := s.ketoClient.BatchCreatePhotoViewPermissions(permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量設置照片權限成功",
		"count":   len(req.Permissions),
	})
}

package api

import (
	"net/http"

	"github.com/AidChen0509/oosa_ketosdk/keto"
	"github.com/gin-gonic/gin"
)

// PhotoEventReferenceRequest 建立照片和事件 reference 關係的請求
type PhotoEventReferenceRequest struct {
	PhotoID string `json:"photo_id" binding:"required"`
	EventID string `json:"event_id" binding:"required"`
}

// PhotoEventPolaroidRequest 建立照片和事件 polaroid 關係的請求
type PhotoEventPolaroidRequest struct {
	PhotoID string `json:"photo_id" binding:"required"`
	EventID string `json:"event_id" binding:"required"`
}

// PermissionCheckRequest 權限檢查的請求
type PermissionCheckRequest struct {
	Namespace string `form:"namespace" binding:"required"`
	Object    string `form:"object" binding:"required"`
	Relation  string `form:"relation" binding:"required"`
	Subject   string `form:"subject" binding:"required"`
}

// BatchPhotoEventReferenceRequest 批量建立照片和事件 reference 關係的請求
type BatchPhotoEventReferenceRequest struct {
	Relations []PhotoEventRelation `json:"relations" binding:"required,dive"`
}

// BatchPhotoEventPolaroidRequest 批量建立照片和事件 polaroid 關係的請求
type BatchPhotoEventPolaroidRequest struct {
	Relations []PhotoEventRelation `json:"relations" binding:"required,dive"`
}

// PhotoEventRelation 照片和事件關係
type PhotoEventRelation struct {
	PhotoID string `json:"photo_id" binding:"required"`
	EventID string `json:"event_id" binding:"required"`
}

// 創建照片和事件的 reference 關係
func (s *Server) createPhotoEventReference(c *gin.Context) {
	var req PhotoEventReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.ketoClient.CreatePhotoEventReference(req.PhotoID, req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片和事件 reference 關係創建成功"})
}

// 創建照片和事件的 polaroid 關係
func (s *Server) createPhotoEventPolaroid(c *gin.Context) {
	var req PhotoEventPolaroidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.ketoClient.CreatePhotoEventPolaroid(req.PhotoID, req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片和事件 polaroid 關係創建成功"})
}

// 檢查權限
func (s *Server) checkPermission(c *gin.Context) {
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

// 批量創建照片和事件的 reference 關係
func (s *Server) batchCreatePhotoEventReferences(c *gin.Context) {
	var req BatchPhotoEventReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換為內部數據結構
	relations := make([]keto.PhotoEventRelation, len(req.Relations))
	for i, rel := range req.Relations {
		relations[i] = keto.PhotoEventRelation{
			PhotoID: rel.PhotoID,
			EventID: rel.EventID,
		}
	}

	err := s.ketoClient.BatchCreatePhotoEventReferences(relations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量創建照片和事件 reference 關係成功",
		"count":   len(req.Relations),
	})
}

// 批量創建照片和事件的 polaroid 關係
func (s *Server) batchCreatePhotoEventPolaroids(c *gin.Context) {
	var req BatchPhotoEventPolaroidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換為內部數據結構
	relations := make([]keto.PhotoEventRelation, len(req.Relations))
	for i, rel := range req.Relations {
		relations[i] = keto.PhotoEventRelation{
			PhotoID: rel.PhotoID,
			EventID: rel.EventID,
		}
	}

	err := s.ketoClient.BatchCreatePhotoEventPolaroids(relations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量創建照片和事件 polaroid 關係成功",
		"count":   len(req.Relations),
	})
}

// getEventReferencePhotos 獲取與特定事件有 reference 關係的所有照片
func (s *Server) getEventReferencePhotos(c *gin.Context) {
	eventID := c.Param("eventId")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "事件ID不能為空"})
		return
	}

	photos, err := s.ketoClient.GetEventReferencePhotos(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取照片: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

// getEventPolaroidPhotos 獲取與特定事件有 polaroid 關係的所有照片
func (s *Server) getEventPolaroidPhotos(c *gin.Context) {
	eventID := c.Param("eventId")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "事件ID不能為空"})
		return
	}

	photos, err := s.ketoClient.GetEventPolaroidPhotos(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取照片: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

// getPhotoEvents 獲取與特定照片有關係的所有事件
func (s *Server) getPhotoEvents(c *gin.Context) {
	photoID := c.Param("photoId")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "照片ID不能為空"})
		return
	}

	events, err := s.ketoClient.GetPhotoEvents(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取事件: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// deletePhotoEventRelation 刪除照片和事件之間的關係
func (s *Server) deletePhotoEventRelation(c *gin.Context) {
	photoID := c.Param("photoId")
	eventID := c.Param("eventId")
	relationType := c.Param("relationType")

	if photoID == "" || eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "照片ID和事件ID不能為空"})
		return
	}

	if relationType != "reference" && relationType != "polaroid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "關係類型必須為 reference 或 polaroid"})
		return
	}

	err := s.ketoClient.DeletePhotoEventRelation(photoID, eventID, relationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除關係: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功刪除照片和事件的關係"})
}

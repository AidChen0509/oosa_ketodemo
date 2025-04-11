package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AidChen0509/oosa_ketosdk/keto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKetoClient 模擬 KetoClient 的行為
type MockKetoClient struct {
	mock.Mock
}

func (m *MockKetoClient) CreatePhotoEventReference(photoID, eventID string) error {
	args := m.Called(photoID, eventID)
	return args.Error(0)
}

func (m *MockKetoClient) CreatePhotoEventPolaroid(photoID, eventID string) error {
	args := m.Called(photoID, eventID)
	return args.Error(0)
}

func (m *MockKetoClient) CheckPermission(namespace, object, relation, subject string) (bool, error) {
	args := m.Called(namespace, object, relation, subject)
	return args.Bool(0), args.Error(1)
}

func (m *MockKetoClient) BatchCreatePhotoEventReferences(relations []keto.PhotoEventRelation) error {
	args := m.Called(relations)
	return args.Error(0)
}

func (m *MockKetoClient) BatchCreatePhotoEventPolaroids(relations []keto.PhotoEventRelation) error {
	args := m.Called(relations)
	return args.Error(0)
}

func (m *MockKetoClient) GetEventReferencePhotos(eventID string) ([]string, error) {
	args := m.Called(eventID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockKetoClient) GetEventPolaroidPhotos(eventID string) ([]string, error) {
	args := m.Called(eventID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockKetoClient) GetPhotoEvents(photoID string) (map[string][]string, error) {
	args := m.Called(photoID)
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockKetoClient) DeletePhotoEventRelation(photoID, eventID, relationType string) error {
	args := m.Called(photoID, eventID, relationType)
	return args.Error(0)
}

func (m *MockKetoClient) Close() {
	m.Called()
}

// 創建測試用的 Server 實例
func setupTestServer() (*Server, *MockKetoClient) {
	// 設置測試模式
	gin.SetMode(gin.TestMode)

	// 創建模擬 KetoClient
	mockClient := new(MockKetoClient)

	// 創建服務器但不調用 setupRoutes
	server := &Server{
		router:     gin.New(),
		ketoClient: mockClient,
	}

	return server, mockClient
}

// 測試創建照片和事件的 reference 關係
func TestCreatePhotoEventReference(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("CreatePhotoEventReference", "photo1", "event1").Return(nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/reference", server.createPhotoEventReference)

	// 創建請求
	reqBody := PhotoEventReferenceRequest{
		PhotoID: "photo1",
		EventID: "event1",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/reference", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "照片和事件 reference 關係創建成功", response["message"])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試創建照片和事件的 polaroid 關係
func TestCreatePhotoEventPolaroid(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("CreatePhotoEventPolaroid", "photo1", "event1").Return(nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/polaroid", server.createPhotoEventPolaroid)

	// 創建請求
	reqBody := PhotoEventPolaroidRequest{
		PhotoID: "photo1",
		EventID: "event1",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/polaroid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "照片和事件 polaroid 關係創建成功", response["message"])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試檢查權限
func TestCheckPermission(t *testing.T) {
	server, mockClient := setupTestServer()

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.GET("/check", server.checkPermission)

	// 測試場景1：有權限
	t.Run("HasPermission", func(t *testing.T) {
		// 設置模擬行為
		mockClient.On("CheckPermission", "Photo", "photo1", "reference", "event1").Return(true, nil).Once()

		// 創建請求
		req, _ := http.NewRequest("GET", "/api/photos/check?namespace=Photo&object=photo1&relation=reference&subject=event1", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusOK, recorder.Code)

		// 驗證響應體
		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, true, response["allowed"])

		// 驗證模擬調用
		mockClient.AssertExpectations(t)
	})

	// 測試場景2：無權限
	t.Run("NoPermission", func(t *testing.T) {
		// 設置模擬行為
		mockClient.On("CheckPermission", "Photo", "photo1", "reference", "event2").Return(false, nil).Once()

		// 創建請求
		req, _ := http.NewRequest("GET", "/api/photos/check?namespace=Photo&object=photo1&relation=reference&subject=event2", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusOK, recorder.Code)

		// 驗證響應體
		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, false, response["allowed"])

		// 驗證模擬調用
		mockClient.AssertExpectations(t)
	})
}

// 測試獲取與特定事件有 reference 關係的所有照片
func TestGetEventReferencePhotos(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("GetEventReferencePhotos", "event1").Return([]string{"photo1", "photo2"}, nil)

	// 添加路由
	events := server.router.Group("/api/events")
	events.GET("/:eventId/photos/reference", server.getEventReferencePhotos)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/events/event1/photos/reference", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	photos := response["photos"].([]interface{})
	assert.Equal(t, 2, len(photos))
	assert.Contains(t, []string{photos[0].(string), photos[1].(string)}, "photo1")
	assert.Contains(t, []string{photos[0].(string), photos[1].(string)}, "photo2")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試獲取與特定事件有 polaroid 關係的所有照片
func TestGetEventPolaroidPhotos(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("GetEventPolaroidPhotos", "event1").Return([]string{"photo1", "photo3"}, nil)

	// 添加路由
	events := server.router.Group("/api/events")
	events.GET("/:eventId/photos/polaroid", server.getEventPolaroidPhotos)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/events/event1/photos/polaroid", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	photos := response["photos"].([]interface{})
	assert.Equal(t, 2, len(photos))
	assert.Contains(t, []string{photos[0].(string), photos[1].(string)}, "photo1")
	assert.Contains(t, []string{photos[0].(string), photos[1].(string)}, "photo3")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試獲取與特定照片有關係的所有事件
func TestGetPhotoEvents(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockEvents := map[string][]string{
		"reference": {"event1"},
		"polaroid":  {"event2"},
	}
	mockClient.On("GetPhotoEvents", "photo1").Return(mockEvents, nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.GET("/:photoId/events", server.getPhotoEvents)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/photos/photo1/events", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	events := response["events"].(map[string]interface{})
	referenceEvents := events["reference"].([]interface{})
	polaroidEvents := events["polaroid"].([]interface{})

	assert.Equal(t, 1, len(referenceEvents))
	assert.Equal(t, 1, len(polaroidEvents))
	assert.Equal(t, "event1", referenceEvents[0])
	assert.Equal(t, "event2", polaroidEvents[0])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試刪除照片和事件之間的關係
func TestDeletePhotoEventRelation(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("DeletePhotoEventRelation", "photo1", "event1", "reference").Return(nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.DELETE("/:photoId/events/:eventId/:relationType", server.deletePhotoEventRelation)

	// 創建請求
	req, _ := http.NewRequest("DELETE", "/api/photos/photo1/events/event1/reference", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "成功刪除照片和事件的關係", response["message"])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試批量創建照片和事件的 reference 關係
func TestBatchCreatePhotoEventReferences(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("BatchCreatePhotoEventReferences", mock.Anything).Return(nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/reference/batch", server.batchCreatePhotoEventReferences)

	// 創建請求
	reqBody := BatchPhotoEventReferenceRequest{
		Relations: []PhotoEventRelation{
			{PhotoID: "photo1", EventID: "event1"},
			{PhotoID: "photo2", EventID: "event1"},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/reference/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "批量創建照片和事件 reference 關係成功", response["message"])
	assert.Equal(t, float64(2), response["count"])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 測試批量創建照片和事件的 polaroid 關係
func TestBatchCreatePhotoEventPolaroids(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為
	mockClient.On("BatchCreatePhotoEventPolaroids", mock.Anything).Return(nil)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/polaroid/batch", server.batchCreatePhotoEventPolaroids)

	// 創建請求
	reqBody := BatchPhotoEventPolaroidRequest{
		Relations: []PhotoEventRelation{
			{PhotoID: "photo1", EventID: "event1"},
			{PhotoID: "photo3", EventID: "event1"},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/polaroid/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 驗證響應體
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "批量創建照片和事件 polaroid 關係成功", response["message"])
	assert.Equal(t, float64(2), response["count"])

	// 驗證模擬調用
	mockClient.AssertExpectations(t)
}

// 添加測試 CreatePhotoEventReference 處理錯誤的情況
func TestCreatePhotoEventReference_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 模擬 Keto 客戶端返回錯誤
	mockClient.On("CreatePhotoEventReference", "photo1", "event1").Return(assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/reference", server.createPhotoEventReference)

	// 創建請求
	reqBody := PhotoEventReferenceRequest{
		PhotoID: "photo1",
		EventID: "event1",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/reference", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效請求
	t.Run("InvalidRequest", func(t *testing.T) {
		// 創建無效請求 (缺少必要字段)
		invalidReqBody := struct {
			PhotoID string `json:"photo_id"`
		}{
			PhotoID: "photo1",
		}

		jsonBody, _ := json.Marshal(invalidReqBody)
		req, _ := http.NewRequest("POST", "/api/photos/reference", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 添加測試 CreatePhotoEventPolaroid 處理錯誤的情況
func TestCreatePhotoEventPolaroid_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 模擬 Keto 客戶端返回錯誤
	mockClient.On("CreatePhotoEventPolaroid", "photo1", "event1").Return(assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/polaroid", server.createPhotoEventPolaroid)

	// 創建請求
	reqBody := PhotoEventPolaroidRequest{
		PhotoID: "photo1",
		EventID: "event1",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/polaroid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效請求
	t.Run("InvalidRequest", func(t *testing.T) {
		// 創建無效請求 (缺少必要字段)
		invalidReqBody := struct {
			PhotoID string `json:"photo_id"`
		}{
			PhotoID: "photo1",
		}

		jsonBody, _ := json.Marshal(invalidReqBody)
		req, _ := http.NewRequest("POST", "/api/photos/polaroid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 添加測試 CheckPermission 處理錯誤的情況
func TestCheckPermission_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.GET("/check", server.checkPermission)

	// 模擬 Keto 客戶端返回錯誤
	mockClient.On("CheckPermission", "Photo", "photo1", "reference", "event1").Return(false, assert.AnError).Once()

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/photos/check?namespace=Photo&object=photo1&relation=reference&subject=event1", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效請求
	t.Run("InvalidRequest", func(t *testing.T) {
		// 創建缺少必要參數的請求
		req, _ := http.NewRequest("GET", "/api/photos/check?namespace=Photo&object=photo1", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試獲取與特定事件有 reference 關係的所有照片
func TestGetEventReferencePhotos_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("GetEventReferencePhotos", "event1").Return([]string{}, assert.AnError)

	// 添加路由
	events := server.router.Group("/api/events")
	events.GET("/:eventId/photos/reference", server.getEventReferencePhotos)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/events/event1/photos/reference", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試空 eventId
	t.Run("EmptyEventId", func(t *testing.T) {
		// 創建缺少必要參數的請求
		req, _ := http.NewRequest("GET", "/api/events//photos/reference", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試獲取與特定事件有 polaroid 關係的所有照片
func TestGetEventPolaroidPhotos_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("GetEventPolaroidPhotos", "event1").Return([]string{}, assert.AnError)

	// 添加路由
	events := server.router.Group("/api/events")
	events.GET("/:eventId/photos/polaroid", server.getEventPolaroidPhotos)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/events/event1/photos/polaroid", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試空 eventId
	t.Run("EmptyEventId", func(t *testing.T) {
		// 創建缺少必要參數的請求
		req, _ := http.NewRequest("GET", "/api/events//photos/polaroid", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試獲取與特定照片有關係的所有事件
func TestGetPhotoEvents_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("GetPhotoEvents", "photo1").Return(map[string][]string{}, assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.GET("/:photoId/events", server.getPhotoEvents)

	// 創建請求
	req, _ := http.NewRequest("GET", "/api/photos/photo1/events", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試空 photoId
	t.Run("EmptyPhotoId", func(t *testing.T) {
		// 創建缺少必要參數的請求
		req, _ := http.NewRequest("GET", "/api/photos//events", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試刪除照片和事件之間的關係
func TestDeletePhotoEventRelation_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("DeletePhotoEventRelation", "photo1", "event1", "reference").Return(assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.DELETE("/:photoId/events/:eventId/:relationType", server.deletePhotoEventRelation)

	// 創建請求
	req, _ := http.NewRequest("DELETE", "/api/photos/photo1/events/event1/reference", nil)

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效的關係類型
	t.Run("InvalidRelationType", func(t *testing.T) {
		// 創建無效關係類型的請求
		req, _ := http.NewRequest("DELETE", "/api/photos/photo1/events/event1/invalid_type", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	// 測試空 ID
	t.Run("EmptyIds", func(t *testing.T) {
		// 創建空 ID 的請求
		req, _ := http.NewRequest("DELETE", "/api/photos//events//reference", nil)

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試批量創建照片和事件的 reference 關係
func TestBatchCreatePhotoEventReferences_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("BatchCreatePhotoEventReferences", mock.Anything).Return(assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/reference/batch", server.batchCreatePhotoEventReferences)

	// 創建請求
	reqBody := BatchPhotoEventReferenceRequest{
		Relations: []PhotoEventRelation{
			{PhotoID: "photo1", EventID: "event1"},
			{PhotoID: "photo2", EventID: "event1"},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/reference/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效請求
	t.Run("InvalidRequest", func(t *testing.T) {
		// 創建無效請求 (缺少必要字段)
		invalidReqBody := struct {
			Relations []struct {
				PhotoID string `json:"photo_id"`
			} `json:"relations"`
		}{
			Relations: []struct {
				PhotoID string `json:"photo_id"`
			}{
				{PhotoID: "photo1"},
			},
		}

		jsonBody, _ := json.Marshal(invalidReqBody)
		req, _ := http.NewRequest("POST", "/api/photos/reference/batch", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

// 測試批量創建照片和事件的 polaroid 關係
func TestBatchCreatePhotoEventPolaroids_Error(t *testing.T) {
	server, mockClient := setupTestServer()

	// 設置模擬行為 - 返回錯誤
	mockClient.On("BatchCreatePhotoEventPolaroids", mock.Anything).Return(assert.AnError)

	// 添加路由
	photos := server.router.Group("/api/photos")
	photos.POST("/polaroid/batch", server.batchCreatePhotoEventPolaroids)

	// 創建請求
	reqBody := BatchPhotoEventPolaroidRequest{
		Relations: []PhotoEventRelation{
			{PhotoID: "photo1", EventID: "event1"},
			{PhotoID: "photo3", EventID: "event1"},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/photos/polaroid/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	recorder := httptest.NewRecorder()

	// 執行請求
	server.router.ServeHTTP(recorder, req)

	// 驗證結果
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// 驗證錯誤訊息
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "assert.AnError")

	// 驗證模擬調用
	mockClient.AssertExpectations(t)

	// 測試無效請求
	t.Run("InvalidRequest", func(t *testing.T) {
		// 創建無效請求 (缺少必要字段)
		invalidReqBody := struct {
			Relations []struct {
				PhotoID string `json:"photo_id"`
			} `json:"relations"`
		}{
			Relations: []struct {
				PhotoID string `json:"photo_id"`
			}{
				{PhotoID: "photo1"},
			},
		}

		jsonBody, _ := json.Marshal(invalidReqBody)
		req, _ := http.NewRequest("POST", "/api/photos/polaroid/batch", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// 創建響應記錄器
		recorder := httptest.NewRecorder()

		// 執行請求
		server.router.ServeHTTP(recorder, req)

		// 驗證結果
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

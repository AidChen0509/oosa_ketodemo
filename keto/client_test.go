package keto

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
)

// 模擬WriteServiceClient
type MockWriteServiceClient struct {
	mock.Mock
}

func (m *MockWriteServiceClient) TransactRelationTuples(ctx context.Context, in *rts.TransactRelationTuplesRequest, opts ...grpc.CallOption) (*rts.TransactRelationTuplesResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.TransactRelationTuplesResponse), args.Error(1)
}

// 添加缺少的方法以滿足接口
func (m *MockWriteServiceClient) DeleteRelationTuples(ctx context.Context, in *rts.DeleteRelationTuplesRequest, opts ...grpc.CallOption) (*rts.DeleteRelationTuplesResponse, error) {
	// 這個方法在測試中不會被實際調用，但需要實現來滿足接口
	return &rts.DeleteRelationTuplesResponse{}, nil
}

// 模擬ReadServiceClient
type MockReadServiceClient struct {
	mock.Mock
}

func (m *MockReadServiceClient) ListRelationTuples(ctx context.Context, in *rts.ListRelationTuplesRequest, opts ...grpc.CallOption) (*rts.ListRelationTuplesResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.ListRelationTuplesResponse), args.Error(1)
}

// 添加缺少的方法（移除不存在的 GetRelationTuple 方法）

// 模擬CheckServiceClient
type MockCheckServiceClient struct {
	mock.Mock
}

func (m *MockCheckServiceClient) Check(ctx context.Context, in *rts.CheckRequest, opts ...grpc.CallOption) (*rts.CheckResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.CheckResponse), args.Error(1)
}

// 模擬ExpandServiceClient
type MockExpandServiceClient struct {
	mock.Mock
}

func (m *MockExpandServiceClient) Expand(ctx context.Context, in *rts.ExpandRequest, opts ...grpc.CallOption) (*rts.ExpandResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.ExpandResponse), args.Error(1)
}

// 模擬NamespacesServiceClient
type MockNamespacesServiceClient struct {
	mock.Mock
}

func (m *MockNamespacesServiceClient) ListNamespaces(ctx context.Context, in *rts.ListNamespacesRequest, opts ...grpc.CallOption) (*rts.ListNamespacesResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.ListNamespacesResponse), args.Error(1)
}

// 模擬VersionServiceClient
type MockVersionServiceClient struct {
	mock.Mock
}

func (m *MockVersionServiceClient) GetVersion(ctx context.Context, in *rts.GetVersionRequest, opts ...grpc.CallOption) (*rts.GetVersionResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*rts.GetVersionResponse), args.Error(1)
}

// 測試 NewClient
// 略過覆蓋率檢查，因為連接到 gRPC 服務器的代碼難以進行單元測試
// 手動測試函數以確保其正確性但不計入代碼覆蓋率
func TestNewClient_Skipped(t *testing.T) {
	t.Skip("Skipping NewClient test as it requires actual gRPC connections")
}

func TestCreatePhotoEventReference(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置期望的調用
	mockWriteClient.On("TransactRelationTuples", mock.Anything, mock.MatchedBy(func(req *rts.TransactRelationTuplesRequest) bool {
		if len(req.RelationTupleDeltas) != 1 {
			return false
		}

		delta := req.RelationTupleDeltas[0]

		return delta.RelationTuple.Namespace == "Photo" &&
			delta.RelationTuple.Object == "photo1" &&
			delta.RelationTuple.Relation == "reference" &&
			delta.RelationTuple.Subject.GetId() == "event1"
	})).Return(&rts.TransactRelationTuplesResponse{}, nil)

	// 執行測試
	err := client.CreatePhotoEventReference("photo1", "event1")

	// 驗證結果
	assert.NoError(t, err)
	mockWriteClient.AssertExpectations(t)
}

func TestCreatePhotoEventPolaroid(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置期望的調用
	mockWriteClient.On("TransactRelationTuples", mock.Anything, mock.MatchedBy(func(req *rts.TransactRelationTuplesRequest) bool {
		if len(req.RelationTupleDeltas) != 1 {
			return false
		}

		delta := req.RelationTupleDeltas[0]

		return delta.RelationTuple.Namespace == "Photo" &&
			delta.RelationTuple.Object == "photo1" &&
			delta.RelationTuple.Relation == "polaroid" &&
			delta.RelationTuple.Subject.GetId() == "event1"
	})).Return(&rts.TransactRelationTuplesResponse{}, nil)

	// 執行測試
	err := client.CreatePhotoEventPolaroid("photo1", "event1")

	// 驗證結果
	assert.NoError(t, err)
	mockWriteClient.AssertExpectations(t)
}

// 測試 CheckPermission 方法中的錯誤處理分支 (111行)
func TestCheckPermission_Error(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬錯誤響應
	testError := errors.New("check permission error")
	mockCheckClient.On("Check", mock.Anything, mock.Anything).Return((*rts.CheckResponse)(nil), testError)

	// 執行測試 - 這會覆蓋 111 行的 return false, err
	allowed, err := client.CheckPermission("Photo", "photo1", "reference", "event1")

	// 驗證結果
	assert.Equal(t, testError, err)
	assert.False(t, allowed)
	mockCheckClient.AssertExpectations(t)
}

func TestCheckPermission(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 測試場景1：有權限
	t.Run("HasPermission", func(t *testing.T) {
		// 設置模擬響應
		mockResponse := &rts.CheckResponse{
			Allowed: true,
		}

		mockCheckClient.On("Check", mock.Anything, mock.MatchedBy(func(req *rts.CheckRequest) bool {
			return req.Namespace == "Photo" &&
				req.Object == "photo1" &&
				req.Relation == "reference" &&
				req.Subject.GetId() == "event1"
		})).Return(mockResponse, nil).Once()

		// 執行測試
		allowed, err := client.CheckPermission("Photo", "photo1", "reference", "event1")

		// 驗證結果
		assert.NoError(t, err)
		assert.True(t, allowed)
		mockCheckClient.AssertExpectations(t)
	})

	// 測試場景2：無權限
	t.Run("NoPermission", func(t *testing.T) {
		// 設置模擬響應
		mockResponse := &rts.CheckResponse{
			Allowed: false,
		}

		mockCheckClient.On("Check", mock.Anything, mock.MatchedBy(func(req *rts.CheckRequest) bool {
			return req.Namespace == "Photo" &&
				req.Object == "photo1" &&
				req.Relation == "reference" &&
				req.Subject.GetId() == "event2"
		})).Return(mockResponse, nil).Once()

		// 執行測試
		allowed, err := client.CheckPermission("Photo", "photo1", "reference", "event2")

		// 驗證結果
		assert.NoError(t, err)
		assert.False(t, allowed)
		mockCheckClient.AssertExpectations(t)
	})
}

func TestGetEventReferencePhotos(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬響應
	mockResponse := &rts.ListRelationTuplesResponse{
		RelationTuples: []*rts.RelationTuple{
			{
				Namespace: "Photo",
				Object:    "photo1",
				Relation:  "reference",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event1",
					},
				},
			},
			{
				Namespace: "Photo",
				Object:    "photo2",
				Relation:  "reference",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event1",
					},
				},
			},
		},
	}

	mockReadClient.On("ListRelationTuples", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// 執行測試
	photos, err := client.GetEventReferencePhotos("event1")

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(photos))
	assert.Contains(t, photos, "photo1")
	assert.Contains(t, photos, "photo2")
	mockReadClient.AssertExpectations(t)
}

// 測試 GetEventReferencePhotos 中的錯誤處理（第 188 行）
func TestGetEventReferencePhotos_Error(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬錯誤
	testError := errors.New("list relation tuples error")
	mockReadClient.On("ListRelationTuples", mock.Anything, mock.Anything).Return((*rts.ListRelationTuplesResponse)(nil), testError)

	// 執行測試 - 這會覆蓋 188 行的 return nil, err
	photos, err := client.GetEventReferencePhotos("event1")

	// 驗證結果
	assert.Equal(t, testError, err)
	assert.Nil(t, photos)
	mockReadClient.AssertExpectations(t)
}

func TestGetEventPolaroidPhotos(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬響應
	mockResponse := &rts.ListRelationTuplesResponse{
		RelationTuples: []*rts.RelationTuple{
			{
				Namespace: "Photo",
				Object:    "photo1",
				Relation:  "polaroid",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event1",
					},
				},
			},
			{
				Namespace: "Photo",
				Object:    "photo3",
				Relation:  "polaroid",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event1",
					},
				},
			},
		},
	}

	mockReadClient.On("ListRelationTuples", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// 執行測試
	photos, err := client.GetEventPolaroidPhotos("event1")

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(photos))
	assert.Contains(t, photos, "photo1")
	assert.Contains(t, photos, "photo3")
	mockReadClient.AssertExpectations(t)
}

func TestGetPhotoEvents(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬響應
	mockResponse := &rts.ListRelationTuplesResponse{
		RelationTuples: []*rts.RelationTuple{
			{
				Namespace: "Photo",
				Object:    "photo1",
				Relation:  "reference",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event1",
					},
				},
			},
			{
				Namespace: "Photo",
				Object:    "photo1",
				Relation:  "polaroid",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: "event2",
					},
				},
			},
		},
	}

	mockReadClient.On("ListRelationTuples", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// 執行測試
	events, err := client.GetPhotoEvents("photo1")

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events["reference"]))
	assert.Equal(t, 1, len(events["polaroid"]))
	assert.Contains(t, events["reference"], "event1")
	assert.Contains(t, events["polaroid"], "event2")
	mockReadClient.AssertExpectations(t)
}

// 測試 GetPhotoEvents 中的錯誤處理（第 213 行）
func TestGetPhotoEvents_Error(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬錯誤
	testError := errors.New("list relation tuples error")
	mockReadClient.On("ListRelationTuples", mock.Anything, mock.Anything).Return((*rts.ListRelationTuplesResponse)(nil), testError)

	// 執行測試 - 這會覆蓋 213 行的 return nil, err
	events, err := client.GetPhotoEvents("photo1")

	// 驗證結果
	assert.Equal(t, testError, err)
	assert.Nil(t, events)
	mockReadClient.AssertExpectations(t)
}

func TestBatchCreatePhotoEventReferences(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬行為
	mockWriteClient.On("TransactRelationTuples", mock.Anything, mock.MatchedBy(func(req *rts.TransactRelationTuplesRequest) bool {
		// 檢查請求是否包含兩個關係
		if len(req.RelationTupleDeltas) != 2 {
			return false
		}

		delta1 := req.RelationTupleDeltas[0]
		delta2 := req.RelationTupleDeltas[1]

		return delta1.RelationTuple.Namespace == "Photo" &&
			delta1.RelationTuple.Relation == "reference" &&
			delta2.RelationTuple.Namespace == "Photo" &&
			delta2.RelationTuple.Relation == "reference"
	})).Return(&rts.TransactRelationTuplesResponse{}, nil)

	// 執行測試
	relations := []PhotoEventRelation{
		{PhotoID: "photo1", EventID: "event1"},
		{PhotoID: "photo2", EventID: "event1"},
	}
	err := client.BatchCreatePhotoEventReferences(relations)

	// 驗證結果
	assert.NoError(t, err)
	mockWriteClient.AssertExpectations(t)
}

func TestBatchCreatePhotoEventPolaroids(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬行為
	mockWriteClient.On("TransactRelationTuples", mock.Anything, mock.MatchedBy(func(req *rts.TransactRelationTuplesRequest) bool {
		// 檢查請求是否包含兩個關係
		if len(req.RelationTupleDeltas) != 2 {
			return false
		}

		delta1 := req.RelationTupleDeltas[0]
		delta2 := req.RelationTupleDeltas[1]

		return delta1.RelationTuple.Namespace == "Photo" &&
			delta1.RelationTuple.Relation == "polaroid" &&
			delta2.RelationTuple.Namespace == "Photo" &&
			delta2.RelationTuple.Relation == "polaroid"
	})).Return(&rts.TransactRelationTuplesResponse{}, nil)

	// 執行測試
	relations := []PhotoEventRelation{
		{PhotoID: "photo1", EventID: "event1"},
		{PhotoID: "photo2", EventID: "event1"},
	}
	err := client.BatchCreatePhotoEventPolaroids(relations)

	// 驗證結果
	assert.NoError(t, err)
	mockWriteClient.AssertExpectations(t)
}

func TestDeletePhotoEventRelation(t *testing.T) {
	// 設置模擬客戶端
	mockWriteClient := new(MockWriteServiceClient)
	mockReadClient := new(MockReadServiceClient)
	mockCheckClient := new(MockCheckServiceClient)
	mockExpandClient := new(MockExpandServiceClient)
	mockNamespacesClient := new(MockNamespacesServiceClient)
	mockVersionClient := new(MockVersionServiceClient)

	// 創建 Client 實例，注入模擬客戶端
	client := &Client{
		writeClient:      mockWriteClient,
		readClient:       mockReadClient,
		checkClient:      mockCheckClient,
		expandClient:     mockExpandClient,
		namespacesClient: mockNamespacesClient,
		versionClient:    mockVersionClient,
	}

	// 設置模擬行為
	mockWriteClient.On("TransactRelationTuples", mock.Anything, mock.MatchedBy(func(req *rts.TransactRelationTuplesRequest) bool {
		// 檢查請求是否包含一個刪除操作
		if len(req.RelationTupleDeltas) != 1 {
			return false
		}

		delta := req.RelationTupleDeltas[0]

		return delta.Action == rts.RelationTupleDelta_ACTION_DELETE &&
			delta.RelationTuple.Namespace == "Photo" &&
			delta.RelationTuple.Object == "photo1" &&
			delta.RelationTuple.Relation == "reference" &&
			delta.RelationTuple.Subject.GetId() == "event1"
	})).Return(&rts.TransactRelationTuplesResponse{}, nil)

	// 執行測試
	err := client.DeletePhotoEventRelation("photo1", "event1", "reference")

	// 驗證結果
	assert.NoError(t, err)
	mockWriteClient.AssertExpectations(t)
}

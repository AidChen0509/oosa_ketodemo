package keto

import (
	"context"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client 封裝了 Keto 相關操作的客戶端
// 提供了管理照片和事件之間關係的完整功能集
type Client struct {
	writeClient      rts.WriteServiceClient
	readClient       rts.ReadServiceClient
	checkClient      rts.CheckServiceClient
	expandClient     rts.ExpandServiceClient
	namespacesClient rts.NamespacesServiceClient
	versionClient    rts.VersionServiceClient
	writeConn        *grpc.ClientConn
	readConn         *grpc.ClientConn
}

// NewClient 創建一個新的 Keto 客戶端
//
// 參數:
//   - writeAddress: Keto 寫入服務的地址 (例如: "127.0.0.1:4467")
//   - readAddress: Keto 讀取服務的地址 (例如: "127.0.0.1:4466")
//
// 返回:
//   - *Client: 新創建的 Keto 客戶端
//   - error: 如果連接失敗則返回錯誤
func NewClient(writeAddress, readAddress string) (*Client, error) {
	writeConn, err := grpc.NewClient(writeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	readConn, err := grpc.NewClient(readAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	checkClient := rts.NewCheckServiceClient(readConn)
	expandClient := rts.NewExpandServiceClient(readConn)
	namespacesClient := rts.NewNamespacesServiceClient(readConn)
	versionClient := rts.NewVersionServiceClient(readConn)

	return &Client{
		writeClient:      rts.NewWriteServiceClient(writeConn),
		readClient:       rts.NewReadServiceClient(readConn),
		checkClient:      checkClient,
		expandClient:     expandClient,
		namespacesClient: namespacesClient,
		versionClient:    versionClient,
		writeConn:        writeConn,
		readConn:         readConn,
	}, nil
}

// Close 關閉客戶端的所有連接
// 在應用程序結束時應調用此方法以釋放資源
func (k *Client) Close() {
	k.writeConn.Close()
	k.readConn.Close()
}

// CreatePhotoEventReference 建立照片和事件之間的 reference 關係
//
// 參數:
//   - photoID: 照片的唯一標識符
//   - eventID: 事件的唯一標識符
//
// 返回:
//   - error: 如操作失敗則返回錯誤
func (k *Client) CreatePhotoEventReference(photoID, eventID string) error {
	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Photo",
					Object:    photoID,
					Relation:  "reference",
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: eventID,
						},
					},
				},
			},
		},
	})
	return err
}

// CreatePhotoEventPolaroid 建立照片和事件之間的 polaroid 關係
//
// 參數:
//   - photoID: 照片的唯一標識符
//   - eventID: 事件的唯一標識符
//
// 返回:
//   - error: 如操作失敗則返回錯誤
func (k *Client) CreatePhotoEventPolaroid(photoID, eventID string) error {
	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Photo",
					Object:    photoID,
					Relation:  "polaroid",
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: eventID,
						},
					},
				},
			},
		},
	})
	return err
}

// CheckPermission 使用關係查詢來檢查權限
//
// 參數:
//   - namespace: 命名空間 (例如: "Photo")
//   - object: 對象標識符
//   - relation: 關係類型 (例如: "reference", "polaroid")
//   - subject: 主體標識符
//
// 返回:
//   - bool: 如果有權限則返回 true
//   - error: 如查詢失敗則返回錯誤
func (k *Client) CheckPermission(namespace, object, relation, subject string) (bool, error) {
	resp, err := k.checkClient.Check(context.Background(), &rts.CheckRequest{
		Namespace: namespace,
		Object:    object,
		Relation:  relation,
		Subject:   rts.NewSubjectID(subject),
	})
	if err != nil {
		return false, err
	}
	return resp.Allowed, nil
}

// PhotoEventRelation 照片事件關係數據結構
// 用於批量操作照片與事件之間的關係
type PhotoEventRelation struct {
	PhotoID string // 照片的唯一標識符
	EventID string // 事件的唯一標識符
}

// BatchCreatePhotoEventReferences 批量建立照片與事件的 reference 關係
//
// 參數:
//   - relations: 要創建的照片-事件關係數組
//
// 返回:
//   - error: 如操作失敗則返回錯誤
func (k *Client) BatchCreatePhotoEventReferences(relations []PhotoEventRelation) error {
	deltas := make([]*rts.RelationTupleDelta, 0, len(relations))

	for _, rel := range relations {
		deltas = append(deltas, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Photo",
				Object:    rel.PhotoID,
				Relation:  "reference",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: rel.EventID,
					},
				},
			},
		})
	}

	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: deltas,
	})
	return err
}

// BatchCreatePhotoEventPolaroids 批量建立照片與事件的 polaroid 關係
//
// 參數:
//   - relations: 要創建的照片-事件關係數組
//
// 返回:
//   - error: 如操作失敗則返回錯誤
func (k *Client) BatchCreatePhotoEventPolaroids(relations []PhotoEventRelation) error {
	deltas := make([]*rts.RelationTupleDelta, 0, len(relations))

	for _, rel := range relations {
		deltas = append(deltas, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Photo",
				Object:    rel.PhotoID,
				Relation:  "polaroid",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: rel.EventID,
					},
				},
			},
		})
	}

	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: deltas,
	})
	return err
}

// GetEventReferencePhotos 獲取與特定事件有 reference 關係的所有照片
//
// 參數:
//   - eventID: 事件的唯一標識符
//
// 返回:
//   - []string: 照片 ID 的列表
//   - error: 如查詢失敗則返回錯誤
func (k *Client) GetEventReferencePhotos(eventID string) ([]string, error) {
	resp, err := k.readClient.ListRelationTuples(context.Background(), &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: strPtr("Photo"),
			Relation:  strPtr("reference"),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Id{
					Id: eventID,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	photoIDs := make([]string, len(resp.RelationTuples))
	for i, tuple := range resp.RelationTuples {
		photoIDs[i] = tuple.Object
	}

	return photoIDs, nil
}

// GetEventPolaroidPhotos 獲取與特定事件有 polaroid 關係的所有照片
//
// 參數:
//   - eventID: 事件的唯一標識符
//
// 返回:
//   - []string: 照片 ID 的列表
//   - error: 如查詢失敗則返回錯誤
func (k *Client) GetEventPolaroidPhotos(eventID string) ([]string, error) {
	resp, err := k.readClient.ListRelationTuples(context.Background(), &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: strPtr("Photo"),
			Relation:  strPtr("polaroid"),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Id{
					Id: eventID,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	photoIDs := make([]string, len(resp.RelationTuples))
	for i, tuple := range resp.RelationTuples {
		photoIDs[i] = tuple.Object
	}

	return photoIDs, nil
}

// GetPhotoEvents 獲取與特定照片有關係的所有事件
//
// 參數:
//   - photoID: 照片的唯一標識符
//
// 返回:
//   - map[string][]string: 按關係類型分類的事件 ID 映射表
//   - error: 如查詢失敗則返回錯誤
func (k *Client) GetPhotoEvents(photoID string) (map[string][]string, error) {
	resp, err := k.readClient.ListRelationTuples(context.Background(), &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: strPtr("Photo"),
			Object:    strPtr(photoID),
		},
	})
	if err != nil {
		return nil, err
	}

	// 按關係類型分類事件ID
	events := map[string][]string{
		"reference": {},
		"polaroid":  {},
	}

	for _, tuple := range resp.RelationTuples {
		switch tuple.Relation {
		case "reference", "polaroid":
			if subject, ok := tuple.Subject.Ref.(*rts.Subject_Id); ok {
				events[tuple.Relation] = append(events[tuple.Relation], subject.Id)
			}
		}
	}

	return events, nil
}

// DeletePhotoEventRelation 刪除照片和事件之間的關係
//
// 參數:
//   - photoID: 照片的唯一標識符
//   - eventID: 事件的唯一標識符
//   - relationType: 關係類型 ("reference" 或 "polaroid")
//
// 返回:
//   - error: 如操作失敗則返回錯誤
func (k *Client) DeletePhotoEventRelation(photoID, eventID, relationType string) error {
	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Photo",
					Object:    photoID,
					Relation:  relationType, // "reference" 或 "polaroid"
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: eventID,
						},
					},
				},
			},
		},
	})
	return err
}

// strPtr 輔助函數，用於返回字符串的指針
func strPtr(s string) *string {
	return &s
}

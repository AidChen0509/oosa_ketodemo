package auth

import (
	"context"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// KetoClient 封裝了 Keto 相關操作
type KetoClient struct {
	writeClient rts.WriteServiceClient
	readClient  rts.ReadServiceClient
	writeConn   *grpc.ClientConn
	readConn    *grpc.ClientConn
}

// NewKetoClient 創建一個新的 Keto 客戶端
func NewKetoClient(writeAddress, readAddress string) (*KetoClient, error) {
	writeConn, err := grpc.NewClient(writeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	readConn, err := grpc.NewClient(readAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &KetoClient{
		writeClient: rts.NewWriteServiceClient(writeConn),
		readClient:  rts.NewReadServiceClient(readConn),
		writeConn:   writeConn,
		readConn:    readConn,
	}, nil
}

// Close 關閉連接
func (k *KetoClient) Close() {
	k.writeConn.Close()
	k.readConn.Close()
}

// CreateFriendRelation 建立朋友關係
func (k *KetoClient) CreateFriendRelation(userID1, userID2 string) error {
	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "User",
					Object:    userID1,
					Relation:  "friend",
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: userID2,
						},
					},
				},
			},
			// 朋友關係是雙向的，所以也添加反向關係
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "User",
					Object:    userID2,
					Relation:  "friend",
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: userID1,
						},
					},
				},
			},
		},
	})
	return err
}

// CreatePhotoViewPermission 創建照片查看權限
func (k *KetoClient) CreatePhotoViewPermission(photoID, userID string) error {
	_, err := k.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Photo",
					Object:    photoID,
					Relation:  "view",
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: userID,
						},
					},
				},
			},
		},
	})
	return err
}

// CheckPermission 使用關係查詢來檢查權限
func (k *KetoClient) CheckPermission(namespace, object, relation, subject string) (bool, error) {
	// 使用 ListRelationTuples 來查詢關係
	resp, err := k.readClient.ListRelationTuples(context.Background(), &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: &namespace,
			Object:    &object,
			Relation:  &relation,
			Subject: &rts.Subject{
				Ref: &rts.Subject_Id{
					Id: subject,
				},
			},
		},
	})
	if err != nil {
		return false, err
	}

	// 如果找到關係，則允許訪問
	return len(resp.RelationTuples) > 0, nil
}

// BatchCreateFriendRelations 批量建立朋友關係
func (k *KetoClient) BatchCreateFriendRelations(relationships []FriendRelationship) error {
	deltas := make([]*rts.RelationTupleDelta, 0, len(relationships)*2)

	for _, rel := range relationships {
		// 添加正向關係
		deltas = append(deltas, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "User",
				Object:    rel.User1ID,
				Relation:  "friend",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: rel.User2ID,
					},
				},
			},
		})

		// 添加反向關係
		deltas = append(deltas, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "User",
				Object:    rel.User2ID,
				Relation:  "friend",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: rel.User1ID,
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

// BatchCreatePhotoViewPermissions 批量設置照片查看權限
func (k *KetoClient) BatchCreatePhotoViewPermissions(permissions []PhotoPermission) error {
	deltas := make([]*rts.RelationTupleDelta, 0, len(permissions))

	for _, perm := range permissions {
		deltas = append(deltas, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Photo",
				Object:    perm.PhotoID,
				Relation:  "view",
				Subject: &rts.Subject{
					Ref: &rts.Subject_Id{
						Id: perm.UserID,
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

// FriendRelationship 朋友關係數據結構
type FriendRelationship struct {
	User1ID string
	User2ID string
}

// PhotoPermission 照片權限數據結構
type PhotoPermission struct {
	PhotoID string
	UserID  string
}

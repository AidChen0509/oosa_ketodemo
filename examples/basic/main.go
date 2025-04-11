package main

import (
	"fmt"
	"log"

	"github.com/AidChen0509/oosa_ketosdk/keto"
)

func main() {
	// 創建 Keto 客戶端
	// 注意：需要替換為您的 Keto 實例地址
	ketoClient, err := keto.NewClient("127.0.0.1:4467", "127.0.0.1:4466")
	if err != nil {
		log.Fatalf("初始化 Keto 客戶端失敗: %v", err)
	}
	defer ketoClient.Close()

	// 示例 1: 創建照片和事件之間的關係
	photoID := "photo123"
	eventID := "event456"

	err = ketoClient.CreatePhotoEventReference(photoID, eventID)
	if err != nil {
		log.Fatalf("建立 reference 關係失敗: %v", err)
	}
	fmt.Println("成功建立照片和事件的 reference 關係")

	// 示例 2: 檢查權限
	allowed, err := ketoClient.CheckPermission("Photo", photoID, "reference", eventID)
	if err != nil {
		log.Fatalf("檢查權限失敗: %v", err)
	}
	fmt.Printf("照片 %s 和事件 %s 的 reference 關係存在: %v\n", photoID, eventID, allowed)

	// 示例 3: 批量操作
	relations := []keto.PhotoEventRelation{
		{PhotoID: "photo123", EventID: "event789"},
		{PhotoID: "photo456", EventID: "event789"},
	}
	err = ketoClient.BatchCreatePhotoEventReferences(relations)
	if err != nil {
		log.Fatalf("批量建立關係失敗: %v", err)
	}
	fmt.Println("成功批量建立照片和事件的關係")

	// 示例 4: 獲取事件關聯的照片
	photos, err := ketoClient.GetEventReferencePhotos(eventID)
	if err != nil {
		log.Fatalf("獲取照片失敗: %v", err)
	}
	fmt.Printf("事件 %s 關聯的照片: %v\n", eventID, photos)

	// 示例 5: 刪除照片和事件之間的關係
	err = ketoClient.DeletePhotoEventRelation(photoID, eventID, "reference")
	if err != nil {
		log.Fatalf("刪除關係失敗: %v", err)
	}
	fmt.Println("成功刪除照片和事件的關係")
}

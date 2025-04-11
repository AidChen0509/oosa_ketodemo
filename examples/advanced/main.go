package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/AidChen0509/oosa_ketosdk/keto"
)

// 模擬完整的照片應用程序場景
func main() {
	// 初始化 Keto 客戶端
	ketoClient, err := keto.NewClient("127.0.0.1:4467", "127.0.0.1:4466")
	if err != nil {
		log.Fatalf("初始化 Keto 客戶端失敗: %v", err)
	}
	defer ketoClient.Close()

	// 模擬多個事件
	events := []string{"event1", "event2", "event3"}

	// 模擬多個照片
	photos := []string{"photo1", "photo2", "photo3", "photo4", "photo5"}

	// 建立不同類型的關係
	fmt.Println("建立照片與事件的關係...")

	// 為 event1 添加 reference 照片
	var relations []keto.PhotoEventRelation
	for i := 0; i < 3; i++ {
		relations = append(relations, keto.PhotoEventRelation{
			PhotoID: photos[i],
			EventID: events[0],
		})
	}

	if err := ketoClient.BatchCreatePhotoEventReferences(relations); err != nil {
		log.Fatalf("批量建立 reference 關係失敗: %v", err)
	}

	// 為 event2 添加 polaroid 照片
	polaroidRelations := []keto.PhotoEventRelation{
		{PhotoID: photos[3], EventID: events[1]},
		{PhotoID: photos[4], EventID: events[1]},
	}

	if err := ketoClient.BatchCreatePhotoEventPolaroids(polaroidRelations); err != nil {
		log.Fatalf("批量建立 polaroid 關係失敗: %v", err)
	}

	// 跨事件的照片關係 - 一張照片屬於多個事件
	if err := ketoClient.CreatePhotoEventReference(photos[2], events[2]); err != nil {
		log.Fatalf("建立跨事件照片關係失敗: %v", err)
	}

	// 展示權限檢查的使用
	fmt.Println("\n檢查權限...")
	checkAndPrintPermission(ketoClient, "Photo", photos[0], "reference", events[0])
	checkAndPrintPermission(ketoClient, "Photo", photos[3], "reference", events[0]) // 應該為 false
	checkAndPrintPermission(ketoClient, "Photo", photos[3], "polaroid", events[1])

	// 查詢關係
	fmt.Println("\n查詢關係...")

	// 同時查詢多個事件的照片 (並發示例)
	var wg sync.WaitGroup
	for _, event := range events {
		wg.Add(1)
		go func(eventID string) {
			defer wg.Done()
			fetchAndPrintPhotos(ketoClient, eventID)
		}(event)
	}
	wg.Wait()

	// 查詢特定照片的所有事件關係
	fetchAndPrintEvents(ketoClient, photos[2]) // 這張照片屬於多個事件

	// 清理某些關係
	fmt.Println("\n清理關係...")
	if err := ketoClient.DeletePhotoEventRelation(photos[0], events[0], "reference"); err != nil {
		log.Printf("刪除關係失敗: %v", err)
	} else {
		fmt.Printf("成功刪除照片 %s 與事件 %s 的 reference 關係\n", photos[0], events[0])
	}

	// 再次查詢確認關係已刪除
	checkAndPrintPermission(ketoClient, "Photo", photos[0], "reference", events[0]) // 應該為 false
}

// 輔助函數: 檢查並打印權限
func checkAndPrintPermission(client *keto.Client, namespace, object, relation, subject string) {
	allowed, err := client.CheckPermission(namespace, object, relation, subject)
	if err != nil {
		log.Printf("檢查權限失敗 (%s %s %s %s): %v", namespace, object, relation, subject, err)
		return
	}
	fmt.Printf("權限檢查: %s:%s-%s->%s = %v\n", namespace, object, relation, subject, allowed)
}

// 輔助函數: 獲取並打印照片
func fetchAndPrintPhotos(client *keto.Client, eventID string) {
	// 獲取 reference 照片
	refPhotos, err := client.GetEventReferencePhotos(eventID)
	if err != nil {
		log.Printf("獲取事件 %s 的 reference 照片失敗: %v", eventID, err)
	} else {
		fmt.Printf("事件 %s 的 reference 照片: %v\n", eventID, refPhotos)
	}

	// 獲取 polaroid 照片
	polPhotos, err := client.GetEventPolaroidPhotos(eventID)
	if err != nil {
		log.Printf("獲取事件 %s 的 polaroid 照片失敗: %v", eventID, err)
	} else {
		fmt.Printf("事件 %s 的 polaroid 照片: %v\n", eventID, polPhotos)
	}
}

// 輔助函數: 獲取並打印事件
func fetchAndPrintEvents(client *keto.Client, photoID string) {
	events, err := client.GetPhotoEvents(photoID)
	if err != nil {
		log.Printf("獲取照片 %s 的事件失敗: %v", photoID, err)
		return
	}

	fmt.Printf("照片 %s 的事件關係:\n", photoID)
	for relationType, eventIDs := range events {
		fmt.Printf("  %s: %v\n", relationType, eventIDs)
	}
}

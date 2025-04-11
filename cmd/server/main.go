package main

import (
	"log"

	"github.com/AidChen0509/oosa_ketosdk/api"
	"github.com/AidChen0509/oosa_ketosdk/keto"
)

func main() {
	// 初始化 Keto 客戶端
	ketoClient, err := keto.NewClient("127.0.0.1:4467", "127.0.0.1:4466")
	if err != nil {
		log.Fatalf("無法連接到 Keto 服務: %v", err)
	}
	defer ketoClient.Close()

	// 這裡開始添加你的 API 服務
	log.Println("事件照片管理後端啟動...")

	// 創建並啟動 API 服務器
	server := api.NewServer(ketoClient)
	log.Fatal(server.Run(":8080"))
}

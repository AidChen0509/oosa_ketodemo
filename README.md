# Ory Keto SDK for Go

這是一個基於 [Ory Keto](https://www.ory.sh/keto/) 的 Go SDK，專門用於管理照片和事件之間的權限關係。此 SDK 封裝了 Keto 的 gRPC API，提供簡單易用的介面來管理關係元組（Relation Tuples）。

## 功能特點

- 建立和管理照片與事件之間的 `reference` 和 `polaroid` 關係
- 批量處理關係操作
- 查詢特定照片或事件的關係
- 權限檢查和驗證
- 並發安全的 API 操作

## 安裝

使用 Go 模組引入此 SDK：

```bash
go get github.com/AidChen0509/oosa_ketosdk
```

## 前提條件

- Go 1.19 或更高版本
- 運行中的 Ory Keto 實例（版本 0.10.0 或更高）

## 快速開始

### 初始化客戶端

```go
import (
    "github.com/AidChen0509/oosa_ketosdk/keto"
)

// 創建 Keto 客戶端
ketoClient, err := keto.NewClient("127.0.0.1:4467", "127.0.0.1:4466")
if err != nil {
    // 處理錯誤
}
defer ketoClient.Close() // 記得釋放資源
```

### 建立關係

```go
// 建立照片和事件之間的 reference 關係
err := ketoClient.CreatePhotoEventReference("photo123", "event456")
if err != nil {
    // 處理錯誤
}

// 建立照片和事件之間的 polaroid 關係
err = ketoClient.CreatePhotoEventPolaroid("photo123", "event456")
if err != nil {
    // 處理錯誤
}
```

### 批量操作

```go
// 批量建立照片與事件的 reference 關係
relations := []keto.PhotoEventRelation{
    {PhotoID: "photo1", EventID: "event1"},
    {PhotoID: "photo2", EventID: "event1"},
    {PhotoID: "photo3", EventID: "event1"},
}

err := ketoClient.BatchCreatePhotoEventReferences(relations)
if err != nil {
    // 處理錯誤
}
```

### 查詢關係

```go
// 獲取與特定事件有 reference 關係的所有照片
photos, err := ketoClient.GetEventReferencePhotos("event1")
if err != nil {
    // 處理錯誤
}

// 獲取與特定照片有關係的所有事件
events, err := ketoClient.GetPhotoEvents("photo1")
if err != nil {
    // 處理錯誤
}
```

### 檢查權限

```go
// 檢查照片和事件之間是否存在特定關係
allowed, err := ketoClient.CheckPermission("Photo", "photo1", "reference", "event1")
if err != nil {
    // 處理錯誤
}

if allowed {
    // 存在關係
} else {
    // 不存在關係
}
```

### 刪除關係

```go
// 刪除照片和事件之間的關係
err := ketoClient.DeletePhotoEventRelation("photo1", "event1", "reference")
if err != nil {
    // 處理錯誤
}
```

## 詳細示例

詳細的使用示例可以在 `examples` 目錄下找到：

- [基本使用](examples/basic/main.go) - 基本 API 使用示例
- [進階示例](examples/advanced/main.go) - 包含並發和複雜關係管理的進階示例

## 如何整合到您的項目

1. 首先安裝 SDK
2. 在您的處理函數中初始化 KetoClient
3. 使用 SDK 提供的方法管理權限和關係

例如，在 Gin 框架中的使用示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/AidChen0509/oosa_ketosdk/keto"
)

func main() {
    r := gin.Default()
    
    // 創建 Keto 客戶端
    ketoClient, err := keto.NewClient("127.0.0.1:4467", "127.0.0.1:4466")
    if err != nil {
        panic(err)
    }
    defer ketoClient.Close()
    
    // 註冊路由
    r.POST("/photos/:photoID/events/:eventID/reference", func(c *gin.Context) {
        photoID := c.Param("photoID")
        eventID := c.Param("eventID")
        
        err := ketoClient.CreatePhotoEventReference(photoID, eventID)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, gin.H{"status": "success"})
    })
    
    r.GET("/events/:eventID/photos", func(c *gin.Context) {
        eventID := c.Param("eventID")
        
        photos, err := ketoClient.GetEventReferencePhotos(eventID)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, gin.H{"photos": photos})
    })
    
    r.Run(":8080")
}
```

## 貢獻

歡迎提交 Issues 和 Pull Requests。

## 授權

此項目採用 MIT 授權。詳見 [LICENSE](LICENSE) 文件。 
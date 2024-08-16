# Andromeda-API (簡易影音串流專案)

## 專案介紹

本專案是以 golang 開發的串流影音 Web API 專案(無介面設計)，因為是第一個 golang 專案，所以程式無採用 gin+gorm 框架套件，使用 go 原生 net/http 進行開發，主要目的是想進一步了解 golang 的語法與精進它的熟悉度。以下會有使用介紹與 DEMO 的連結，歡迎使用看看。(但 DEMO 網站目前架在 GCP 的 VM 上，三個月送的 300 美用完就會關閉 😆)

## 技術相關

- Go 1.22.5
- MySQL 8.0
- GCP
    - Storage (影片上傳)
    - Transcode API (影片轉檔)
- Basic Auth (後台 API 簡易授權)

### 套件

```
github.com/go-sql-driver/mysql                 (連結 MySQL 資料庫)
github.com/joho/godotenv                       (使用 .env 檔設定環境變數)
github.com/go-playground/validator/v10         (驗證器)
github.com/go-playground/universal-translator  (驗證的翻譯器)
github.com/google/uuid                         (影片唯一名稱)
cloud.google.com/go/storage                    (GCP Storage)
cloud.google.com/go/video                      (GCP Transcoder API)
```

## 目錄架構

```
Project
   ├─── common                                     # 通用工具
   │      │
   │      ├─── vod ─── vod.go                      # 影音上傳與轉檔
   │      │
   │      ├─── data_type.go                        # 自訂類型
   │      │
   │      └─── utils.go                            # 工具函式
   │
   ├─── config ─── config.go                       # 設定
   │
   ├─── controllers                                # 相關 controller
   │      │
   │      ├─── admin                               # 後台相關的 controller
   │      │      │
   │      │      ├─── video_controller.go          # 影片資料的 controller
   │      │      │
   │      │      └─── video_file_controller.go     # 影片檔案的 controller
   │      │
   │      ├─── validator                           # 相關驗證器
   │      │      ├─── rules
   │      │      │      ├─── video.go              # video_controller.go 的驗證規則
   │      │      │      │
   │      │      │      └─── video_file.go         # video_file_controller.go 的驗證規則
   │      │      │
   │      │      └─── validator.go                 # 驗證器
   │      │
   │      └─── home_controller.go                  # 前台首頁 controller
   │ 
   ├─── middleware                                 # 相關 middleware
   │      │
   │      ├─── basic_auth_middleware.go            # 後台 API 使用的 Basic Auth 
   │      │
   │      └─── middlewares.go                      # middleware 的入口，內容是呼叫所有 middleware 與 API 的錯誤處理
   │      
   ├─── models                                     # 相關 model
   │      │
   │      ├─── base.go                             # 資料庫的設定
   │      │
   │      ├─── video.go                            # 影片資料的 model
   │      │
   │      └─── video_file.go                       # 影片檔案的 model
   │
   ├─── routes ─── routes.go                       # API 路由
   │
   ├─── services                                   # 商業邏輯
   │      │
   │      ├─── service_home.go                     # 前台首頁的 API 邏輯
   │      │
   │      ├─── service_video.go                    # 後台影片資料的 API 邏輯
   │      │
   │      └─── service_video_file.go               # 後台影片檔案的 API 邏輯
   │      
   ├─── .env.example                               # 環境設定參考檔
   │
   └─── main.go                                    # 程式入口

```

## 本地使用

1. clone Repository
2. env 相關
    1. 複製 .env.example 命名為 .env
    2. 設定 Server 連線資訊 (Port or TLS 憑證與私鑰的路徑)
    3. 設定 MySQL 資料庫的資訊
    4. .設定 GCP 的 json 憑證路徑與 Storage、Transcode API 相關設定
    5. 設定 Basic Auth 帳密
    6. 若要上正式環境，可把 .env 的 APP_ENV 設定成 production，避免程式錯誤訊息外漏
3. 將 database/mysql/andromeda.sql 匯入 MySql 8.0 資料庫
4. 執行 go run main.go
5. 將 docs/postman/collection.json 匯入 Postman
    1. 把 Andromeda 目錄設定 Variables 的 domain 為你本地的 URL ( ex: http://localhost:8080 )
    2. 將 後台 API 目錄設定 Authorization 為 BasicAuth
6. 即可使用

## 線上 DEMO

[前台 API Demo 頁](https://vod.olivermg.fun/api/home) (僅顯示影片狀態設定 publish + 檔案轉檔完成的資料)

### POSTMAN 後台 API 操作流程

1. 後台 API / 影片 / 新增
2. 後台 API / 影片 / 列表
3. 後台 API / 影片 / 單筆 (POSTMAN 會紀錄 video_id)
4. 後台 API / 檔案 / 上傳與轉檔
5. 後台 API / 檔案 / 單筆 (前端可每秒打一次 check 裝態，transformed 代表轉檔完成)
6. 後台 API / 影片 / 更新
7. 前台 API / 首頁顯示所有公開影片
     1. 複製 hls_path or mpd_path 的連結 (前台 API 會顯示 影片設定 publish + 檔案轉檔完成的資料)
8. 到 https://livepush.io/hls-player/index.html (可以使用這個網頁播放器來測試，如果複製 hls_path 用 HLS PLAYER 播放，如果複製 mpd_path 用 DASH PALYER 播放)

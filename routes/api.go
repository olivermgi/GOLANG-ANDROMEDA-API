package routes

import (
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	controllers "github.com/olivermgi/golang-crud-practice/controllers"
	admin_controllers "github.com/olivermgi/golang-crud-practice/controllers/admin"
)

// 註冊 API 路由
func RegisterAPIRoutes() {
	http.HandleFunc("/", notFoundHandler)

	// 後台 API
	http.HandleFunc("POST /api/internal/videos", admin_controllers.StoreVideo)
	http.HandleFunc("GET /api/internal/videos", admin_controllers.IndexVideo)
	http.HandleFunc("GET /api/internal/videos/{video_id}", admin_controllers.ShowVideo)
	http.HandleFunc("PUT /api/internal/videos/{video_id}", admin_controllers.UpdateVideo)
	http.HandleFunc("DELETE /api/internal/videos/{video_id}", admin_controllers.DestroyVideo)

	http.HandleFunc("POST /api/internal/videos/{video_id}/files", admin_controllers.StoreVideoFile)
	http.HandleFunc("GET /api/internal/videos/{video_id}/files", admin_controllers.ShowVideoFile)
	http.HandleFunc("DELETE /api/internal/videos/{video_id}/files", admin_controllers.DestroyVideoFile)

	// 前台 API
	http.HandleFunc("GET /api/home", controllers.Home)
}

// 404 Response
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	common.Abort(http.StatusNotFound, "此 API 不存在")
}

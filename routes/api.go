package routes

import (
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/controllers"
)

// 註冊 API 路由
func RegisterAPIRoutes() {
	http.HandleFunc("/", notFoundHandler)

	http.HandleFunc("POST /api/internal/videos", controllers.StoreVideo)
	http.HandleFunc("GET /api/internal/videos", controllers.IndexVideo)
	http.HandleFunc("GET /api/internal/videos/{video_id}", controllers.ShowVideo)
	http.HandleFunc("PUT /api/internal/videos/{video_id}", controllers.UpdateVideo)
	http.HandleFunc("DELETE /api/internal/videos/{video_id}", controllers.DestroyVideo)

	// http.HandleFunc("GET /api/companies", controllers.IndexCompany)
	// http.HandleFunc("POST /api/companies", controllers.StoreCompany)
	// http.HandleFunc("GET /api/companies/{company_id}", controllers.ShowCompany)
	// http.HandleFunc("PUT /api/companies/{company_id}", controllers.UpdateCompany)
	// http.HandleFunc("DELETE /api/companies/{company_id}", controllers.DestroyCompany)
}

// 404 Response
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	common.Abort(http.StatusNotFound, "此 API 不存在", nil)
}

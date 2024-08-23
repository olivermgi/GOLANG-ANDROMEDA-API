package controllers

import (
	"net/http"

	"github.com/olivermgi/golang-andromeda-api/common"
	"github.com/olivermgi/golang-andromeda-api/services"
)

// 首頁
func Home(w http.ResponseWriter, r *http.Request) {
	service := &services.ServiceHome{}
	videos := service.Home()

	common.Response(http.StatusOK, "首頁資料取得成功", videos, w)
}

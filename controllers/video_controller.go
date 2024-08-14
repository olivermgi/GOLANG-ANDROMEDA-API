package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/controllers/validator"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video"
	"github.com/olivermgi/golang-crud-practice/services"
)

// 顯示影片列表 (分頁+排序)
func IndexVideo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ruleData := &rules.VideoIndex{
		Page:       common.StringToInt(r.FormValue("page")),
		PerPage:    common.StringToInt(r.FormValue("per_page")),
		Sort:       r.FormValue("sort"),
		SortColumn: r.FormValue("sort_column"),
	}

	validator.ValidateOrAbort(ruleData)

	service := &services.ServiceVideo{}
	paginations := service.Index(ruleData)

	common.Response(http.StatusOK, "影片資料列表取得成功", paginations, w)
}

// 新增影片資料
func StoreVideo(w http.ResponseWriter, r *http.Request) {
	var ruleData *rules.VideoStore
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ruleData)
	if err != nil {
		common.Abort(http.StatusForbidden, "資料格式不正確")
	}

	validator.ValidateOrAbort(ruleData)

	service := &services.ServiceVideo{}
	video := service.Store(ruleData)

	common.Response(http.StatusCreated, "影片資料新增成功", video, w)
}

// 顯示單筆影片資料
func ShowVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "video_id 路徑參數不正確")
	}

	ruleData := &rules.VideoShow{VideoId: videoId}
	validator.ValidateOrAbort(ruleData)

	service := &services.ServiceVideo{}
	video := service.GetOrAbort(ruleData.VideoId)

	common.Response(http.StatusOK, "單筆影片資料取得成功", video, w)
}

// 更新單筆公司資料
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "video_id 路徑參數不正確")
	}

	var ruleData *rules.VideoUpdate
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&ruleData)
	if err != nil {
		common.Abort(http.StatusForbidden, "資料格式不正確")
	}

	ruleData.VideoId = videoId
	validator.ValidateOrAbort(ruleData)

	service := &services.ServiceVideo{}
	service.GetOrAbort(ruleData.VideoId)
	video := service.Update(ruleData)

	common.Response(http.StatusOK, "影片資料更新成功", video, w)
}

// 刪除一筆影片資料
func DestroyVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "video_id 路徑參數不正確")
	}

	ruleData := &rules.VideoDelete{VideoId: videoId}
	validator.ValidateOrAbort(ruleData)

	service := &services.ServiceVideo{}
	service.GetOrAbort(ruleData.VideoId)
	service.Delete(ruleData.VideoId)

	common.Response(http.StatusOK, "影片資料刪除成功", nil, w)
}

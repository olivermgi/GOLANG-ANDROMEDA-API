package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/services"
)

// 顯示影片列表 (分頁+排序)
func IndexVideo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	videoListData := services.VideoListData{
		Page:       common.StringToInt(r.FormValue("page")),
		PerPage:    common.StringToInt(r.FormValue("per_page")),
		Sort:       r.FormValue("sort"),
		SortColumn: r.FormValue("sort_column"),
	}

	common.ValidateStruct(videoListData)

	paginations := services.IndexVideo(videoListData)

	common.Response(http.StatusOK, "影片資料列表取得成功", paginations, w)
}

// 新增影片資料
func StoreVideo(w http.ResponseWriter, r *http.Request) {
	var videoData services.VideoData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&videoData)
	if err != nil {
		common.Abort(http.StatusForbidden, "JSON 格式不正確", nil)
	}

	common.ValidateStruct(videoData)

	video := services.StoreVideo(videoData)

	common.Response(http.StatusCreated, "影片資料新增成功", video, w)
}

// 顯示單筆影片資料
func ShowVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "參數類型錯誤", nil)
	}

	common.ValidateStruct(struct {
		VideoId int `validate:"required,min=1" field:"video_id "`
	}{videoId})

	video := services.GetVideo(videoId)

	common.Response(http.StatusOK, "單筆影片資料取得成功", video, w)
}

// 更新單筆公司資料
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "參數類型錯誤", nil)
	}

	var videoData services.VideoData
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&videoData)
	if err != nil {
		common.Abort(http.StatusForbidden, "JSON 格式不正確", nil)
	}

	common.ValidateStruct(videoData)

	video := services.UpdateVideo(videoId, videoData)

	common.Response(http.StatusOK, "影片資料更新成功", video, w)
}

// 刪除一筆公司資料
func DestroyVideo(w http.ResponseWriter, r *http.Request) {
	videoId, err := strconv.Atoi(r.PathValue("video_id"))
	if err != nil {
		common.Abort(http.StatusForbidden, "參數類型錯誤", nil)
	}

	services.DeleteVideo(videoId)

	common.Response(http.StatusOK, "影片資料刪除成功", nil, w)
}

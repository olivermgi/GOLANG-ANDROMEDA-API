package services

import (
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/models"
)

type VideoData struct {
	Status      string `validate:"required,oneof=publish unpublish" field:"狀態"`
	Title       string `validate:"required,max=50" field:"標題"`
	Description string `validate:"omitempty,max=255" field:"描述"`
}

type VideoListData struct {
	Page       int    `validate:"required,min=1" field:"頁數"`
	PerPage    int    `validate:"required,oneof=10 20 30" field:"筆數"`
	Sort       string `validate:"required,oneof=asc desc" field:"排序方式"`
	SortColumn string `validate:"required,oneof=id updated_at" field:"排序欄位"`
}

var model *models.Video

func IndexVideo(videListData VideoListData) map[string]interface{} {
	videos, total, last_page := model.Paginate(videListData.Page, videListData.PerPage,
		videListData.SortColumn, videListData.Sort)

	return map[string]interface{}{
		"page":      videListData.Page,
		"per_page":  videListData.PerPage,
		"total":     total,
		"last_page": last_page,
		"items":     videos,
	}
}

func StoreVideo(videoData VideoData) *models.Video {
	data := models.Video{
		Status:      videoData.Status,
		Title:       videoData.Title,
		Description: videoData.Description,
	}

	video := model.Insert(data)

	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料新增失敗", nil)
	}

	return video
}

func GetVideo(video_id int) *models.Video {
	return model.Get(video_id)
}

func GetVideoOrAbort(video_id int) *models.Video {
	video := GetVideo(video_id)

	if video == nil {
		common.Abort(http.StatusNotFound, "無此影片資料資料", nil)
	}

	return video
}

func UpdateVideo(videoId int, videoData VideoData) *models.Video {
	GetVideoOrAbort(videoId)

	data := models.Video{
		Status:      videoData.Status,
		Title:       videoData.Title,
		Description: videoData.Description,
	}

	video := model.Update(videoId, data)
	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料更新失敗", nil)
	}

	return video
}

func DeleteVideo(videoId int) {
	GetVideoOrAbort(videoId)

	is_success := model.Delete(videoId)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗", nil)
	}
}

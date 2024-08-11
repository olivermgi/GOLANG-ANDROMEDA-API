package services

import (
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video"
	"github.com/olivermgi/golang-crud-practice/models"
)

var model *models.Video

func IndexVideo(passedData *rules.VideoIndex) map[string]interface{} {
	videos, total, last_page := model.Paginate(passedData.Page,
		passedData.PerPage, passedData.SortColumn, passedData.Sort)

	return map[string]interface{}{
		"page":      passedData.Page,
		"per_page":  passedData.PerPage,
		"total":     total,
		"last_page": last_page,
		"items":     videos,
	}
}

func StoreVideo(passedData *rules.VideoStore) *models.Video {
	dbData := models.Video{
		Status:      passedData.Status,
		Title:       passedData.Title,
		Description: passedData.Description,
	}

	video := model.Insert(dbData)

	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料新增失敗")
	}

	return video
}

func GetVideo(videoId int) *models.Video {
	return model.Get(videoId)
}

func GetVideoOrAbort(videoId int) *models.Video {
	video := GetVideo(videoId)

	if video == nil {
		common.Abort(http.StatusNotFound, "無此影片資料資料")
	}

	return video
}

func UpdateVideo(passedData *rules.VideoUpdate) *models.Video {
	GetVideoOrAbort(passedData.VideoId)

	data := models.Video{
		Status:      passedData.Status,
		Title:       passedData.Title,
		Description: passedData.Description,
	}

	video := model.Update(passedData.VideoId, data)
	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料更新失敗")
	}

	return video
}

func DeleteVideo(videoId int) {
	GetVideoOrAbort(videoId)

	is_success := model.Delete(videoId)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗")
	}
}

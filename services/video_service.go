package services

import (
	"github.com/olivermgi/golang-crud-practice/models"
)

type VideoData struct {
	Status      string `validate:"required,oneof=publish unpublish" validate_field:"狀態"`
	Title       string `validate:"required,max=50" validate_field:"標題"`
	Description string `validate:"omitempty,max=255" validate_field:"描述"`
}

var model *models.Video

func StoreVideo(videoData VideoData) *models.Video {
	data := models.Video{
		Status:      videoData.Status,
		Title:       videoData.Title,
		Description: videoData.Description,
	}

	video := model.Insert(data)

	return video
}

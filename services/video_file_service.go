package services

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/olivermgi/golang-crud-practice/common"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video_file"
	"github.com/olivermgi/golang-crud-practice/models"
)

var videoFileModel *models.VideoFile

func StoreVideoFile(passedData *rules.VideoFileStore) *models.VideoFile {
	if GetVideoFile(passedData.VideoId) != nil {
		common.Abort(http.StatusForbidden, "影片檔案資料已存在")
	}

	// 產生影片檔名
	var name string
	for {
		name = uuid.New().String()[:8]
		if videoFileModel.GetByName(name) == nil {
			break
		}
		log.Println("影片檔案名稱生成重複，名稱：", name)
	}

	dbData := models.VideoFile{
		VideoId: passedData.VideoId,
		Name:    name,
	}

	videoFile := videoFileModel.Insert(dbData)

	if videoFile == nil {
		common.Abort(http.StatusForbidden, "影片檔案資料新增失敗")
	}

	// 處理檔案

	return videoFile
}

func GetVideoFile(videoId int) *models.VideoFile {
	return videoFileModel.GetByVideoId(videoId)
}

func GetVideoFileOrAbort(videoId int) *models.VideoFile {
	videoFile := GetVideoFile(videoId)

	if videoFile == nil {
		common.Abort(http.StatusNotFound, "無此影片檔案資料")
	}

	return videoFile
}

func DeleteVideoFile(videoId int) {
	GetVideoFileOrAbort(videoId)

	is_success := videoFileModel.SoftDelete(videoId)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗")
	}

	// 發送關閉正在轉檔的 Job 信號
}

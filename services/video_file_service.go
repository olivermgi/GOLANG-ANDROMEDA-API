package services

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/olivermgi/golang-crud-practice/common"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video_file"
	"github.com/olivermgi/golang-crud-practice/models"
)

type VideoFileService struct {
	model *models.VideoFile
}

func (s *VideoFileService) Store(passedData *rules.VideoFileStore) *models.VideoFile {
	if s.Get(passedData.VideoId) != nil {
		common.Abort(http.StatusForbidden, "影片檔案資料已存在")
	}

	// 產生影片檔名
	var name string
	for {
		name = uuid.New().String()[:8]
		if s.model.GetByName(name) == nil {
			break
		}
		log.Println("影片檔案名稱生成重複，名稱：", name)
	}

	dbData := models.VideoFile{
		VideoId: passedData.VideoId,
		Name:    name,
	}

	videoFile := s.model.Insert(dbData)

	if videoFile == nil {
		common.Abort(http.StatusForbidden, "影片檔案資料新增失敗")
	}

	// 處理檔案

	return videoFile
}

func (s *VideoFileService) Get(videoId int) *models.VideoFile {
	return s.model.GetByVideoId(videoId)
}

func (s *VideoFileService) GetOrAbort(videoId int) *models.VideoFile {
	videoFile := s.Get(videoId)

	if videoFile == nil {
		common.Abort(http.StatusNotFound, "無此影片檔案資料")
	}

	return videoFile
}

func (s *VideoFileService) Delete(videoId int) {
	s.GetOrAbort(videoId)

	is_success := s.model.SoftDelete(videoId)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗")
	}

	// 發送關閉正在轉檔的 Job 信號
}

package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/common/vod"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video_file"
	"github.com/olivermgi/golang-crud-practice/models"
)

type ServiceVideoFile struct {
	model *models.VideoFile
}

func (s *ServiceVideoFile) UploadAndTransformVideoFile(videoFile *models.VideoFile, file multipart.File) {
	if s.UploadVideoFile(videoFile, file) {
		s.TransformVideoFile(videoFile)
	} else {
		s.model.UpdateStatus(videoFile.VideoId, "upload_fail")
	}
}

func (s *ServiceVideoFile) TransformVideoFile(videoFile *models.VideoFile) bool {
	if videoFile == nil || videoFile.Status != "uploaded" {
		return false
	}

	filename := videoFile.Name
	extension := filepath.Ext(filename)

	inputPath := fmt.Sprintf("%s/%s", filename[:1], filename)
	outputPaht := fmt.Sprintf("%s/%s/", filename[:1], filename[0:len(filename)-len(extension)])

	if !s.model.UpdateStatus(videoFile.VideoId, "transforming") {
		return false
	}

	jobID, err := vod.Transcoder.TransformVideoFile(inputPath, outputPaht)
	if err != nil {
		return false
	}

	go func() {
		var count uint8
		var is_success = false
		for count < 120 {
			jobState, _ := vod.Transcoder.GetJobState(jobID)
			if jobState == "SUCCEEDED" {
				s.model.UpdateStatus(videoFile.VideoId, "transformed")
				is_success = true
				break
			}

			time.Sleep(time.Second)
			count++
		}

		if !is_success {
			s.model.UpdateStatus(videoFile.VideoId, "transform_failed")
		}

		vod.Transcoder.DeleteJob(jobID)
	}()

	videoFile.Status = "transforming"
	return true
}

func (s *ServiceVideoFile) UploadVideoFile(videoFile *models.VideoFile, file multipart.File) bool {
	filename := videoFile.Name
	path := fmt.Sprintf("mp4/%s/%s", filename[:1], filename)

	if !s.model.UpdateStatus(videoFile.VideoId, "uploading") {
		return false
	}

	if vod.Uploader.UploadFile(file, path) != nil {
		return false
	}

	if !s.model.UpdateStatus(videoFile.VideoId, "uploaded") {
		return false
	}

	videoFile.Status = "uploaded"
	return true
}

func (s *ServiceVideoFile) Store(passedData *rules.VideoFileStore) *models.VideoFile {
	if s.Get(passedData.VideoId) != nil {
		common.Abort(http.StatusForbidden, "影片檔案資料已存在")
	}

	// 產生影片檔名
	var name string
	for {
		name = uuid.New().String()[:8] + ".mp4"
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

	return videoFile
}

func (s *ServiceVideoFile) Get(videoId int) *models.VideoFile {
	return s.model.GetByVideoId(videoId)
}

func (s *ServiceVideoFile) GetOrAbort(videoId int) *models.VideoFile {
	videoFile := s.Get(videoId)

	if videoFile == nil {
		common.Abort(http.StatusNotFound, "無此影片檔案資料")
	}

	return videoFile
}

func (s *ServiceVideoFile) Delete(videoId int) {
	s.GetOrAbort(videoId)

	is_success := s.model.SoftDelete(videoId)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗")
	}

	// 發送關閉正在轉檔的 Job 信號
	// 將轉好的串流檔移至其他資料夾
}

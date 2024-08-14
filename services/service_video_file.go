package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"sync"
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
		common.Abort(http.StatusForbidden, "影片檔案已存在")
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
		common.Abort(http.StatusForbidden, "影片檔案新增失敗")
	}

	return videoFile
}

func (s *ServiceVideoFile) Get(videoId int) *models.VideoFile {
	return s.model.GetByVideoId(videoId)
}

func (s *ServiceVideoFile) GetOrAbort(videoId int) *models.VideoFile {
	videoFile := s.Get(videoId)

	if videoFile == nil {
		common.Abort(http.StatusNotFound, "無此影片檔案")
	}

	return videoFile
}

func (s *ServiceVideoFile) Delete(videoId int) {
	videoFile := s.GetOrAbort(videoId)

	if !slices.Contains([]string{"stanby", "uploaded", "transformed", "delete_failed"}, videoFile.Status) {
		common.Abort(http.StatusForbidden, "影片檔案資料正在處理，無法刪除檔案")
	}

	is_success := s.model.UpdateStatus(videoFile.VideoId, "deleting")
	if !is_success {
		common.Abort(http.StatusForbidden, "影片檔案刪除失敗")
	}

	videoFile.Status = "deleting"
	go s.hiddenVideo(videoFile)
}

func (s *ServiceVideoFile) hiddenVideo(videoFile *models.VideoFile) {
	var wg sync.WaitGroup

	filename := videoFile.Name
	extension := filepath.Ext(filename)
	filePath := fmt.Sprintf("%s/%s", filename[:1], filename)
	streamPath := fmt.Sprintf("%s/%s/", filename[:1], filename[0:len(filename)-len(extension)])
	mp4Dir := "mp4/"
	streamsDir := "streams/"
	deletedDir := "deleted/"

	wg.Add(2)
	deletedCount := 0
	go func() {
		defer wg.Done()
		vod.Uploader.MoveFile(mp4Dir+filePath, deletedDir+filePath)
		deletedCount++
	}()

	go func() {
		defer wg.Done()
		vod.Uploader.MoveFolder(streamsDir+streamPath, deletedDir+streamPath)
		deletedCount++
	}()
	wg.Wait()
	if deletedCount != 2 {
		s.model.UpdateStatus(videoFile.VideoId, "delete_failed")
		return
	}

	s.model.SoftDelete(videoFile.VideoId)
}

package services

import (
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	rules "github.com/olivermgi/golang-crud-practice/controllers/validator/rules/video"
	"github.com/olivermgi/golang-crud-practice/models"
)

type ServiceVideo struct {
	model *models.Video
}

func (s *ServiceVideo) Index(passedData *rules.VideoIndex) map[string]interface{} {
	videos, total, last_page := s.model.Paginate(passedData.Page,
		passedData.PerPage, passedData.SortColumn, passedData.Sort)

	return map[string]interface{}{
		"page":      passedData.Page,
		"per_page":  passedData.PerPage,
		"total":     total,
		"last_page": last_page,
		"items":     videos,
	}
}

func (s *ServiceVideo) Store(passedData *rules.VideoStore) *models.Video {
	dbData := models.Video{
		Status:      passedData.Status,
		Title:       passedData.Title,
		Description: passedData.Description,
	}

	video := s.model.Insert(dbData)

	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料新增失敗")
	}

	return video
}

func (s *ServiceVideo) Get(id int) *models.Video {
	return s.model.Get(id)
}

func (s *ServiceVideo) GetOrAbort(id int) *models.Video {
	video := s.Get(id)

	if video == nil {
		common.Abort(http.StatusNotFound, "無此影片資料")
	}

	return video
}

func (s *ServiceVideo) Update(passedData *rules.VideoUpdate) *models.Video {
	s.GetOrAbort(passedData.VideoId)

	data := models.Video{
		Status:      passedData.Status,
		Title:       passedData.Title,
		Description: passedData.Description,
	}

	video := s.model.Update(passedData.VideoId, data)
	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料更新失敗")
	}

	return video
}

func (s *ServiceVideo) Delete(id int) {
	s.GetOrAbort(id)

	is_success := s.model.SoftDelete(id)
	if !is_success {
		common.Abort(http.StatusForbidden, "影片資料刪除失敗")
	}
}

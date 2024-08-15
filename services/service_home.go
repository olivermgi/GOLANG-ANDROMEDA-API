package services

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/olivermgi/golang-crud-practice/models"
)

type ServiceHome struct {
	model *models.Video
}

func (s *ServiceHome) Home() []map[string]interface{} {
	videos := s.model.AllPublish()

	videoMaps := make([]map[string]interface{}, 0)
	jsonStr, _ := json.Marshal(videos)
	err := json.Unmarshal(jsonStr, &videoMaps)
	if err != nil {
		return make([]map[string]interface{}, 0)
	}

	for key, videoMap := range videoMaps {
		videoFile := videoMap["video_file"].(map[string]interface{})
		filename := videoFile["name"].(string)
		extension := filepath.Ext(filename)
		videoFile["hls_path"] = fmt.Sprintf("https://static.olivermg.fun/streams/%s/%s/manifest.m3u8", filename[:1], filename[0:len(filename)-len(extension)])
		videoFile["mpd_path"] = fmt.Sprintf("https://static.olivermg.fun/streams/%s/%s/manifest.mpd", filename[:1], filename[0:len(filename)-len(extension)])
		delete(videoFile, "name")
		videoMaps[key]["video_file"] = videoFile
	}

	return videoMaps
}

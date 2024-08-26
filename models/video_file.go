package models

import (
	"time"
)

type VideoFile struct {
	Id        int    `json:"-"`
	Status    string `json:"status,omitempty"`
	Name      string `json:"name"`
	VideoId   int    `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

// 新增公司單筆資料
func (c *VideoFile) Insert(data VideoFile) *VideoFile {
	var id int
	query := `INSERT INTO video_files(name, video_id, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	err := DB.QueryRow(query, data.Name, data.VideoId).Scan(&id)
	if err != nil {
		return nil
	}

	return &VideoFile{
		Id:      int(id),
		Status:  "standby",
		Name:    data.Name,
		VideoId: data.VideoId,
	}
}

// 以 name 抓取影片單筆檔案資料
func (c *VideoFile) GetByVideoId(videoId int) *VideoFile {
	var videoFile VideoFile
	err := DB.QueryRow("SELECT id, status, name, video_id FROM video_files WHERE video_id = $1 AND deleted_at IS NULL", videoId).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name, &videoFile.VideoId)

	if err != nil {
		return nil
	}

	return &videoFile
}

// 以 name 抓取影片單筆檔案資料
func (c *VideoFile) GetByName(name string) *VideoFile {
	var videoFile VideoFile
	err := DB.QueryRow("SELECT id, status, name FROM video_files WHERE name = $1 AND deleted_at IS NULL", name).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name)

	if err != nil {
		return nil
	}

	return &videoFile
}

func (c *VideoFile) UpdateStatus(videoId int, status string) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE video_files SET status = $1, updated_at = $2 WHERE video_id = $3 AND deleted_at IS NULL",
		status, now, videoId)

	return err == nil
}

func (c *VideoFile) SoftDelete(videoId int) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE video_files SET status = $1, deleted_at = $2 WHERE video_id = $3 AND deleted_at IS NULL",
		"deleted", now, videoId)

	return err == nil
}

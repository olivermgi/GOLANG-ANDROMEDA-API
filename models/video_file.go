package models

import (
	"time"
)

type VideoFile struct {
	Id        int    `json:"-"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	VideoId   int    `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

// 新增公司單筆資料
func (c *VideoFile) Insert(data VideoFile) *VideoFile {
	result, err := DB.Exec("INSERT INTO video_files(name, video_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		data.Name, data.VideoId)
	if err != nil {
		return nil
	}
	id, err := result.LastInsertId()
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
	err := DB.QueryRow("SELECT id, status, name, video_id FROM video_files WHERE video_id = ? AND deleted_at IS NULL", videoId).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name, &videoFile.VideoId)

	if err != nil {
		return nil
	}

	return &videoFile
}

// 以 name 抓取影片單筆檔案資料
func (c *VideoFile) GetByName(name string) *VideoFile {
	var videoFile VideoFile
	err := DB.QueryRow("SELECT id, status, name FROM video_files WHERE name = ? AND deleted_at IS NULL", name).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name)

	if err != nil {
		return nil
	}

	return &videoFile
}

func (c *VideoFile) UpdateStatus(videoId int, status string) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE video_files SET status = ?, updated_at = ? WHERE video_id = ? AND deleted_at IS NULL",
		status, now, videoId)

	return err == nil
}

func (c *VideoFile) SoftDelete(videoId int) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE video_files SET status = ?, deleted_at = ? WHERE video_id = ? AND deleted_at IS NULL",
		"deleted", now, videoId)

	return err == nil
}

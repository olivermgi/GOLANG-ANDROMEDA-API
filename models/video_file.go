package models

import (
	"log"
	"time"
)

type VideoFile struct {
	Id        int    `json:"id"`
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
		log.Println("新增影片檔案資料失敗，錯誤訊息：", err)
		return nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("取得新增的影片檔案 ID 失敗，錯誤訊息：", err)
		return nil
	}

	return &VideoFile{
		Id:     int(id),
		Status: "standby",
		Name:   data.Name,
	}
}

// 以 name 抓取影片單筆檔案資料
func (c *VideoFile) GetByVideoId(videoId int) *VideoFile {
	var videoFile VideoFile
	err := DB.QueryRow("SELECT id, status, name FROM video_files WHERE video_id = ? AND deleted_at IS NULL", videoId).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name)

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

// 以 video_id 軟刪除公司單筆資料
func (c *VideoFile) SoftDelete(videoId int) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE video_files SET deleted_at = ? WHERE video_id = ? AND deleted_at IS NULL",
		now, videoId)

	if err != nil {
		log.Println("軟刪除影片檔案資料失敗，錯誤訊息：", err)
		return false
	}

	return true
}

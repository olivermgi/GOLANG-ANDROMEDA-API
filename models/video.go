package models

import (
	"fmt"
	"log"
	"math"
	"time"
)

type Video struct {
	Id          int         `json:"id"`
	Status      string      `json:"status"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	VideoFile   interface{} `json:"video_file,omitempty"`
	CreatedAt   string      `json:"-"`
	UpdatedAt   string      `json:"-"`
	DeletedAt   string      `json:"-"`
}

// 新增影片單筆資料
func (c *Video) Insert(data Video) *Video {
	result, err := DB.Exec("INSERT INTO videos(status, title, description, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		data.Status, data.Title, data.Description)
	if err != nil {
		return nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}

	return &Video{
		Id:          int(id),
		Status:      data.Status,
		Title:       data.Title,
		Description: data.Description,
	}
}

// 抓取影片多筆資料，並回傳頁面格式
func (c *Video) Paginate(page int, perPage int, sortColume string, sort string) (videos []Video, total int, lastPage int) {
	queryStr := fmt.Sprintf(
		"SELECT id, status, title, updated_at FROM videos deleted_at IS NULL ORDER BY %s %s, id DESC LIMIT ? OFFSET ?",
		sortColume, sort,
	)

	page = (page - 1) * perPage
	rows, err := DB.Query(queryStr, perPage, page)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	videos = make([]Video, 0)
	for rows.Next() {
		var video Video
		err := rows.Scan(&video.Id, &video.Status, &video.Title, &video.UpdatedAt)
		if err != nil {
			return make([]Video, 0), 0, 0
		}
		videos = append(videos, video)
	}

	count := 0
	DB.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	lastPage = int(math.Ceil(float64(count) / float64(perPage)))
	return videos, count, lastPage
}

// 抓取影片全部資料
func (c *Video) All() []Video {
	queryStr := fmt.Sprintln("SELECT id, status, title, updated_at FROM videos deleted_at IS NULL")

	rows, err := DB.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	videos := make([]Video, 0)
	for rows.Next() {
		var video Video
		err := rows.Scan(&video.Id, &video.Status, &video.Title, &video.UpdatedAt)
		if err != nil {
			return make([]Video, 0)
		}
		videos = append(videos, video)
	}

	return videos
}

// 以 id 抓取影片單筆資料
func (c *Video) Get(id int) *Video {
	var video Video
	err := DB.QueryRow("SELECT id, status, title, description FROM videos WHERE id = ? AND deleted_at IS NULL", id).
		Scan(&video.Id, &video.Status, &video.Title, &video.Description)

	if err != nil {
		return nil
	}

	var videoFile VideoFile
	err = DB.QueryRow("SELECT id, status, name FROM video_files WHERE video_id = ? AND deleted_at IS NULL", id).
		Scan(&videoFile.Id, &videoFile.Status, &videoFile.Name)

	if err != nil {
		video.VideoFile = struct{}{}
	} else {
		video.VideoFile = videoFile
	}

	return &video
}

// 以 id 更新影片單筆資料
func (c *Video) Update(id int, data Video) *Video {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE videos SET status = ?, title = ?, description = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		data.Status, data.Title, data.Description, now, id)

	if err != nil {
		return nil
	}

	return &Video{
		Id:          id,
		Status:      data.Status,
		Title:       data.Title,
		Description: data.Description,
	}
}

// 以 id 軟刪除公司單筆資料
func (c *Video) SoftDelete(id int) bool {
	now := time.Now().Format(time.DateTime)

	_, err := DB.Exec("UPDATE videos SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, id)

	return err == nil
}

// 以 id 刪除影片單筆資料
func (c *Video) Delete(id int) bool {
	_, err := DB.Exec("DELETE FROM videos WHERE id = ? AND deleted_at IS NULL", id)
	return err == nil
}

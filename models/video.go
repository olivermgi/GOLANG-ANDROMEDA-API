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
		"SELECT id, status, title, description updated_at FROM videos WHERE deleted_at IS NULL ORDER BY %s %s, id DESC LIMIT ? OFFSET ?",
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
		err := rows.Scan(&video.Id, &video.Status, &video.Title, &video.Description, &video.UpdatedAt)
		if err != nil {
			return make([]Video, 0), 0, 0
		}
		videos = append(videos, video)
	}

	count := 0
	DB.QueryRow("SELECT COUNT(*) FROM videos WHERE deleted_at IS NULL").Scan(&count)
	lastPage = int(math.Ceil(float64(count) / float64(perPage)))
	return videos, count, lastPage
}

func (c *Video) AllPublish() []Video {
	queryStr := fmt.Sprintln(`
	SELECT videos.id, videos.status, title, description, video_files.name, videos.updated_at FROM videos 
	INNER JOIN video_files ON videos.id = video_files.video_id
	WHERE videos.status = 'publish' AND videos.deleted_at IS NULL AND video_files.status = 'transformed' ORDER BY videos.updated_at DESC`)

	rows, err := DB.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	videos := make([]Video, 0)
	for rows.Next() {
		var video Video
		var videoFile VideoFile
		err := rows.Scan(&video.Id, &video.Status, &video.Title, &video.Description, &videoFile.Name, &video.UpdatedAt)
		if err != nil {
			return make([]Video, 0)
		}
		video.VideoFile = videoFile
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

// 以 id 軟刪除單筆影片資料
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

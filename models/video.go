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
}

// 新增影片單筆資料
func (c *Video) Insert(data Video) *Video {
	result, err := DB.Exec("INSERT INTO videos(status, title, description, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		data.Status, data.Title, data.Description)
	if err != nil {
		log.Println("新增影片資料失敗，錯誤訊息：", err)
		return nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("取得新增的影片 ID 失敗，錯誤訊息：", err)
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
		"SELECT id, status, title, updated_at FROM videos ORDER BY %s %s, id DESC LIMIT ? OFFSET ?",
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
			log.Println("取得多筆公司資料失敗，錯誤訊息：", err)
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
	queryStr := fmt.Sprintln("SELECT id, status, title, updated_at FROM videos")

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
			log.Println("取得多筆公司資料失敗，錯誤訊息：", err)
			return make([]Video, 0)
		}
		videos = append(videos, video)
	}

	return videos
}

// 以 id 抓取影片單筆資料
func (c *Video) Get(id int) *Video {
	var video Video
	err := DB.QueryRow("SELECT id, status, title, description FROM videos WHERE id = ?", id).
		Scan(&video.Id, &video.Status, &video.Title, &video.Description)

	if err != nil {
		log.Println("取得單筆影片資料失敗，錯誤訊息：", err)
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

	_, err := DB.Exec("UPDATE videos SET status = ?, title = ?, description = ?, updated_at = ? WHERE id = ?",
		data.Status, data.Title, data.Description, now, id)

	if err != nil {
		log.Println("更新影片資料失敗，錯誤訊息：", err)
		return nil
	}

	return &Video{
		Id:          id,
		Status:      data.Status,
		Title:       data.Title,
		Description: data.Description,
	}
}

// 以 id 刪除影片單筆資料
func (c *Video) Delete(id int) bool {
	_, err := DB.Exec("DELETE FROM videos WHERE id = ?", id)
	if err != nil {
		log.Println("刪除影片資料失敗，錯誤訊息：", err)
		return false
	}

	return true
}

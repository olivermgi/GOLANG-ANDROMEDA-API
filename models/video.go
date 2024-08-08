package models

import (
	"fmt"
	"log"
)

type Video struct {
	Id          int    `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

// 新增公司單筆資料
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
	fmt.Println("公司新增成功")

	return &Video{
		Id:          int(id),
		Status:      data.Status,
		Title:       data.Title,
		Description: data.Description,
	}
}

// // 抓取所有公司資料，並且以 id 欄位降冪排序
// func (c *Company) All() []Company {
// 	companies := make([]Company, 0)
// 	rows, err := DB.Query("SELECT * FROM companies ORDER BY id DESC")
// 	if err != nil {
// 		log.Println("取得多筆公司資料失敗，錯誤訊息：", err)
// 		return companies
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var company Company
// 		err := rows.Scan(&company.Id, &company.Name, &company.Address)
// 		if err != nil {
// 			log.Println("取得多筆公司資料失敗，錯誤訊息：", err)
// 			return companies
// 		}
// 		companies = append(companies, company)
// 	}

// 	return companies
// }

// // 以 id 抓取公司單筆資料
// func (c *Company) Get(id int) *Company {
// 	var company Company
// 	err := DB.QueryRow("SELECT * FROM companies WHERE id = ?", id).Scan(&company.Id, &company.Name, &company.Address)
// 	if err != nil {
// 		log.Println("取得單筆公司資料失敗，錯誤訊息：", err)
// 		return nil
// 	}

// 	return &company
// }

// // 以 id 更新公司單筆資料
// func (c *Company) Update(id int, data Company) *Company {
// 	_, err := DB.Exec("UPDATE companies SET name = ?, address = ? WHERE id = ?", data.Name, data.Address, id)
// 	if err != nil {
// 		log.Println("更新公司資料失敗，錯誤訊息：", err)
// 		return nil
// 	}

// 	return &Company{
// 		Id:      int(id),
// 		Name:    data.Name,
// 		Address: data.Address,
// 	}
// }

// // 以 id 刪除公司單筆資料
// func (c *Company) Delete(id int) bool {
// 	_, err := DB.Exec("DELETE FROM companies WHERE id = ?", id)
// 	if err != nil {
// 		log.Println("更新公司資料失敗，錯誤訊息：", err)
// 		return false
// 	}

// 	return true
// }

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/services"
)

// var model *models.Video

// // 顯示影片列表 (分頁)
// func IndexVideo(w http.ResponseWriter, r *http.Request) {
// 	companies := model.All()
// 	common.Response(companies, http.StatusOK, "", w)
// }

// 新增影片資料
func StoreVideo(w http.ResponseWriter, r *http.Request) {
	var videoData services.VideoData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&videoData)
	if err != nil {
		common.Abort(http.StatusForbidden, "JSON 格式不正確", nil)
	}

	common.ValidateStruct(videoData)

	video := services.StoreVideo(videoData)
	if video == nil {
		common.Abort(http.StatusForbidden, "影片資料新增失敗", nil)
	}

	common.Response(http.StatusCreated, "影片資料新增成功", video, w)
}

// // 顯示單筆公司資料
// func ShowVideo(w http.ResponseWriter, r *http.Request) {
// 	company_id, err := strconv.Atoi(r.PathValue("company_id"))
// 	if err != nil {
// 		common.Response(struct{}{}, http.StatusForbidden, "參數類型錯誤", w)
// 		return
// 	}

// 	company := model.Get(company_id)
// 	if company == nil {
// 		common.Response(struct{}{}, http.StatusNotFound, "無此公司資料", w)
// 		return
// 	}

// 	common.Response(company, http.StatusOK, "", w)
// }

// // 更新單筆公司資料
// func UpdateVideo(w http.ResponseWriter, r *http.Request) {
// 	company_id, err := strconv.Atoi(r.PathValue("company_id"))
// 	if err != nil {
// 		common.Response(struct{}{}, http.StatusForbidden, "參數類型錯誤", w)
// 		return
// 	}

// 	var companyData models.Company
// 	decoder := json.NewDecoder(r.Body)
// 	err = decoder.Decode(&companyData)
// 	if err != nil {
// 		common.Response(struct{}{}, http.StatusForbidden, "JSON 格式不正確", w)
// 		return
// 	}

// 	if companyData.Name == "" || companyData.Address == "" {
// 		common.Response(struct{}{}, http.StatusForbidden, "資料格式不正確", w)
// 		return
// 	}

// 	if model.Get(company_id) == nil {
// 		common.Response(struct{}{}, http.StatusNotFound, "無此公司資料", w)
// 		return
// 	}

// 	company := model.Update(company_id, companyData)
// 	if company == nil {
// 		common.Response(struct{}{}, http.StatusForbidden, "公司更新失敗", w)
// 		return
// 	}

// 	common.Response(company, http.StatusOK, "公司更新成功", w)
// }

// // 刪除一筆公司資料
// func DestroyVideo(w http.ResponseWriter, r *http.Request) {
// 	company_id, err := strconv.Atoi(r.PathValue("company_id"))
// 	if err != nil {
// 		common.Response(struct{}{}, http.StatusForbidden, "參數類型錯誤", w)
// 		return
// 	}

// 	if model.Get(company_id) == nil {
// 		common.Response(struct{}{}, http.StatusNotFound, "無此公司資料", w)
// 		return
// 	}

// 	is_success := model.Delete(company_id)
// 	if !is_success {
// 		common.Response(struct{}{}, http.StatusForbidden, "公司刪除失敗", w)
// 		return
// 	}

// 	common.Response(struct{}{}, http.StatusOK, "公司刪除成功", w)
// }

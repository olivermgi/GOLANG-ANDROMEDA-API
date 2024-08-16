package rules

type VideoIndex struct {
	Page       int    `validate:"required,min=1" field:"頁數"`
	PerPage    int    `validate:"required,oneof=10 20 30" field:"筆數"`
	Sort       string `validate:"required,oneof=asc desc" field:"排序方式"`
	SortColumn string `validate:"required,oneof=id updated_at" field:"排序欄位"`
}

type VideoStore struct {
	Status      string `validate:"required,oneof=publish unpublish" field:"狀態"`
	Title       string `validate:"required,max=50" field:"標題"`
	Description string `validate:"omitempty,max=255" field:"描述"`
}

type VideoShow struct {
	VideoId int `validate:"required,min=1" field:"video_id"`
}

type VideoUpdate struct {
	VideoId     int    `validate:"required,min=1" field:"video_id"`
	Status      string `validate:"required,oneof=publish unpublish" field:"狀態"`
	Title       string `validate:"required,max=50" field:"標題"`
	Description string `validate:"omitempty,max=255" field:"描述"`
}

type VideoDelete struct {
	VideoId int `validate:"required,min=1" field:"video_id"`
}

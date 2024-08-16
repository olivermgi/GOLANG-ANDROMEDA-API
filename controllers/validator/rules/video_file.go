package rules

import "mime/multipart"

type VideoFileStore struct {
	VideoId int                   `validate:"required,min=1" field:"video_id"`
	File    multipart.File        `validate:"file_exists" field:"此檔案"`
	Header  *multipart.FileHeader `validate:"required,max_file_size=52428800,file_mimes=video/mp4" field:"此檔案"`
}

type VideoFileShow struct {
	VideoId int `validate:"required,min=1" field:"video_id"`
}

type VideoFileDelete struct {
	VideoId int `validate:"required,min=1" field:"video_id"`
}

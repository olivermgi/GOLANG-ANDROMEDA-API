package validator

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	locale_zh_tw "github.com/go-playground/locales/zh_Hant_TW"  // 語言環境包
	ut "github.com/go-playground/universal-translator"          // 翻譯器
	"github.com/go-playground/validator/v10"                    // 驗證器
	"github.com/go-playground/validator/v10/translations/zh_tw" // 語言包
	"github.com/olivermgi/golang-crud-practice/common"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	validate = validator.New()
	utObj := ut.New(locale_zh_tw.New())
	trans, _ = utObj.GetTranslator("zh_hant_tw")

	validate.RegisterValidation("file_exists", fileExists)
	validate.RegisterValidation("max_file_size", maxFileSize)
	validate.RegisterValidation("file_mimes", fileMines)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("field")
	})

	validate.RegisterTranslation("file_exists", trans, func(ut ut.Translator) error {
		return ut.Add("file_exists", "{0}為必填欄位", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("file_exists", fe.Field())
		return t
	})

	validate.RegisterTranslation("max_file_size", trans, func(ut ut.Translator) error {
		return ut.Add("max_file_size", "{0}大小必須小於 {1} MB", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		maxSize := fmt.Sprintf("%d", common.StringToInt(fe.Param())/1024/1024)
		t, _ := ut.T("max_file_size", fe.Field(), maxSize)
		return t
	})

	validate.RegisterTranslation("file_mimes", trans, func(ut ut.Translator) error {
		return ut.Add("file_mimes", "{0}的格式只能符合：{1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		mimes := strings.Replace(fe.Param(), " ", ", ", -1)
		t, _ := ut.T("file_mimes", fe.Field(), mimes)
		return t
	})

	err := zh_tw.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatalln("驗證器初始化錯誤：", err)
	}
}

func ValidateOrAbort(ruleData interface{}) {
	errs := validate.Struct(ruleData)
	if errs != nil {
		errorData := make(common.ErrorMap)
		errorData["validation"] = make(map[string]string)
		for _, err := range errs.(validator.ValidationErrors) {
			key := strings.ToLower(err.StructField())
			errText := err.Translate(trans)
			validationMap := errorData["validation"].(map[string]string)
			validationMap[key] = errText
		}
		common.AbortWithData(http.StatusUnprocessableEntity, "輸入資料驗證失敗", errorData)
	}
}

func fileExists(fl validator.FieldLevel) bool {
	file := fl.Field().Interface().(multipart.File)
	fmt.Println("test file:", file)
	return file != nil
}

func maxFileSize(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	maxSize, _ := strconv.ParseInt(fl.Param(), 10, 64)
	fileSize := fileHeader.Size

	return fileSize <= maxSize
}

func fileMines(fl validator.FieldLevel) bool {
	if fl.Param() == "" {
		return true
	}

	fileHeader := fl.Field().Interface().(multipart.FileHeader)

	acceptedMines := strings.Split(fl.Param(), " ")
	fileMine := fileHeader.Header.Get("Content-Type")

	return slices.Contains(acceptedMines, fileMine)
}

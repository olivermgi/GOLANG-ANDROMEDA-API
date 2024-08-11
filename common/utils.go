package common

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	locale_zh_tw "github.com/go-playground/locales/zh_Hant_TW"  // 語言環境包
	ut "github.com/go-playground/universal-translator"          // 翻譯器
	"github.com/go-playground/validator/v10"                    // 驗證器
	"github.com/go-playground/validator/v10/translations/zh_tw" // 語言包
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	validate = validator.New()
	ut := ut.New(locale_zh_tw.New())
	trans, _ = ut.GetTranslator("zh_hant_tw")

	err := zh_tw.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatalln("驗證器初始化錯誤：", err)
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("field")
	})
}

func StringToInt(str string) int {
	number, _ := strconv.Atoi(str)
	return number
}

func DumpDie(data interface{}) {

	//panic(&HttpJsonError{StatusCode: statusCode, Message: message, ErrorData: errorData})
}

func Abort(statusCode int, message string) {
	errorData := make(ErrorMap)
	AbortWithData(statusCode, message, errorData)
}

func AbortWithData(statusCode int, message string, errorData ErrorMap) {
	panic(&HttpJsonError{StatusCode: statusCode, Message: message, ErrorData: errorData})
}

func Response(statusCode int, message string, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	responseData := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}

	isError := false
	if !(statusCode >= 200 && statusCode < 300) {
		isError = true
	}

	if data == nil {
		data = struct{}{}
	}

	if isError {
		responseData["errors"] = data
	} else {
		responseData["data"] = data
	}

	json.NewEncoder(w).Encode(responseData)
}

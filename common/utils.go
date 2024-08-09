package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

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
		return field.Tag.Get("validate_field")
	})
}

type ErrorMap map[string]interface{}

func ValidateStruct(data interface{}) {
	errs := validate.Struct(data)
	if errs != nil {
		errorData := make(ErrorMap)
		errorData["validation"] = make(map[string]string)
		for _, err := range errs.(validator.ValidationErrors) {
			key := strings.ToLower(err.StructField())
			errText := err.Translate(trans)
			validationMap := errorData["validation"].(map[string]string)
			validationMap[key] = errText
			// errorData["validation"][key] = errText
		}

		fmt.Println(errorData)
		Abort(http.StatusUnprocessableEntity, "輸入資料驗證失敗", errorData)
	}
}

type HttpJsonError struct {
	StatusCode int
	Message    string
	ErrorData  ErrorMap
}

func (e *HttpJsonError) Error() string {
	return fmt.Sprintf("StatusCode:%d, Message:%s", e.StatusCode, e.Message)
}

func Abort(statusCode int, message string, errorData ErrorMap) {
	if errorData == nil {
		errorData = make(ErrorMap)
	}
	panic(&HttpJsonError{StatusCode: statusCode, Message: message, ErrorData: errorData})
}

func Response(statusCode int, message string, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	responseData := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}
	fmt.Println(data)
	_, is_error := data.(ErrorMap)
	if is_error {
		responseData["errors"] = data
	} else {
		responseData["data"] = data
	}

	json.NewEncoder(w).Encode(responseData)
}

package common

import (
	"encoding/json"
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

type ValidateErrors map[string]map[string]string

func ValidateStruct(data interface{}) ValidateErrors {
	errs := validate.Struct(data)
	if errs != nil {
		validateErrors := make(ValidateErrors)
		validateErrors["validation"] = make(map[string]string)
		for _, err := range errs.(validator.ValidationErrors) {
			key := strings.ToLower(err.StructField())
			errText := err.Translate(trans)
			validateErrors["validation"][key] = errText
		}

		return validateErrors
	}

	return nil
}

func Response(data interface{}, statusCode int, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	responseData := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}

	_, is_error := data.(ValidateErrors)
	if is_error {
		responseData["errors"] = data
	} else {
		responseData["data"] = data
	}

	json.NewEncoder(w).Encode(responseData)
}

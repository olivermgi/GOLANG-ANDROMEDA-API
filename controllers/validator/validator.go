package validator

import (
	"log"
	"net/http"
	"reflect"
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

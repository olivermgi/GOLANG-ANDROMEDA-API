package middleware

import (
	"fmt"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
)

type ErrorResponseMiddleware struct {
	Next http.Handler
}

// 集中管理例外的錯誤
func (b *ErrorResponseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		statusCode := 500
		message := "未知的錯誤"
		errors := common.ErrorMap{"error": fmt.Sprint(err)}
		httpJsonError, ok := err.(*common.HttpJsonError)
		if ok {
			statusCode = httpJsonError.StatusCode
			message = httpJsonError.Message
			errors = httpJsonError.ErrorData
		}

		common.Response(statusCode, message, errors, w)
	}()

	b.Next = &BasicAuthMiddleware{}

	b.Next.ServeHTTP(w, r)
}

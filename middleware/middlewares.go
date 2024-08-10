package middleware

import (
	"fmt"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/common"
)

type Middlewares struct {
	handlers []http.Handler
}

func handlers() []http.Handler {
	return []http.Handler{
		&BasicAuthMiddleware{},
	}
}

func (m *Middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if m.handlers == nil {
		m.handlers = make([]http.Handler, 0)
		m.handlers = append(m.handlers, handlers()...)
	}

	for _, handler := range m.handlers {
		handler.ServeHTTP(w, r)
	}

	http.DefaultServeMux.ServeHTTP(w, r)
}

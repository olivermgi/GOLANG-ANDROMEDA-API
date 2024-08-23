package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/olivermgi/golang-andromeda-api/common"
	"github.com/olivermgi/golang-andromeda-api/config"
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

		errMessage := ""
		if !config.IsProduction() {
			errMessage = fmt.Sprint(err)
		}

		statusCode := 500
		message := "未知的錯誤"
		errors := common.ErrorMap{"error": errMessage}
		httpJsonError, ok := err.(*common.HttpJsonError)
		if ok {
			statusCode = httpJsonError.StatusCode
			message = httpJsonError.Message
			errors = httpJsonError.ErrorData
		} else {
			log.Println(err)
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			fmt.Printf("Stack trace:\n%s\n", buf[:n])
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

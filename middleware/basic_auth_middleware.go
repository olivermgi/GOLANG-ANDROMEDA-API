package middleware

import (
	"net/http"
	"strings"

	"github.com/olivermgi/golang-crud-practice/common"
	"github.com/olivermgi/golang-crud-practice/config"
)

type BasicAuthMiddleware struct {
}

// 以 Basic Auth 的方式，簡單增加驗證機制，只會判斷路徑為 /admin 開頭的 API
func (m *BasicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/internal") {
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		common.Abort(http.StatusUnauthorized, "認證不合法")
	}

	authInfo := config.GetAuthConfig()

	if username != authInfo["username"] || password != authInfo["password"] {
		common.Abort(http.StatusUnauthorized, "帳號密碼不正確")
	}
}

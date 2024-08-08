package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/olivermgi/golang-crud-practice/common"
)

type BasicAuthMiddleware struct {
	Next http.Handler
}

// 以 Basic Auth 的方式，簡單增加驗證機制，只會判斷路徑為 /admin 開頭的 API
func (b *BasicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if b.Next == nil {
		b.Next = http.DefaultServeMux
	}

	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/internal") {
		b.Next.ServeHTTP(w, r)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("認證不合法")
		common.Response(struct{}{}, http.StatusUnauthorized, "認證不合法", w)
		return
	}

	if username != "admin" || password != "123456" {
		log.Println("帳號密碼不正確")
		common.Response(struct{}{}, http.StatusUnauthorized, "帳號密碼不正確", w)
		return
	}

	b.Next.ServeHTTP(w, r)
}

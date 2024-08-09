package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/middleware"
	_ "github.com/olivermgi/golang-crud-practice/models"
	"github.com/olivermgi/golang-crud-practice/routes"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("恢復了panic:", r)
		}
	}()

	server := http.Server{
		Addr:    ":8080",
		Handler: &middleware.ErrorResponseMiddleware{},
	}

	routes.RegisterAPIRoutes()

	log.Println("網頁伺服器正在運行中...")
	server.ListenAndServe()
}

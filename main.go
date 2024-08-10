package main

import (
	"log"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/middleware"
	_ "github.com/olivermgi/golang-crud-practice/models"
	"github.com/olivermgi/golang-crud-practice/routes"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &middleware.ErrorResponseMiddleware{},
	}

	routes.RegisterAPIRoutes()

	log.Println("網頁伺服器正在運行中...")
	server.ListenAndServe()
}

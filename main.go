package main

import (
	"log"
	"net/http"

	_ "github.com/olivermgi/golang-crud-practice/controllers/validator"
	"github.com/olivermgi/golang-crud-practice/middleware"
	_ "github.com/olivermgi/golang-crud-practice/models"
	"github.com/olivermgi/golang-crud-practice/routes"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &middleware.Middlewares{},
	}

	routes.RegisterAPIRoutes()

	log.Println("網頁伺服器正在運行中...")
	server.ListenAndServe()
}

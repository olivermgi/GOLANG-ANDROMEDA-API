package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/olivermgi/golang-crud-practice/config"
	_ "github.com/olivermgi/golang-crud-practice/controllers/validator"
	"github.com/olivermgi/golang-crud-practice/middleware"
	_ "github.com/olivermgi/golang-crud-practice/models"
	"github.com/olivermgi/golang-crud-practice/routes"
)

func main() {
	port := config.GetPort()

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: &middleware.Middlewares{},
	}

	routes.RegisterAPIRoutes()

	log.Println("網頁伺服器正在運行中...")
	server.ListenAndServe()
}

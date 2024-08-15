package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/olivermgi/golang-crud-practice/config"
	_ "github.com/olivermgi/golang-crud-practice/controllers/validator"
	"github.com/olivermgi/golang-crud-practice/middleware"
	_ "github.com/olivermgi/golang-crud-practice/models"
	"github.com/olivermgi/golang-crud-practice/routes"
)

func main() {
	serverConfig := config.GetServerConfig()

	port := serverConfig["port"]
	secure, _ := strconv.ParseBool(serverConfig["secure"])
	certificatePath := serverConfig["certificate_path"]
	privateKeyPath := serverConfig["private_key_path"]

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: &middleware.Middlewares{},
	}

	routes.RegisterAPIRoutes()

	if secure {
		log.Println("HTTPS 網頁伺服器正在運行中...")
		server.ListenAndServeTLS(certificatePath, privateKeyPath)
		return
	}
	log.Println("HTTP 網頁伺服器正在運行中...")
	server.ListenAndServe()
}

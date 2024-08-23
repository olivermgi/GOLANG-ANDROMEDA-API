package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/olivermgi/golang-andromeda-api/config"
	_ "github.com/olivermgi/golang-andromeda-api/controllers/validator"
	"github.com/olivermgi/golang-andromeda-api/middleware"
	_ "github.com/olivermgi/golang-andromeda-api/models"
	"github.com/olivermgi/golang-andromeda-api/routes"
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

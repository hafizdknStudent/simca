package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/fsnotify.v1"

	"simka/config"
	"simka/handler"
	"simka/simka"
	"simka/utilities"
)

func main() {
	config.LoadEnv()
	config.LoadConfig()

	utilities.GenerateOauthCookie()

	jar := utilities.NewJar()
	simakaService := simka.NewService(jar)
	simakaHandler := handler.NewSimakaHandler(simakaService)

	googleHandler := handler.NewGoogleHandler()

	router := gin.Default()

	router.Static("assets", "templates/assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", googleHandler.GetGoogleLogin)
	router.GET("/login-google", googleHandler.PostGoogleLogin)
	router.GET("/google_callback", googleHandler.GetGoogleCallBack)
	router.GET("/simaka", simakaHandler.GetLoginSimaka)
	router.POST("/login-simaka", simakaHandler.PostLoginSimaka)
	router.GET("/success", simakaHandler.SuccessLogin)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasSuffix(event.Name, "app_offline.htm") {
					fmt.Println("Exiting due to app_offline.htm being present")
					os.Exit(0)
				}
			}
		}
	}()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err := watcher.Add(currentDir); err != nil {
		fmt.Println("ERROR", err)
	}

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
	// router.Run(":8080")
}

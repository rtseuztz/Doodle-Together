package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func serveHome(c *gin.Context) {
	log.Println(c.Request.URL)
	if c.Request.URL.Path != "/" {
		http.Error(c.Writer, "Not found", http.StatusNotFound)
		return
	}
	if c.Request.Method != http.MethodGet {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(c.Writer, c.Request, "index.html")
}
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		router := gin.New()
		router.Use(gin.Logger())
		router.LoadHTMLGlob("templates/*.tmpl.html")
		router.Static("/static", "static")
		hub := newHub()
		go hub.run()
		router.GET("/", serveHome) //.HandleFunc("/", serveHome)
		router.GET("/wss", func(c *gin.Context) {
			serveWs(hub, c.Writer, c.Request)
		})

		router.Run(":" + "8000")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	hub := newHub()
	go hub.run()
	router.GET("/", serveHome) //.HandleFunc("/", serveHome)
	router.GET("/wss", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	router.Run(":" + port)
}

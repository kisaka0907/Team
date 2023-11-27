package main

import (
	"net/http"
	"web-chat/controllers"
	"web-chat/initializers"
	"web-chat/middleware"
	"web-chat/websocket"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", gin.H{
			"title": "Auth",
		})
	})

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	// Home
	r.GET("/home", middleware.RequireAuth, controllers.Lists)

	// search friend
	r.POST("/home/search_friend", middleware.RequireAuth, controllers.SearchFriend)
	// add friend
	r.POST("/home/add_friend", middleware.RequireAuth, controllers.AddFriend)
	// create group
	r.POST("/home/create_group", middleware.RequireAuth, controllers.CreateGroup)
	// delete friend
	r.POST("/home/destroy_friend", middleware.RequireAuth, controllers.DeleteFriend)
	// delete room
	r.POST("/home/destroy_group", middleware.RequireAuth, controllers.DeleteGroup)

	//
	//
	// chat history
	r.GET("/chat/chat", middleware.RequireAuth, controllers.ListChatHistory)
	// chat list
	r.GET("/chat/chatlist", middleware.RequireAuth, controllers.ChatList)
	r.POST("/chat/create", middleware.RequireAuth, controllers.CreateChat)

	// ws
	hub := websocket.NewHub()
	go hub.Run()
	r.GET("/ws/:id", func(c *gin.Context) {
		controllers.ServeRoomWs(c)
	})
	r.Run()
}

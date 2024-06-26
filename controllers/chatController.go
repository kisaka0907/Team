package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web-chat/initializers"
	"web-chat/models"

	"github.com/gin-gonic/gin"
)

// 誰作？
func CreateChat(c *gin.Context) {
	user, _ := c.Get("user")
	content := &models.Req_reseiver{}
	fmt.Print("ccstart")
	if err := c.Bind(content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Print("Content:")
	// fmt.Print(content.Content)
	// fmt.Print("ID:")
	// fmt.Print(content.RoomID)

	// チャットデータを構造体に格納
	chat := &models.Chat_history{
		Content: content.Content,
		// RoomID:  87,
		RoomID: content.RoomID,
		UserID: user.(models.Users).ID,
	}
	// chat_historyテーブルにデータを挿入
	result := initializers.DB.Create(&chat)
	// エラー処理s
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "chat.html", gin.H{
			"title":  user.(models.Users).UserName,
			"result": "todoの作成に失敗しました",
		})
		return
	}
	// path := "/chat/chat?id=" + fmt.Sprint(chat.RoomID)
	// path := "/chat/chat?id=87"
	// // チャットの作成後,ページを更新
	// c.Redirect(http.StatusSeeOther, path)
	c.JSON(http.StatusOK, chat)
}

// func FriendsChatList(c *gin.Context) {
// 	user, _ := c.Get("user")
// 	// フレンドチャット
// 	var friends []models.Friends
// 	initializers.DB.Where("user_id = ?", user.(models.User).ID).Find(&friends)
// 	// グループチャット
// 	var groups []models.Groups
// 	initializers.DB.Where("user_id = ?", user.(models.User).ID).Find(&groups)

// 	c.HTML(http.StatusOK, "chatlist.html", gin.H{
// 		"title":   user.(models.User).UserName + "'s Chatlist",
// 		"friends": friends,
// 		"groups":  groups,
// 	})
// }

// ryo作
func ChatList(c *gin.Context) {
	user, _ := c.Get("user")

	var chatroom []models.Friends
	var room1 []models.Rooms
	var grouproom []models.Groups
	var room2 []models.Rooms

	initializers.DB.Where("user_id1 = ? OR user_id2 = ?", user.(models.Users).ID, user.(models.Users).ID).Find(&chatroom)
	initializers.DB.Where("users_refer = ?", user.(models.Users).ID).Find(&grouproom)

	// 配列(スライス)の要素毎のループ処理
	// for インデックス番号, 値の変数 := range 配列orスライス{繰り返し処理}
	for _, v := range chatroom {
		var room models.Rooms
		initializers.DB.Where("id = ?", v.RoomsRefer).Find(&room)
		room1 = append(room1, room) // append(スライス,追加する値)
	}

	for _, v := range grouproom {
		var room models.Rooms
		initializers.DB.Where("id = ?", v.RoomsRefer).Find(&room)
		//変更部分
		if room.RoomName != "" {
			room2 = append(room2, room)
		}
	}

	c.HTML(http.StatusOK, "chatlist.html", gin.H{
		"title": user.(models.Users).UserName + "'s Chat",
		// "chatroom": chatroom,
		"room1": room1,
		"room2": room2,
	})
}

// 誰作？
func ListChatHistory(c *gin.Context) {
	// user, _ := c.Get("user")
	// userID := user.(models.User).ID
	var chat_history []models.Chat_history
	var room models.Rooms
	room_id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Fatalln(err)
	}
	// initializers.DB.Find(&chat_history, "room_id = ?", user.(models.Chat_history).ID)
	initializers.DB.Where("room_id = ?", room_id).Find(&chat_history)
	initializers.DB.Where("id = ?", room_id).Find(&room)
	// fmt.Print(chat_history)
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"title":        room.RoomName + "'s Chat History",
		"chat_history": chat_history,
		"room":         room,
	})
}

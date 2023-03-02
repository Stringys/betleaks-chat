package models

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Stringys/betleaks-chat/db"
	"github.com/gorilla/websocket"
)

type User struct {
	Username string
	Conn     *websocket.Conn
	Room     *Chat
	PgDB     *db.PostgresDB
}

func NewUser(username string, conn *websocket.Conn, room *Chat, pgDB *db.PostgresDB) *User {
	return &User{
		Username: username,
		Conn:     conn,
		Room:     room,
		PgDB:     pgDB,
	}
}

func (u *User) StartReadingMessages(chatID, fromID int) {
	defer func() {
		u.Room.Leave <- u
	}()

	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			log.Printf("Error reading message in websocket: %v", err)
			break
		} else {
			messageRead := string(message)

			msg, err := u.PgDB.Client.Message.Create().
				SetFromID(fromID).
				SetBody(messageRead).
				Save(context.Background())
			if err != nil {
				log.Printf("Error saving message in database: %v", err)
				break
			}

			chat, err := u.PgDB.Client.Chat.
				UpdateOneID(chatID).
				AddMessageIDs(msg.ID).
				Save(context.Background())
			if err != nil {
				log.Printf("Error updating chat in database: %v", err)
			}

			log.Printf("Chat with id: %d has been updated", chat.ID)

			u.Room.Messages <- NewMessage(messageRead, u.Username)
		}
	}
}

func (u *User) WriteMessage(msg *Message) {
	msgEncoded, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error encoding msg json: %v", err)
	}

	if err = u.Conn.WriteMessage(websocket.TextMessage, msgEncoded); err != nil {
		log.Printf("Error writing message: %v", err)
	}
}

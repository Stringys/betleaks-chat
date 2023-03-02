package service

import (
	"fmt"

	"github.com/Stringys/betleaks-chat/db"
	"github.com/Stringys/betleaks-chat/service/models"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	PgDB  *db.PostgresDB
	chats map[int]*models.Chat
}

func NewChatService(pgDB *db.PostgresDB) ChatService {
	return ChatService{PgDB: pgDB}
}

func (s *ChatService) CreateChat(to, from string, ws *websocket.Conn) error {
	receiver, err := s.PgDB.FindUserByUsername(to)
	if err != nil {
		// WRITE A LOG
		return err
	}

	sender, err := s.PgDB.FindUserByUsername(from)
	if err != nil {
		// WRITE A LOG
		return err
	}

	chat, err := s.PgDB.FindChatByUsernames(to, from)
	if chat != nil || err != nil {
		// WRITE A LOG
		return err
	}

	chat, err = s.PgDB.CreateChat(receiver.ID, sender.ID)
	if err != nil {
		// WRITE A LOG
		return err
	}

	room := models.NewChat(chat.ID)
	s.chats[room.ID] = room
	go room.Run()

	user := models.NewUser(sender.Username, ws, room, s.PgDB)
	room.Join <- user
	user.StartReadingMessages(chat.ID, sender.ID)

	return nil
}

func (s *ChatService) JoinChat(chatID int, username string, ws *websocket.Conn) error {
	userFromDB, err := s.PgDB.FindUserByUsername(username)
	if err != nil {
		return err
	}

	if room, ok := s.chats[chatID]; ok {
		user := models.NewUser(username, ws, room, s.PgDB)
		room.Join <- user
		user.StartReadingMessages(chatID, userFromDB.ID)
		return nil
	} else {
		return fmt.Errorf("chat not found")
	}
}

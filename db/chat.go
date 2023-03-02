package db

import (
	"context"

	"github.com/Stringys/betleaks-chat/ent"
	ChatEnt "github.com/Stringys/betleaks-chat/ent/chat"
	UserEnt "github.com/Stringys/betleaks-chat/ent/user"
)

func (pg *PostgresDB) GetChatMembers(chatID int) ([]*ent.User, error) {
	chat, err := pg.Client.Chat.Get(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	members, err := chat.QueryMembers().All(context.Background())
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (pg *PostgresDB) GetChats(username string) ([]*ent.Chat, error) {
	user, err := pg.Client.User.Query().Where(UserEnt.UsernameEQ(username)).First(context.Background())
	if err != nil {
		return nil, err
	}

	chats, err := user.QueryChats().All(context.Background())
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (pg *PostgresDB) FindChatByUsernames(to, from string) (*ent.Chat, error) {
	chat, err := pg.Client.Chat.Query().Where(
		ChatEnt.HasMembersWith(
			UserEnt.UsernameEQ(to),
			UserEnt.UsernameEQ(from),
		),
	).Only(context.Background())
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (pg *PostgresDB) CreateChat(to, from int) (*ent.Chat, error) {
	chat, err := pg.Client.Chat.Create().AddMemberIDs(to, from).Save(context.Background())
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (pg *PostgresDB) FindUserByUsername(username string) (*ent.User, error) {
	user, err := pg.Client.User.Query().Where(UserEnt.UsernameEQ(username)).First(context.Background())
	if err != nil {
		return nil, err
	}

	return user, err
}

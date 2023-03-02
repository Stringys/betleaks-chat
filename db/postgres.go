package db

import (
	"context"
	"log"

	"github.com/Stringys/betleaks-chat/ent"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Client *ent.Client
}

func NewPostgresDB(dataSource string) (*PostgresDB, error) {
	client, err := ent.Open("postgres", dataSource)
	if err != nil {
		log.Printf("failed opening connection to postgres: %v", err)
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Printf("failed creating schema: %v", err)
		return nil, err
	}

	return &PostgresDB{Client: client}, nil
}

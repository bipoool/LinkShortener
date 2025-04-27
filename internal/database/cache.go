package database

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
)

type Cache interface {
	Get(key string) (string, error)
	Upsert(key string, value string) error
	Delete(key string) error
	Close()
}

type Dice struct {
	db *dicedb.Client
}

func (dice *Dice) Get(key string) (string, error) {
	resp := dice.db.Fire(&wire.Command{Cmd: "GET", Args: []string{key}})

	if resp.Status == wire.Status_ERR {
		return "", fmt.Errorf(resp.Message)
	}

	return resp.GetGETRes().Value, nil
}

func (dice *Dice) Upsert(key string, value string) error {
	resp := dice.db.Fire(&wire.Command{Cmd: "SET", Args: []string{key, value}})
	if resp.Status == wire.Status_ERR {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func (dice *Dice) Delete(key string) error {
	resp := dice.db.Fire(&wire.Command{Cmd: "DELETE", Args: []string{key}})
	if resp.Status == wire.Status_ERR {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func (dice *Dice) Close() {
	dice.db.Close()
}

func NewCache(host string, port string) Cache {
	portInt, _ := strconv.Atoi(port)
	dice, err := dicedb.NewClient(host, portInt)
	if err != nil {
		slog.Error("Failed to initialize Cache DB")
		panic(err)
	}

	cache := &Dice{
		db: dice,
	}

	return cache
}

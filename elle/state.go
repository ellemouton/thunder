package main

import (
	"database/sql"
	"fmt"

	"github.com/ellemouton/thunder/elle/db"
	"github.com/ellemouton/thunder/syncchain"
	syncchain_client "github.com/ellemouton/thunder/syncchain/client"
)

type State struct {
	db *sql.DB

	syncchainClient syncchain.Client
}

func newState() (*State, error) {
	s := new(State)

	db, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %s", err)
	}
	s.db = db

	s.syncchainClient, err = syncchain_client.New()
	if err != nil {
		return nil, fmt.Errorf("connecting to syncchain: %s", err)
	}

	return s, nil
}

func (s *State) GetDB() *sql.DB {
	return s.db
}

func (s *State) SyncchainClient() syncchain.Client {
	return s.syncchainClient
}

func (s *State) cleanup() {
	if err := s.db.Close(); err != nil {
		fmt.Errorf("closing db: %v", err)
	}

	if err := s.syncchainClient.Close(); err != nil {
		fmt.Errorf("closing syncchain: %v", err)
	}
}

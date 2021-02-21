package main

import (
	"database/sql"
	"fmt"

	"github.com/ellemouton/thunder/db"
)

type State struct {
	db *sql.DB
}

func newState() (*State, error) {
	s := new(State)

	db, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %s", err)
	}
	s.db = db

	return s, nil
}

func (s *State) GetDB() *sql.DB {
	return s.db
}

func (s *State) cleanup() {
	if err := s.db.Close(); err != nil {
		fmt.Errorf("closing db: %v", err)
	}
}

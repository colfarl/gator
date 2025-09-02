package main

import (
	"github.com/colfarl/gator/internal/config"
	"github.com/colfarl/gator/internal/database"
)

type state struct {
	CurrentState			*config.Config	
	db						*database.Queries
}

func newState(c *config.Config, q *database.Queries) state {
	return state{
		CurrentState: c,
		db: q,
	}
}

func (s *state) updateUser(newUser string) error {

	err := s.CurrentState.SetUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

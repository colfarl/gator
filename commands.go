package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/colfarl/gator/internal/database"
	"github.com/google/uuid"
)

// ========== Command Struct ================

type command struct {
	Name		string
	Args		[]string
}

func argsToCommand(args []string) (command, error) {

	if len(args) < 2 {
		return command{}, fmt.Errorf("USAGE: <program> <command> [args]")
	}
	
	cmdName := args[1]
	
	if len(args) == 2 {
		return command{
					Name: cmdName,
					Args: nil,
				}, nil
	}
	return command{
				Name: cmdName,
				Args: args[2:], 
			}, nil
}

// ============ Commands Struct ================

type commands struct {
	SupportedCommands			map[string]func(*state, command) error
}

func newCommands() commands {
	var c commands
	c.initialize()
	return c
}

func (c *commands) run(s *state, cmd command) error {
	cmdName := cmd.Name
	cmdFunction, ok := c.SupportedCommands[cmdName]
	if !ok {
		return fmt.Errorf("command: %s does not exist", cmdName)
	 }

	err := cmdFunction(s, cmd) 
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.SupportedCommands[name] = f	
}

func (c *commands) initialize() {
	c.SupportedCommands = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
}

// =========== Command Handlers ===============

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) != 1 || cmd.Args == nil {
		return fmt.Errorf("USAGE: login <user-name>")
	}
	
	newUser := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), sql.NullString{String: newUser, Valid: true})
	if  err != nil {
		return fmt.Errorf("user does not exist")
	}

	err = s.updateUser(newUser)
	if err != nil {
		return err
	}

	fmt.Printf("user is now: %s\n", newUser)
	return nil
}

func handlerRegister(s * state, cmd command) error {
	
	if len(cmd.Args) != 1 || cmd.Args == nil {
		return fmt.Errorf("USAGE: register <user-name>")
	}
	
	newName := cmd.Args[0]
	newUser := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{String: newName, Valid: true},
	}

	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return err;
	}

	s.updateUser(newName)
	fmt.Printf("new user %s created\n", newName)
	printUser(user)
	return nil
}

func handlerReset(s * state, cmd command) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("USAGE: reset")
	}
	
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s * state, cmd command) error {	

	if len(cmd.Args) != 0 {
		return fmt.Errorf("USAGE: users")
	}
	
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, v := range names {
		name := v.String
		fmt.Printf(" * %s", name)
		if name == s.CurrentState.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}


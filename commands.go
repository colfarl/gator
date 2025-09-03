package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	c.register("agg", handlerAgg)
	c.register("feeds", handlerFeeds)
	c.register("browse", middlewareLoggedIn(handlerBrowse))
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
}

// ============================== Command Handlers ==============================  

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

func handlerAgg(s * state, cmd command) error {	

	if len(cmd.Args) != 1 {
		return fmt.Errorf("USAGE: agg <time-between-reqs: 1h, 1m, 1s...>")
	}
	
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func handlerFeeds(s *state, cmd command) error {
	
	if len(cmd.Args) != 0 {
		return fmt.Errorf("USAGE: addfeed <feed-name> <feed-url>")
	}
		
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, v := range allFeeds {
		creatorName, err := s.db.GetUserNameByID(context.Background(), v.UserID.UUID)
		if err != nil {
			return err
		}
		
		fmt.Println()
		fmt.Println("Feed Name:", v.Name)
		fmt.Println("URL:", v.Url.String)
		fmt.Println("Creator Name:", creatorName.String)
		fmt.Println()
	}
	return nil
}

// ============================== "LOGGED IN FUNCTIONS" ============================== 
func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("USAGE: follow <url>")
	}
	
		
	feedID, err := s.db.GetFeedIdByURL(context.Background(), sql.NullString{String: cmd.Args[0], Valid: cmd.Args[0] != ""})
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feedID,
		UserID:	user.ID,
	}

	createdFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s; Now following feed: %s", createdFollow.UserName.String, createdFollow.FeedName)
	return nil 
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("USAGE: following")
	}

	allFollowing, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	
	fmt.Println("You are Currently Following:")
	for _, v := range allFollowing {
		fmt.Printf("	- '%s'\n", v.FeedName)
	}

	return nil
}

func handlerAddFeed(s * state, cmd command, user database.User) error {	

	if len(cmd.Args) != 2 {
		return fmt.Errorf("USAGE: addfeed <feed-name> <feed-url>")
	}
	
	params := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
		Url: sql.NullString{String: cmd.Args[1], Valid: cmd.Args[1] != ""},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	}

	inserted, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	
	paramsFollowing := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: inserted.ID,
		UserID:	user.ID,
	}
	
	_, err = s.db.CreateFeedFollow(context.Background(), paramsFollowing)
	if err != nil {
		return err
	}

	fmt.Println("Successfully added feed:", inserted.Name)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	
	if len(cmd.Args) != 1 {
		return fmt.Errorf("USAGE: unfollow <url>")
	}

	feedID, err := s.db.GetFeedIdByURL(context.Background(), sql.NullString{String: cmd.Args[0], Valid: cmd.Args[0] != ""})
	if err != nil {
		return err
	}
	
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)	
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	
	if len(cmd.Args) > 1 {
		return fmt.Errorf("USAGE: browse [limit]")
	}
	
	var limit int32
	if len(cmd.Args) == 1 {
		num, err := strconv.Atoi(cmd.Args[0])
		limit = int32(num)
		if err != nil {
			return err
		}
	} else {
		limit = 2
	}
	
	params := database.GetPostsForUserParams{
		UserID: user.ID, 
		Limit: limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println()
		prettyPost(post)
		fmt.Println()
	}
	
	return nil
}


package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/colfarl/gator/internal/config"
	"github.com/colfarl/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfgInitial, err :=  config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfgInitial.DBURL) 
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}	
	dbQueries := database.New(db)

	cfg := newState(&cfgInitial, dbQueries)		
	cmds := newCommands()	

	cmd, err := argsToCommand(os.Args) 
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	err = cmds.run(&cfg, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

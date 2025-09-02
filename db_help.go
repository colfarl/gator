package main

import (
	"fmt"

	"github.com/colfarl/gator/internal/database"
)

func printUser(u database.User) {
	fmt.Printf("ID: %v\n", u.ID)
	fmt.Printf("Time Created: %v\n", u.CreatedAt)
	fmt.Printf("Time Updated: %v\n", u.UpdatedAt)
	fmt.Printf("Name: %v\n", u.Name.String)
}

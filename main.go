package main

import (
	"fmt"
	"github.com/colfarl/gator/internal/config"
)

func main() {

	cfg, err :=  config.Read()
	if err != nil {
		fmt.Println(err)
		return 
	}

	err = config.SetUser("colin", cfg)
	if err != nil {
		fmt.Println(err)
		return 
	}
	
	
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return 
	}
	
	fmt.Println(cfg)
}

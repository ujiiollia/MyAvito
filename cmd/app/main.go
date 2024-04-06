package main

import (
	"app/internal/config"
	"fmt"
)

func main() {
	// config
	cfg := config.MustLoad()
	fmt.Println(cfg.Address)
	//todo log
	//todo logic
	//todo start serv
}

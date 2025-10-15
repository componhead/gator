package main

import (
	"fmt"

	"github.com/componhead/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Emiliano")
	cfg = config.Read()
	fmt.Println(cfg)
}

package main

import (
	"github.com/madalinpopa/memoflash/internal/api"
)

func main() {
	config := api.NewConfig()
	apiServer := api.NewApiServer(config)
	apiServer.Run()
}

package main

import (
	"github.com/karotte128/karotteapi/v2/api"
	_ "github.com/karotte128/karotteapi/v2/builtins"
)

func main() {
	config := api.Config{}

	api.InitAPI(config)
}

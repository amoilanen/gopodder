package main

import (
	"github.com/amoilanen/gopodder/cmd"
	"github.com/amoilanen/gopodder/pkg/config"
)

func main() {
	config.InitViper()
	cmd.Execute()
}

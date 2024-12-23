package main

import (
	"log"
	"os"

	"github.com/murtaza-u/alertfy/internal/conf"
	"github.com/murtaza-u/alertfy/internal/hook"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Validate()
	if err != nil {
		log.Fatalf("failed to validate provided config: %s", err.Error())
	}
	hook, err := hook.New(*conf)
	if err != nil {
		log.Fatal(err)
	}
	hook.Listen()
}

package main

import (
	"log"
	"os"

	"github.com/murtaza-u/amify/internal/conf"
	"github.com/murtaza-u/amify/internal/hook"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	hook, err := hook.New(*conf)
	if err != nil {
		log.Fatal(err)
	}

	err = hook.Listen()
	if err != nil {
		log.Fatal(err)
	}
}

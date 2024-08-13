package main

import (
	"log"
	"os"

	"github.com/murtaza-u/amify/internal/conf"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(conf)
}

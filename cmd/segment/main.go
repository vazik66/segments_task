package main

import (
	"avito-segment/internal"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	cfg := &internal.Config{}
	if err := cleanenv.ReadConfig("./.env", cfg); err != nil {
		log.Fatalf("Could not read env, %v", err)
	}

	s, err := internal.NewRPCServer(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(s.Run())
}

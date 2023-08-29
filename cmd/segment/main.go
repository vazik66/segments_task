package main

import (
	"avito-segment/internal"
	"context"
	"log"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cfg, err := internal.NewConfig()
    if err != nil {
        log.Fatalln(err)
    }

	s, err := internal.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(s.Run(ctx))
}

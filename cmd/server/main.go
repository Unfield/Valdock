package main

import (
	"fmt"
	"log"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/internal/api"
	"github.com/Unfield/cascade"
)

func main() {
	log.Print("Running Valdock v0.0.1-messing-around")

	cfg := config.ValdockConfig{}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8080

	cfgLoader := cascade.NewLoader(
		cascade.WithEnvPrefix("VALDOCK"),
		cascade.WithFlags(),
	)

	if err := cfgLoader.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	apiRouter := api.NewAPIRouter(&cfg)

	log.Printf("Valdock running on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	if err := apiRouter.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/quard/asciiwrite/internal/rest_api"
	"github.com/quard/asciiwrite/internal/storage"
)

var opts struct {
	Run rest_api.Opts `command:"run"`
}

func main() {
	parser := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	stor, err := storage.NewFirebaseFontStorage()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := rest_api.NewRestAPIServer(opts.Run, stor)
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}

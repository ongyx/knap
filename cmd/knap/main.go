package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/exporter"
	"github.com/ongyx/knap/internal/schema"
)

var L = log.New(os.Stderr, "", 0)

func main() {
	if len(os.Args) < 2 {
		L.Println("usage: knap <path to vault> <path to export ZIP file to>")
		os.Exit(1)
	}

	idn := schema.Identity{
		ID:    uuid.New(),
		Name:  "test",
		Email: "test@test.invalid",
	}

	vp := os.Args[1]
	ex, err := exporter.New(idn, vp)
	if err != nil {
		L.Fatalln(err)
	}

	ep := os.Args[2]
	f, err := os.Create(ep)
	if err != nil {
		L.Fatalln(err)
	}

	if err := ex.Export(f); err != nil {
		L.Fatalln(err)
	}

	L.Println("Export done!")
}

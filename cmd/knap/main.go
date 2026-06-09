package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/schema"
)

var L = log.New(os.Stderr, "", 0)

func main() {
	if len(os.Args) < 2 {
		L.Println("usage: knap <file>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		L.Printf("error: failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	id, _ := uuid.NewRandom()
	d := schema.NewDocument(id)
	if err := d.SetTimestamps(f.Name()); err != nil {
		L.Printf("error: failed to set timestamps from file: %s\n", err)
		os.Exit(1)
	}

	src, err := io.ReadAll(f)
	if err != nil {
		L.Printf("error: failed to read from file: %s\n", err)
		os.Exit(1)
	}

	cv := converter.New(nil)
	d.Data, err = cv.Convert(src)
	if err != nil {
		L.Printf("error: failed to convert markdown: %s\n", err)
		os.Exit(1)
	}

	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		L.Printf("error: failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ongyx/knap/internal/schema"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: knap <file>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error: failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	d := schema.NewDocument()
	if err := d.ParseReader(f); err != nil {
		fmt.Printf("error: failed to parse md: %s\n", err)
		os.Exit(1)
	}

	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		fmt.Printf("error: failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ongyx/knap/internal/collections"
	"github.com/ongyx/knap/internal/exporter"
	"github.com/ongyx/knap/internal/schema"
	"github.com/ongyx/knap/internal/util"
	flag "github.com/spf13/pflag"
)

var (
	// Flag values
	name, email string
	help, force bool
	ignore      []string

	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
)

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintln(o, "Knap exports an Obsidian vault to a collection that can be imported in Outline.")
		fmt.Fprintln(o)
		fmt.Fprintln(o, "usage: knap <path to vault> <path to zipfile>")
		flag.PrintDefaults()
	}
}

func main() {
	flag.BoolVarP(&help, "help", "h", false, "show this help message")
	flag.StringVar(&name, "name", "", "the name to export with")
	flag.StringVar(&email, "email", "", "the email to export with")
	flag.StringSliceVar(&ignore, "ignore", nil, "the directories to ignore in the vault when scanning (can be specified multiple times)")
	flag.BoolVarP(&force, "force", "f", false, "overwrite the zip file if it exists")
	flag.Parse()

	args := flag.Args()
	if help || len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if name == "" || email == "" {
		logger.Println("Warning: Either name or email was not specified. Test values will be used.")
	}

	vaultPath := util.Must(filepath.Abs(args[0]))
	zipPath := util.Must(filepath.Abs(args[1]))

	// Not having the truncate flag here was a big gotcha...
	fl := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if !force {
		fl |= os.O_EXCL
	}
	out, err := os.OpenFile(zipPath, fl, 0666)
	if err != nil {
		logger.Fatalf("Error: Failed to create zipfile: %s\n", err)
	}

	opts := &exporter.ExporterOptions{
		VaultPath: vaultPath,
		Identity: &schema.Identity{
			Name:  name,
			Email: email,
		},
		Logger: logger,
		Ignore: collections.NewSet(ignore...),
	}
	exporter, err := exporter.New(opts)
	if err != nil {
		logger.Fatalf("Error: Failed to init exporter: %s\n", err)
	}

	if err := exporter.Export(out); err != nil {
		logger.Fatalf("Error: Failed to export vault: %s\n", err)
	}

	if err := out.Close(); err != nil {
		logger.Fatalf("Error: Failed to close zipfile: %s\n", err)
	}
}

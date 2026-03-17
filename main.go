package main

import (
	"flag"
	"log"

	"portfolio/internal/builder"
)

func main() {
	dev := flag.Bool("dev", false, "run dev server with live reload on :3000")
	flag.Parse()

	cfg := builder.Config{
		ContentDir: "content",
		OutputDir:  "docs",
	}

	if *dev {
		runDevServer(cfg)
		return
	}

	if err := builder.Build(cfg); err != nil {
		log.Fatal(err)
	}
}

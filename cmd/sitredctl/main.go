package main

import (
	"os"

	"git.sr.ht/~jamesponddotco/sitred/cmd/sitredctl/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args))
}

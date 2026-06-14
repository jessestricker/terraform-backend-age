package main

import (
	"os"

	"github.com/jessestricker/terraform-backend-age/internal"
)

func main() {
	os.Exit(internal.Main())
}

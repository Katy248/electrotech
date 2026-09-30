// Server executable
package main

import (
	"electrotech/internal/di"
	"electrotech/storage"
)

func main() {
	storage.Init(true)

	app := di.NewApp()
	app.Run()
}

// Server executable
package main

import (
	"electrotech/internal/di"
)

func main() {
	app := di.NewApp()
	app.Run()
}

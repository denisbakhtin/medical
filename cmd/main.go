// Package main is used for building main http server executable
package main

import (
	"flag"

	"github.com/denisbakhtin/medical/controllers"
)

func main() {
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	app := controllers.NewApplication(*mode)
	app.StartPeriodicTasks()
	router := app.SetupRouter()

	app.Logger.Infof("Starting http server in %s mode, port: 8010\n", *mode)
	app.Logger.Fatal(router.Run(":8010"))
}

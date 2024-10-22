package main

import (
	"flag"
	"log"

	"github.com/denisbakhtin/medical/controllers"
	"github.com/denisbakhtin/medical/system"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()
	system.Init(*mode)
	system.SetupPeriodicTasks(*mode)

	router := controllers.SetupRouter()
	log.Printf("INFO: Starting http server in %s mode, port: 8010\n", *mode)
	log.Fatal(router.Run(":8010"))
}

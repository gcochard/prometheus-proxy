package proxy

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// Application is a framework for writing proxy prometheus collectors
type Application struct {
	// Extra flags to parse
	Flagset flag.FlagSet

	// Callback called once flags are parsed
	CreateFactory func() CollectorFactory
}

// Main emulates func main() { } for simpler applications
func Main(app Application) {
	app.Flagset = *flag.NewFlagSet("", flag.ContinueOnError)
	port := app.Flagset.Int("port", 2112, "the port to listen on")
	listen := app.Flagset.String("listen", "0.0.0.0", "the address to listen on")
	err := app.Flagset.Parse(os.Args[1:])
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = ListenAndServe(*listen, strconv.Itoa(*port), app.CreateFactory())
	if err != nil {
		log.Fatal(err)
	}
}

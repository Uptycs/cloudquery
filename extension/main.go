package main

import (
	"flag"
	"log"
	"time"

	"github.com/kolide/osquery-go"
)

var (
	socket    = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	keyFile   = flag.String("key-file-path", "", "Path to service account credential file")
	projectId = flag.String("project-id", "", "Project Id")
	zone      = flag.String("zone", "us-east4-c", "zone")
	timeout   = flag.Int("timeout", 3, "Seconds to wait for autoloaded extensions")
	interval  = flag.Int("interval", 3, "Seconds delay between connectivity checks")
)

//go:generate node ./../utilities/extension-codegen/generateGcpExtensions.js ${PWD}
func main() {
	flag.Parse()
	if *socket == "" {
		log.Fatalln("Missing required --socket argument")
	}
	if *keyFile == "" {
		log.Fatalln("Missing required --key-file-path argument")
	}
	if *projectId == "" {
		log.Fatalln("Missing required --project")
	}
	serverTimeout := osquery.ServerTimeout(
		time.Second * time.Duration(*timeout),
	)
	serverPingInterval := osquery.ServerPingInterval(
		time.Second * time.Duration(*interval),
	)

	server, err := osquery.NewExtensionManagerServer(
		"example_extension",
		*socket,
		serverTimeout,
		serverPingInterval,
	)

	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	registerPlugins(server)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

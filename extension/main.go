package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/kolide/osquery-go"
)

var (
	socket        = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	verbose       = flag.Bool("verbose", false, "Enable verbose logging")
	homeDirectory = flag.String("home-directory", "", "Path to the extensions home directory")
	timeout       = flag.Int("timeout", 3, "Seconds to wait for autoloaded extensions")
	interval      = flag.Int("interval", 3, "Seconds delay between connectivity checks")
)

func main() {
	flag.Parse()
	if *socket == "" {
		log.Fatalln("Missing required --socket argument")
	}

	if *homeDirectory == "" {
		// read from environment
		homeDir := os.Getenv("CLOUDQUERY_EXT_HOME")
		if homeDir == "" {
			homeDir = "/etc/cloudquery"
		}
		homeDirectory = &homeDir
	}

	serverTimeout := osquery.ServerTimeout(
		time.Second * time.Duration(*timeout),
	)
	serverPingInterval := osquery.ServerPingInterval(
		time.Second * time.Duration(*interval),
	)

	server, err := osquery.NewExtensionManagerServer(
		"cloudquery_extension",
		*socket,
		serverTimeout,
		serverPingInterval,
	)

	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	readExtensionConfigurations(*homeDirectory + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "extension_config.json")
	readTableConfigurations(*homeDirectory)
	registerPlugins(server)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

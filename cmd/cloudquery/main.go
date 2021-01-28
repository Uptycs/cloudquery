/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Uptycs/cloudquery/extension"
	"github.com/kolide/osquery-go"
)

var (
	socket   = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	timeout  = flag.Int("timeout", 3, "Seconds to wait for autoloaded extensions")
	interval = flag.Int("interval", 3, "Seconds delay between connectivity checks")
)

func main() {
	flag.Parse()
	if *socket == "" {
		log.Fatalln("Missing required --socket argument")
	}

	serverTimeout := osquery.ServerTimeout(
		time.Second * time.Duration(*timeout),
	)
	serverPingInterval := osquery.ServerPingInterval(
		time.Second * time.Duration(*interval),
	)

	homeDirectory := os.Getenv("CLOUDQUERY_EXT_HOME")
	if homeDirectory == "" {
		homeDirectory = "/opt/cloudquery"
	}

	server, err := osquery.NewExtensionManagerServer(
		"cloudquery_extension",
		*socket,
		serverTimeout,
		serverPingInterval,
	)

	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	extension.InitializeLogger(verbose)
	extension.ReadExtensionConfigurations(homeDirectory + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "extension_config.json")
	extension.ReadTableConfigurations(homeDirectory)
	extension.RegisterPlugins(server)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

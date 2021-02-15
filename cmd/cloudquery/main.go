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
	"context"
	"flag"
	"fmt"
	"github.com/Uptycs/cloudquery/utilities"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	osquery "github.com/Uptycs/basequery-go"
	"github.com/Uptycs/cloudquery/extension"
)

var (
	socket   = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	timeout  = flag.Int("timeout", 10, "Seconds to wait for autoloaded extensions")
	interval = flag.Int("interval", 10, "Seconds delay between connectivity checks")
)

func main() {
	flag.Parse()
	if *socket == "" {
		log.Fatalln("Missing required --socket argument")
	}

	homeDirectory := os.Getenv("CLOUDQUERY_EXT_HOME")
	if homeDirectory == "" {
		homeDirectory = "/opt/cloudquery"
	}

	server, err := osquery.NewExtensionManagerServer(
		"cloudquery_extension",
		*socket,
		osquery.ServerVersion("1.0.0"),
		osquery.ServerTimeout(time.Second*time.Duration(*timeout)),
		osquery.ServerPingInterval(time.Second*time.Duration(*interval)),
	)

	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	extension.ReadExtensionConfigurations(homeDirectory+string(os.PathSeparator)+"config"+string(os.PathSeparator)+"extension_config.json", *verbose)
	extension.ReadTableConfigurations(homeDirectory)
	extension.RegisterPlugins(server)

	// Set up cancellation context and waitgroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Wait for interrupt signal to gracefully shutdown the server with waitgroup
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Sprintf("Failed to start extension manager server: %s", err))
		}
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// Start event tables
	for _, eventTable := range extension.GetEventTables() {
		go eventTable.Start(ctx, wg, *socket, time.Second*time.Duration(*timeout))
	}

	<-quit
	utilities.GetLogger().Info("Shutting down cloudquery")

	cancelFunc() // Signal cancellation to context.Context
	// Wait for all thread to exit
	wg.Wait()
	// We are done
	utilities.GetLogger().Info("Graceful shut down done for cloudquery")
}

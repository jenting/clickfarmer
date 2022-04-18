// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"storj.io/clickfarmer/database"
	"storj.io/clickfarmer/webserver"
)

var (
	// Used for flags.
	databaseRPCAddr string
	httpAddr        string
	cacheInterval   time.Duration

	rootCmd = &cobra.Command{
		Use:   "clickfarmer",
		Short: "An application for farming clicks from around the world",
	}

	databaseCmd = &cobra.Command{
		Use:   "database",
		Short: "Run the database RPC server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := database.Run(ctx, databaseRPCAddr)
			if err != nil {
				panic(err)
			}
		},
	}

	webserverCmd = &cobra.Command{
		Use:   "webserver",
		Short: "Run one webserver",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			fmt.Println("running")

			err := webserver.Run(ctx, httpAddr, databaseRPCAddr, "./web", cacheInterval)
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	databaseCmd.Flags().StringVarP(&databaseRPCAddr, "rpc-addr", "r", ":8080", "rpc server address to listen on")

	webserverCmd.Flags().StringVarP(&databaseRPCAddr, "rpc-addr", "r", "localhost:8080", "database rpc server address to connect to")
	webserverCmd.Flags().DurationVarP(&cacheInterval, "cache-interval", "c", 1*time.Second, "cache refresh interval")
	webserverCmd.Flags().StringVarP(&httpAddr, "http-addr", "a", ":3000", "http port to listen on")

	rootCmd.AddCommand(databaseCmd)
	rootCmd.AddCommand(webserverCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

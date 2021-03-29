// root.go - Init CLI and launch the bitcoin vanity adress generator
// Copyright (c) 2021 Tanguy Maraux. Author Tanguy Maraux All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import (
	"fmt"
	"os"
	"time"
	"vanity-btc/addresses"

	"github.com/spf13/cobra"
)

var (
	pattern     string
	threads     int
	verbose     bool
	number      bool
	chronometer bool
)

// main command
var rootCmd = &cobra.Command{
	Use:   "vanity-btc",
	Short: "vanity-btc is a bitcoin vanity address generator",
	Long:  "vanity-btc is a bitcoin vanity address generator.\n\nIt is based on this method : https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addressses",
	Run: func(cmd *cobra.Command, args []string) {
		// launch the program if the pattern is valid
		if addresses.CheckPattern(pattern) {
			// get time before generating addresses
			start := time.Now()

			// generate the bitcoin vanity address
			addresses.GetAddress(pattern, verbose, threads, number)
			// display time duration
			if chronometer {
				duration := time.Since(start)
				fmt.Println("Execution:", duration-time.Millisecond)
			}
		} else {
			fmt.Println("Pattern not valid! It must not contains any '0', 'O', 'l', 'I'")
		}
	},
}

// init command flags
func init() {
	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "pattern wanted in the btc address")
	rootCmd.Flags().IntVarP(&threads, "threads", "t", 4, "number of threads to use")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	rootCmd.Flags().BoolVarP(&number, "number", "n", false, "enable counting the number of addresses generated")
	rootCmd.Flags().BoolVarP(&chronometer, "chronometer", "c", false, "enable chronometer")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

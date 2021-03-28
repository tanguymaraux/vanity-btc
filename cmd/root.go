package cmd

import (
	"fmt"
	"os"
	"time"
	"vanity-btc/addresses"

	"github.com/spf13/cobra"
)

var pattern string
var threads int
var verbose bool
var difficulty bool
var chronometer bool

var rootCmd = &cobra.Command{
	Use:   "vanity-btc",
	Short: "vanity-btc is a personnaly bitcoin address generator",
	Long:  "vanity-btc is a personnaly bitcoin address generator.\n\nIt is based on this method : https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addressses",
	Run: func(cmd *cobra.Command, args []string) {
		if addresses.CheckPattern(pattern) {
			start := time.Now()

			addresses.GetPattern(pattern, verbose, threads, difficulty)
			if chronometer {
				duration := time.Since(start)
				fmt.Println("Execution:", duration)
			}
		} else {
			fmt.Println("Pattern not valid! It must not contains any '0', 'O', 'l', 'I'")
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "pattern wanted in the btc address")
	rootCmd.Flags().IntVarP(&threads, "threads", "t", 4, "number of threads to use")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	rootCmd.Flags().BoolVarP(&difficulty, "difficulty", "d", false, "enable difficulty counting")
	rootCmd.Flags().BoolVarP(&chronometer, "chronometer", "c", false, "enable chronometer")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

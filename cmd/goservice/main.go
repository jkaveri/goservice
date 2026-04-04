package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goservice",
	Short: "A CLI tool for Go service development",
	Long: `A comprehensive CLI tool for Go service development that provides
commands for linting, code formatting, testing, and coverage analysis.`,
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(prettyCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(coverageCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

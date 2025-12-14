/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ackieeee/dnv/internal/kv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "dnv <first-file> <second-file>",
	Short:         "Compare key=value pairs across two files.",
	Long:          `dnv loads KEY=VALUE pairs from the provided files and prints differences in a diff-like format.`,
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		firstPath, secondPath := args[0], args[1]

		firstValues, err := kv.ParseFile(firstPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", firstPath, err)
		}

		secondValues, err := kv.ParseFile(secondPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", secondPath, err)
		}

		result := kv.Compare(firstValues, secondValues)
		if result.IsMatch() {
			cmd.Println("All keys and values match.")
			return nil
		}

		cmd.Printf("--- %s\n", firstPath)
		cmd.Printf("+++ %s\n", secondPath)

		for _, key := range result.MissingInSecond {
			if value, ok := firstValues[key]; ok {
				cmd.Printf("- %s=%s\n", key, value)
				continue
			}
			cmd.Printf("- %s\n", key)
		}

		for _, key := range result.MissingInFirst {
			if value, ok := secondValues[key]; ok {
				cmd.Printf("+ %s=%s\n", key, value)
				continue
			}
			cmd.Printf("+ %s\n", key)
		}

		for _, diff := range result.Differing {
			cmd.Printf("- %s=%s\n", diff.Key, diff.FirstValue)
			cmd.Printf("+ %s=%s\n", diff.Key, diff.SecondValue)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
}

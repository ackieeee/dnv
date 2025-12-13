/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ackieeee/dnv/internal/kv"
	"github.com/spf13/cobra"
)

var dnvCmd = &cobra.Command{
	Use:   "dnv <first-file> <second-file>",
	Short: "Compare key=value pairs across two files",
	Long: `dnv loads KEY=VALUE pairs from the provided files, reports missing keys,
and highlights values that differ between the two sources.`,
	Args: cobra.ExactArgs(2),
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

		for _, key := range result.MissingInSecond {
			cmd.Printf("Missing in %s: %s\n", secondPath, key)
		}

		for _, key := range result.MissingInFirst {
			cmd.Printf("Missing in %s: %s\n", firstPath, key)
		}

		for _, diff := range result.Differing {
			cmd.Printf("Value mismatch for %s: %s=%q %s=%q\n", diff.Key, firstPath, diff.FirstValue, secondPath, diff.SecondValue)
		}

		return fmt.Errorf("found differences between %s and %s", firstPath, secondPath)
	},
}

func init() {
	rootCmd.AddCommand(dnvCmd)

	dnvCmd.SilenceUsage = true
}

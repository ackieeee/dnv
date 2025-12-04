package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ackieeee/dnv/internal/kv"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <file1> <file2>\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Compares key=value pairs between two files.")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	firstPath, secondPath := args[0], args[1]

	firstValues, err := kv.ParseFile(firstPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", firstPath, err)
		os.Exit(1)
	}

	secondValues, err := kv.ParseFile(secondPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", secondPath, err)
		os.Exit(1)
	}

	result := kv.Compare(firstValues, secondValues)
	if result.IsMatch() {
		fmt.Println("All keys and values match.")
		return
	}

	for _, key := range result.MissingInSecond {
		fmt.Printf("Missing in %s: %s\n", secondPath, key)
	}

	for _, key := range result.MissingInFirst {
		fmt.Printf("Missing in %s: %s\n", firstPath, key)
	}

	for _, diff := range result.Differing {
		fmt.Printf("Value mismatch for %s: %s=%q %s=%q\n", diff.Key, firstPath, diff.FirstValue, secondPath, diff.SecondValue)
	}

	os.Exit(1)
}

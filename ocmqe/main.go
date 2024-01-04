package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/openshift-qe/openshift-rosa-cli/cmd/create"
	"github.com/openshift-qe/openshift-rosa-cli/cmd/delete"
	"github.com/openshift-qe/openshift-rosa-cli/cmd/tag"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "ocmqe",
	Short: "Command line tool for ocmqe testing.",
	Long:  "Command line tool for ocm qe to prepare data",
}

func init() {
	// Add the command line flags:
	// fs := root.PersistentFlags()
	// flags.AddDebugFlag(fs)

	// Register the subcommands:
	root.AddCommand(create.Cmd)
	root.AddCommand(delete.Cmd)
	root.AddCommand(tag.Cmd)
}

func main() {
	// Execute the root command:
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "Did you mean this?") {
			fmt.Fprintf(os.Stderr, "Failed to execute root command: %s\n", err)
		}
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TarwsCmd defines the root command
var TarwsCmd = &cobra.Command{
	Use:   "tarws",
	Short: "Stream a Tarball to AWS",
}

// Execute starts the root command
func Execute() {
	err := TarwsCmd.Execute()
	if err != nil {
		fmt.Printf("Error executing tarws %s", err)
		os.Exit(-1)
	}
}

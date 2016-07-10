package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup data",
		Run:   backup,
	}
)

func backup(cmd *cobra.Command, args []string) {
	fmt.Printf("Backupd command %+v\n", args)
}

func init() {
	TarwsCmd.AddCommand(backupCmd)
}

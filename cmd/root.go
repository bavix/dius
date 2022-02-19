package cmd

import (
	"github.com/bavix/dius/internal/du"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "dius [PATH]",
	Short: "Fast calculation of folder/file sizes. Alternative `du -h -d 1`",
	Long: `The command is a replacement for the "du" command for huge volumes.
Speed is achieved using multi-threaded counting (goroutines).`,
	Run: du.Execute,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

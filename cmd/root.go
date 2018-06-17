package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultTagsDB      = "/etc/tags.db"
	DefaultLockTimeout = 10 * time.Second
	BucketName         = "tags"
	AssignmentSymbol   = "="
)

var (
	tagsDB      string
	lockTimeout time.Duration
	rootCmd     = &cobra.Command{
		Use:   "tagrr",
		Short: "tagrr is a tagging tool",
		Long:  `tbd`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&tagsDB, "database", "d", DefaultTagsDB, "tags db file")
	rootCmd.PersistentFlags().DurationVarP(&lockTimeout, "timeout", "t", DefaultLockTimeout, "duration to wait for db to be available")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

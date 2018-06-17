package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultTagsDB      = "/etc/tags.db"
	DefaultLockTimeout = 5 * time.Second
	BucketName         = "tags"
	AssignmentSymbol   = "="
)

var (
	tagsDB      string
	lockTimeout time.Duration
	rootCmd     = &cobra.Command{
		Use:   "tagrr",
		Short: "tagrr is a simple and transactional tags database",
		Long: `Use tagrr to get and set tags (keys with optional values).
Many tagrr processes can concurrently read the db. Only one
tagrr process can write to the db, and not while any other
tagrr processes are reading them.`,
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

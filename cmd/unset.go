package cmd

import (
	"fmt"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/dlipovetsky/tagrr/dbutil"
	"github.com/spf13/cobra"
)

var (
	unsetCmd = &cobra.Command{
		Use:   "unset",
		Short: "unset tags",
		Long:  `Unset tags.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db, err := bolt.Open(tagsDB, 0600, &bolt.Options{Timeout: lockTimeout})
			if err != nil {
				log.Fatalf("Error: failed to open tags db %q: %s\n", tagsDB, err)
			}
			defer db.Close()
			err = UnsetCmd(db, lockTimeout, args)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
		},
	}
)

func UnsetCmd(db *bolt.DB, lockTimeout time.Duration, tags []string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("failed to initialize tags db: %s", err)
		}

		for _, k := range tags {
			err = dbutil.Delete(b, k)
			if err != nil {
				return fmt.Errorf("failed to unset key %q: %s", k, err)
			}
		}
		return nil
	})
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}

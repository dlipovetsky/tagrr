package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/dlipovetsky/tagrr/dbutil"
	"github.com/spf13/cobra"
)

var (
	setCmd = &cobra.Command{
		Use:   "set",
		Short: "set tags",
		Long:  `tbd`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db, err := bolt.Open(tagsDB, 0600, &bolt.Options{Timeout: lockTimeout})
			if err != nil {
				log.Fatalf("Error: failed to open tags db %q: %s\n", tagsDB, err)
			}
			defer db.Close()
			err = SetCmd(db, lockTimeout, args)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
		},
	}
)

func SetCmd(db *bolt.DB, lockTimeout time.Duration, args []string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("failed to initialize tags db %q: %s", tagsDB, err)
		}

		for _, arg := range args {
			k, v, err := parseAssignment(arg)
			if err != nil {
				return fmt.Errorf("failed to parse assignment %q: %s", arg, err)
			}
			err = dbutil.Put(b, k, v)
			if err != nil {
				return fmt.Errorf("failed to set key %q to value %q: %s", k, v, err)
			}
		}
		return nil
	})
}

func parseAssignment(arg string) (key, value string, err error) {
	parsed := strings.Split(arg, AssignmentSymbol)
	if len(parsed) != 2 {
		return "", "", err
	}
	return parsed[0], parsed[1], nil
}

func init() {
	rootCmd.AddCommand(setCmd)
}

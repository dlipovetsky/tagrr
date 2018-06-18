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
	inputTags map[string]string
	setCmd    = &cobra.Command{
		Use:   "set",
		Short: "set tags",
		Long:  `tbd`,
		Args: func(cmd *cobra.Command, args []string) error {
			minNArgs := cobra.MinimumNArgs(1)
			if err := minNArgs(cmd, args); err != nil {
				return err
			}
			inputTags = make(map[string]string)
			for _, arg := range args {
				k, v, err := parseTag(arg)
				if err != nil {
					return err
				}
				inputTags[k] = v
			}
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			db, err := bolt.Open(tagsDB, 0600, &bolt.Options{Timeout: lockTimeout})
			if err != nil {
				log.Fatalf("Error: failed to open tags db %q: %s\n", tagsDB, err)
			}
			defer db.Close()
			err = SetCmd(db, lockTimeout, inputTags)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
		},
	}
)

func SetCmd(db *bolt.DB, lockTimeout time.Duration, tags map[string]string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("failed to initialize tags db: %s", err)
		}

		for k, v := range tags {
			err = dbutil.Put(b, k, v)
			if err != nil {
				return fmt.Errorf("failed to set key %q to value %q: %s", k, v, err)
			}
		}
		return nil
	})
}

func parseTag(arg string) (key, value string, err error) {
	parsed := strings.Split(arg, AssignmentSymbol)
	if len(parsed) != 2 {
		return "", "", fmt.Errorf("failed to parse tag %q, expected %q", arg, fmt.Sprintf("<key>%s<value> or <key>%s", AssignmentSymbol, AssignmentSymbol))
	}
	return parsed[0], parsed[1], nil
}

func init() {
	rootCmd.AddCommand(setCmd)
}

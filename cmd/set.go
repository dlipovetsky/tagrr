package cmd

import (
	"fmt"
	"log"
	"strings"

	bolt "github.com/coreos/bbolt"
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
				log.Fatal("failed to open tags db %q: %s", tagsDB, err)
			}
			defer db.Close()

			err = db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte(BucketName))
				if err != nil {
					log.Fatal("failed to initialize tags db %q: %s", tagsDB, err)
				}

				for _, arg := range args {
					k, v, err := parseAssignment(arg)
					if err != nil {
						fmt.Errorf("failed to parse assignment %q: %s", arg, err)
					}
					err = b.Put([]byte(k), []byte(v))
					if err != nil {
						log.Fatal("failed to update tags db: %s", err)
					}
				}
				return nil
			})
		},
	}
)

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

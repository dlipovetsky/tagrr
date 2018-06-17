package cmd

import (
	"bytes"
	"fmt"
	"log"

	bolt "github.com/coreos/bbolt"
	"github.com/spf13/cobra"
)

const (
	DefaultPrefix = ""
	DefaultAll    = false
)

var (
	prefix string
	all    bool
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get tags",
		Long:  `tbd`,
		Args: func(cmd *cobra.Command, args []string) error {
			if all && len(prefix) > 0 {
				return fmt.Errorf("do not use the `all` and `prefix` flag together")
			}
			if all && len(args) > 0 {
				return fmt.Errorf("do not use keys with the `all` flag")
			}
			if len(prefix) > 0 && len(args) > 0 {
				return fmt.Errorf("do not use keys with the `prefix` flag")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			db, err := bolt.Open(tagsDB, 0600, &bolt.Options{ReadOnly: true, Timeout: lockTimeout})
			if err != nil {
				log.Fatalf("failed to open tags db %q: %s\n", tagsDB, err)
			}
			defer db.Close()

			result := make(map[string]string)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(BucketName))

				// No tags in this database
				if b == nil {
					return nil
				}

				if all {
					b.ForEach(func(k, v []byte) error {
						result[string(k)] = string(v)
						return nil
					})
					return nil
				}

				if len(prefix) > 0 {
					c := b.Cursor()
					for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
						result[string(k)] = string(v)
					}
					return nil
				}

				for _, k := range args {
					if v := b.Get([]byte(k)); v != nil {
						result[k] = string(v)
					}
				}

				return nil
			})

			for k, v := range result {
				fmt.Printf("%s:%s\n", k, v)
			}
		},
	}
)

func init() {
	getCmd.PersistentFlags().BoolVarP(&all, "all", "a", DefaultAll, "tags db file")
	getCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", DefaultPrefix, "tags db file")
	rootCmd.AddCommand(getCmd)
}

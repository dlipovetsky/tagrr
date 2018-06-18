package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/dlipovetsky/tagrr/dbutil"
	"github.com/spf13/cobra"
)

const (
	DefaultPrefix = ""
	DefaultAll    = false
	DefaultFormat = "simple"
)

var AllowedFormats = []string{"simple", "json"}

var (
	prefix string
	all    bool
	format string
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get tags",
		Long:  `tbd`,
		Args: func(cmd *cobra.Command, args []string) error {
			if all && len(prefix) > 0 {
				return fmt.Errorf("must use either the `all` or the `prefix` flag")
			}
			if all && len(args) > 0 {
				return fmt.Errorf("must omit keys when using the `all` flag")
			}
			if len(prefix) > 0 && len(args) > 0 {
				return fmt.Errorf("must omit keys when using the `prefix` flag")
			}
			if !all && len(prefix) == 0 && len(args) == 0 {
				return fmt.Errorf("must specify at least one key")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			db, err := bolt.Open(tagsDB, 0600, &bolt.Options{ReadOnly: true, Timeout: lockTimeout})
			if err != nil {
				log.Fatalf("Error: failed to open tags db %q: %s\n", tagsDB, err)
			}
			defer db.Close()

			var result map[string]string

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(BucketName))

				// No tags in this database
				if b == nil {
					return nil
				}

				if all {
					result = dbutil.GetAll(b)
					return nil
				}

				if len(prefix) > 0 {
					result = dbutil.GetPrefix(b, prefix)
					return nil
				}

				result = dbutil.GetKeys(b, args)
				return nil
			})

			switch format {
			case "json":
				printJSON(os.Stdout, result)
			case "simple":
				printSimple(os.Stdout, result)
			default:
				log.Fatalf("Error: unknown format %q, allowed formats are: %s", format, strings.Join(AllowedFormats, ","))
			}
		},
	}
)

func printSimple(out io.Writer, result map[string]string) {
	for k, v := range result {
		fmt.Fprintf(out, "%s%s%s\n", k, AssignmentSymbol, v)
	}
}

func printJSON(out io.Writer, result map[string]string) {
	j, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatalf("failed to format JSON: %s", err)
	}
	j = append(j, byte('\n'))
	out.Write(j)
}

func init() {
	getCmd.PersistentFlags().BoolVarP(&all, "all", "a", DefaultAll, "tags db file")
	getCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", DefaultPrefix, "tags db file")
	getCmd.PersistentFlags().StringVarP(&format, "format", "o", DefaultFormat, fmt.Sprintf("output format, allowed formats: %s", strings.Join(AllowedFormats, ",")))
	rootCmd.AddCommand(getCmd)
}

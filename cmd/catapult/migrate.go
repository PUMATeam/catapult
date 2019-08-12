package cmd

import (
	log "github.com/sirupsen/logrus"

	migration "github.com/PUMATeam/catapult/internal/database/migration"
	"github.com/spf13/cobra"
)

var reset bool

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "execute migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		argsMig := args[:0]
		for _, arg := range args {
			switch arg {
			case "migrate", "--reset":
			default:
				argsMig = append(argsMig, arg)
			}
		}

		var err error
		if reset {
			err = migration.Reset()
		} else {
			err = migration.Migrate(argsMig)
		}
		if err != nil {
            log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "migrate down to version 0 then up to latest. WARNING: all data will be lost!")

}

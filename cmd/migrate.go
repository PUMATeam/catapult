package cmd

import (
	"log"

	migration "github.com/PUMATeam/catapult/database/migration"
	"github.com/spf13/cobra"
)

var reset bool

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "exectute migrations",
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
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "migrate down to version 0 then up to latest. WARNING: all data will be lost!")

}

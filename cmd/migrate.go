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
		log.Println("migrate called")
		argsMig := args[:0]
		for _, arg := range args {
			switch arg {
			case "migrate", "--reset":
			default:
				argsMig = append(argsMig, arg)
			}
		}

		if reset {
			migration.Reset()
		} else {
			migration.Migrate(argsMig)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&reset, "reset", false, "migrate down to version 0 then up to latest. WARNING: all data will be lost!")

}

package cmd

import (
	"errors"
	"fmt"

	pocket "github.com/brimstone/go-pocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get items from Pocket",
	Long: `Get the items in pocket, in a variety of ways
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		state := pocket.GetStateUnread
		if viper.GetBool("archive") && !viper.GetBool("unread") {
			state = pocket.GetStateArchive
		}
		if viper.GetBool("archive") && viper.GetBool("unread") {
			state = pocket.GetStateAll
		}
		if !viper.GetBool("archive") && !viper.GetBool("unread") {
			return errors.New("Must specify at least archive or unread")
		}

		things, err := p.Get(&pocket.GetOptions{
			State: state,
		})
		if err != nil {
			return err
		}
		for key, thing := range things.List {
			fmt.Printf("%s: %#v\n", key, thing.ResolvedTitle)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	// TODO add a format output type
	/*
		viper.BindEnv("output", "OUTPUT")
		getCmd.Flags().StringP("output", "o", "pretty", "Format of output: pretty, json, csv [$OUTPUT]")
		viper.BindPFlag("output", getCmd.Flags().Lookup("output"))
	*/

	viper.BindEnv("archive", "ARCHIVE")
	getCmd.Flags().BoolP("archive", "a", false, "Return archived items [$ARCHIVE]")
	viper.BindPFlag("archive", getCmd.Flags().Lookup("archive"))

	viper.BindEnv("unread", "UNREAD")
	getCmd.Flags().BoolP("unread", "u", true, "Return unread items [$UNREAD]")
	viper.BindPFlag("unread", getCmd.Flags().Lookup("unread"))
}

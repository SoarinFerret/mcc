package cmd

import (
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/soarinferret/mcc/internal/config"
	"github.com/soarinferret/mcc/internal/meshcentral"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"p"},
	Short:   "Manage local profiles",
	Long:    ``,
}

var profileDefaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"switch", "d"},
	Short:   "Set a new default profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.SetDefaultProfile(args[0], true)
		if err != nil {
			pExit("Failed to switch profile:", err)
		}
		pterm.Info.Println("Switched default profile to: ", args[0])
	},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all profiles",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		printProfileTable(config.GetProfiles())
	},
}

var profileRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove a profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.RemoveProfile(args[0])
		pterm.Info.Println("Removed profile: ", args[0])
	},
}

var profileAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "create"},
	Short:   "Add a new profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		server, _ := cmd.Flags().GetString("server")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		isDefault, _ := cmd.Flags().GetBool("default")
		token, _ := cmd.Flags().GetString("token")
		dontGenerateToken, _ := cmd.Flags().GetBool("dont-generate-token")

		// if password is empty, ask for it
		if password == "" {
			password, _ = pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter your password:")
		}

		if !dontGenerateToken {
			if token != "" {
				meshcentral.SetMfaToken(token)
			}

			// Generate a login token if not using username/password
			err := error(nil)
			username, password, err = meshcentral.GenerateLoginToken(server, username, password)
			if err != nil {
				pExit("Error generating login token: ", err)
			}
		}
		p := config.AddProfile(name, isDefault, server, username, password)

		printProfileTable([]config.Profile{*p})
	},
}

var profileTokenConvertCmd = &cobra.Command{
	Use:     "convert-token",
	Aliases: []string{"ct"},
	Short:   "Convert a username/password profile to a login token profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		p, _ := cmd.Flags().GetString("profile")
		token, _ := cmd.Flags().GetString("token")

		// Get the profile
		profile, err := config.GetProfile(p)
		if err != nil {
			pExit("Error getting profile: ", err)
		}

		if token != "" {
			meshcentral.SetMfaToken(token)
		}

		// Generate a login token
		username, password, err := meshcentral.GenerateLoginToken(profile.Server, profile.Username, profile.Password)
		if err != nil {
			pExit("Error generating login token: ", err)
		}

		// Update the profile with the new token
		profile.Username = username
		profile.Password = password
		err = config.UpdateProfile(*profile)
		if err != nil {
			pExit("Error updating profile: ", err)
		}

		pterm.Info.Println("Converted profile to use login token successfully.")
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileDefaultCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRmCmd)
	profileCmd.AddCommand(profileTokenConvertCmd)

	profileAddCmd.Flags().StringP("name", "n", "", "The name of the profile to add")
	profileAddCmd.Flags().BoolP("default", "d", false, "Set this profile as the default profile")
	profileAddCmd.Flags().StringP("server", "s", "", "Mesh Central Server URL")
	profileAddCmd.Flags().StringP("username", "u", "", "Mesh Central Username")
	profileAddCmd.Flags().StringP("password", "p", "", "Mesh Central Password")
	profileAddCmd.Flags().BoolP("dont-generate-token", "", false, "Don't replace username and password with a login token")
	profileAddCmd.Flags().StringP("token", "t", "", "MFA Token for Mesh Central")
	profileTokenConvertCmd.Flags().StringP("profile", "p", "", "The profile to convert to a login token")
	profileTokenConvertCmd.Flags().StringP("token", "t", "", "MFA Token for Mesh Central")
	profileAddCmd.MarkFlagRequired("name")
	profileAddCmd.MarkFlagRequired("server")
	profileAddCmd.MarkFlagRequired("username")
	//profileAddCmd.MarkFlagRequired("password")

}

func printProfileTable(profiles []config.Profile) {
	// print profiles in a table
	profileData := [][]string{}

	// add header
	profileData = append(profileData, []string{"Name", "Server", "Username", "IsDefault"})

	// add profile data
	for _, p := range profiles {
		d := config.GetDefaultProfileName()
		/*isDefault := "false"
		if strings.Compare(p.Name, d) == 0 {
			isDefault = "true"
			}*/

		profileData = append(
			profileData,
			[]string{
				p.Name,
				p.Server,
				p.Username,
				strconv.FormatBool((strings.Compare(p.Name, d) == 0)),
			},
		)
	}
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(profileData).Render()
}

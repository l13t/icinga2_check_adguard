package cmd

import (
	"github.com/spf13/cobra"

	"icinga2_check_adguard/checkadguard"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "check_adguard",
	Short: "CLI user interface to work with CAN bus",
	Long:  "This command starts CLI user interface to work with CAN bus",
	Run: func(cmd *cobra.Command, args []string) {
		if argsErr := cobra.NoArgs(cmd, args); argsErr != nil {
			cmd.Help()
			os.Exit(0)
		}
		host, _ = cmd.Flags().GetString("host")
		port, _ = cmd.Flags().GetString("port")
		username, _ = cmd.Flags().GetString("username")
		password, _ = cmd.Flags().GetString("password")
		metrics, _ = cmd.Flags().GetBool("metrics")
		mode, _ = cmd.Flags().GetString("mode")
		ssl, _ = cmd.Flags().GetBool("ssl")
		insecure, _ = cmd.Flags().GetBool("insecure")
		timeout, _ = cmd.Flags().GetInt("timeout")

		// if host == "" {
		// 	host = "localhost"
		// }
		// if port == "" {
		// 	port = "3000"
		// }

		switch check_adguard.Mode {
		case "simple",
			"detailed":
			check_adguard.CheckAdGuard(host, port, username, password, mode, timeout, ssl, insecure, metrics)
		default:
			fmt.Printf("Only `simple` or `detailed` modes are supported.")
			os.Exit(0)
		}
	},
}

// Execute will check input variables, run required calls to AdGuard API and return status of the AdGuard.
// If it's required by input variable `metrics`, it would generate Nagios compatible metrics in status output.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().String("host", "localhost", "Host address of the AdGuard API endpoint")
	rootCmd.PersistentFlags().String("port", "3000", "Port number for AdGuard API endpoint")
	rootCmd.PersistentFlags().String("username", "admin", "Username used to do API calls (by default: `admin`)")
	rootCmd.PersistentFlags().String("password", "", "Password to access API endpoint (by default: empty)")
	rootCmd.PersistentFlags().Bool("metrics", false, "Defines if we need to enable metrics in plugin output (by default: False)")
	rootCmd.PersistentFlags().String("mode", "simple", "Defines how detailed would be output for a check")
	rootCmd.PersistentFlags().Bool("ssl", false, "Defines if AdGuard uses HTTPS on API side")
	rootCmd.PersistentFlags().Bool("insecure", false, "If set to `true`, does not validate certificate")
	rootCmd.PersistentFlags().Bool("timeout", 5, "Timeout for http request")
}

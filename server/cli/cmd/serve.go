package cmd

import (
	"github.com/bahadrix/corpushub/server/operator"
	"github.com/bahadrix/corpushub/server/rest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start rest server",
	Long: `Starts rest server at given host and port.`,
	Run: func(cmd *cobra.Command, args []string) {

		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		debugMode, _ := cmd.Flags().GetBool("debug")


		op, err := operator.NewOperator(DataPath, nil)

		if err != nil {
			log.Fatal("Can not create operator", err)
		}

		err = rest.Start(host, port, debugMode, op)

		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("host","0.0.0.0", "Defaults to all interfaces")
	serveCmd.Flags().Int("port", 8081, "Port to listen")
	serveCmd.Flags().BoolP("debug", "d", false, "Set to enable debug mode")


}

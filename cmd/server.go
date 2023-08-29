package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vamika-digital/wms-api-server/config"
	"github.com/vamika-digital/wms-api-server/pkg/database"
	"github.com/vamika-digital/wms-api-server/pkg/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web application server",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a database connection
		dbConn, err := database.NewMySQLConnection(config.AppConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %s", err)
			return
		}
		// Create a new server instance with the provided port and database connection
		srv := server.NewServer(config.AppConfig.RestServer.Address, config.AppConfig.RestServer.Port, dbConn)
		// Run the server
		srv.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

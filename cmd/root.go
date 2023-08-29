package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vamika-digital/wms-api-server/config"
)

var rootCmd = &cobra.Command{
	Use:   "WMS",
	Short: "Warehouse Managment System",
	Long:  "Warehouse Managment System",
}

func init() {
	cobra.OnInitialize(config.InitConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

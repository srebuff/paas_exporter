package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srebuff/paas_exporter/paas"
	"log"
	"strings"
)

var (
	dsn_mongo string
	dsl_mongo string
)
var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "mongo相关操作",
	Long:  "mongo ping等操作",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.Contains(dsn, "mongo") {
			err := paas.ConnectMongoDB(dsn)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println("connect mongodb success")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mongoCmd)
	mongoCmd.PersistentFlags().StringVarP(&dsn_mongo, "dsn", "d", "", "mongodb dsn")
	mongoCmd.PersistentFlags().StringVarP(&dsl_mongo, "sql", "s", "", "mongodb dsl statement")
}

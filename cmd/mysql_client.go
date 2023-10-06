package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srebuff/paas_exporter/paas"
	"log"
	"strings"
)

var (
	dsn string
	sql string
)
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql相关操作",
	Long:  "mysql ping等操作",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.Contains(dsn, "mysql") {
			//remove dsn prefix msysql://
			dsn = strings.TrimPrefix(dsn, "mysql://")
			err := paas.ConnectMysql(dsn)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println("connect mysql success")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", "", "database dsn")
	mysqlCmd.PersistentFlags().StringVarP(&sql, "sql", "s", "", "sql statement")
}

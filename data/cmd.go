//nolint:all
package data

import (
	"github.com/spf13/cobra"
)

// CmdData the service command.
var CmdData = &cobra.Command{
	Use:   "data",
	Short: "generate the kratos data code from proto",
	Long:  "generate the kratos data code from proto",
	Run:   run,
}

var (
	db                string // input mysql or postgres
	dsn               string // 数据库连接
	targetTables      string // 指定表
	outPutDataPath    string // data输出路径
	outPutServicePath string // service输出路径
)

func init() {
	CmdData.Flags().StringVarP(&db, "db", "", "", "db")
	CmdData.Flags().StringVarP(&dsn, "dsn", "", "", "dsn")
	CmdData.Flags().StringVarP(&targetTables, "tables", "", "", "tables")
	CmdData.Flags().StringVarP(&outPutDataPath, "outPutDataPath", "", "./internal/data", "outPutDataPath")
	CmdData.Flags().StringVarP(&outPutServicePath, "outPutServicePath", "", "./internal/service", "outPutServicePath")
}

func run(_ *cobra.Command, args []string) {
	NewData(db, dsn, targetTables, outPutDataPath, outPutServicePath).Run()
}

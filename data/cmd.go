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
	db             string // input mysql or postgres
	dsn            string // 数据库连接
	targetTables   string // 指定表
	outPutDataPath string // data输出路径
)

func init() {
	CmdData.Flags().StringVarP(&db, "db", "d", "", "db")
	CmdData.Flags().StringVarP(&dsn, "dsn", "s", "", "dsn")
	CmdData.Flags().StringVarP(&targetTables, "tables", "t", "", "tables")
	CmdData.Flags().StringVarP(&outPutDataPath, "outPutDataPath", "", "./internal/data", "outPutDataPath")
}

func run(_ *cobra.Command, args []string) {
	NewData(db, dsn, targetTables, outPutDataPath).Run()
}

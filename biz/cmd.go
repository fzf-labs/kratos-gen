package biz

import (
	"github.com/spf13/cobra"
)

// CmdLogic the service command.
var CmdBiz = &cobra.Command{
	Use:   "biz",
	Short: "generate the kratos service biz code from proto",
	Long:  "generate the kratos service biz code from proto",
	Run:   run,
}

var (
	inPutPbPath   string // pb输入路径
	outPutBizPath string // biz输出路径
)

//nolint:gochecknoinits
func init() {
	CmdBiz.Flags().StringVarP(&inPutPbPath, "inPutPbPath", "i", "./api", "inPutPbPath")
	CmdBiz.Flags().StringVarP(&outPutBizPath, "outPutBizPath", "b", "./internal/biz", "outPutBizPath")
}

// run the service command.
func run(_ *cobra.Command, _ []string) {
	NewBiz(inPutPbPath, outPutBizPath).Run()
}

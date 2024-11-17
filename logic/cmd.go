package logic

import (
	"github.com/spf13/cobra"
)

// CmdLogic the service command.
var CmdLogic = &cobra.Command{
	Use:   "logic",
	Short: "generate the kratos service biz code from proto",
	Long:  "generate the kratos service biz code from proto",
	Run:   run,
}

var (
	inPutPbPath       string // pb输入路径
	outPutBizPath     string // biz输出路径
	outPutServicePath string // service输出路径
)

//nolint:gochecknoinits
func init() {
	CmdLogic.Flags().StringVarP(&inPutPbPath, "inPutPbPath", "i", "./api", "inPutPbPath")
	CmdLogic.Flags().StringVarP(&outPutBizPath, "outPutBizPath", "b", "./internal/biz", "outPutBizPath")
	CmdLogic.Flags().StringVarP(&outPutServicePath, "outPutServicePath", "s", "./internal/service", "outPutServicePath")
}

// run the service command.
func run(_ *cobra.Command, _ []string) {
	NewLogic(inPutPbPath, outPutBizPath, outPutServicePath).Run()
}

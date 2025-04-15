package service

import (
	"github.com/spf13/cobra"
)

// CmdLogic the service command.
var CmdService = &cobra.Command{
	Use:   "service",
	Short: "generate the kratos service biz code from proto",
	Long:  "generate the kratos service biz code from proto",
	Run:   run,
}

var (
	inPutPbPath       string // pb输入路径
	outPutServicePath string // service输出路径
)

//nolint:gochecknoinits
func init() {
	CmdService.Flags().StringVarP(&inPutPbPath, "inPutPbPath", "i", "./api", "inPutPbPath")
	CmdService.Flags().StringVarP(&outPutServicePath, "outPutServicePath", "s", "./internal/service", "outPutServicePath")
}

// run the service command.
func run(_ *cobra.Command, _ []string) {
	NewService(inPutPbPath, outPutServicePath).Run()
}

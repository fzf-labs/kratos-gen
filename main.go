package main

import (
	"log"

	"github.com/fzf-labs/kratos-gen/data"
	"github.com/fzf-labs/kratos-gen/logic"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kratos-gen",
	Short:   "kratos-gen: gen service biz data code",
	Long:    `kratos-gen: gen service biz data code`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(logic.CmdLogic)
	rootCmd.AddCommand(data.CmdData)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

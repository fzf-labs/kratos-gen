package main

import (
	"log"

	"github.com/fzf-labs/kratos-gen/data"
	"github.com/fzf-labs/kratos-gen/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kratos-gen",
	Short:   "kratos-gen: gen service biz data code",
	Long:    `kratos-gen: gen service biz data code`,
	Version: "v2.8.2",
}

func init() {
	rootCmd.AddCommand(data.CmdData)
	rootCmd.AddCommand(service.CmdService)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

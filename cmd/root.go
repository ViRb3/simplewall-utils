package cmd

import (
	"github.com/spf13/cobra"
	"simplewall-utils/cmd/allow"
)

var cmd = &cobra.Command{
	Use:   "simplewall-utils",
	Short: "Simple utilities for simplewall",
}

func Execute() error {
	return cmd.Execute()
}

func init() {
	cmd.AddCommand(allow.Cmd)
}

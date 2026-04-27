package flags

import (
	"github.com/spf13/cobra"
)

var v3ApiFlag bool

func AddV3ApiFlag(cmd *cobra.Command) {
	usage := "If set, use V3 api when possible."
	cmd.PersistentFlags().BoolVar(&v3ApiFlag, "v3-api", false, usage)
	cmd.PersistentFlags().MarkHidden("v3-api")
}

func V3Api() bool {
	return v3ApiFlag
}

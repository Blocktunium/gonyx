package cli

import (
	"github.com/abolfazlbeh/zhycan/internal/command"
	"github.com/spf13/cobra"
)

// AttachCommands - Attach default cli commands
func AttachCommands(cmd *cobra.Command) {
	cmd.AddCommand(command.NewRunServerCmd()) // Run Server Command
}

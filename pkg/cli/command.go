package cli

import (
	"github.com/Blocktunium/gonyx/internal/command"
	"github.com/spf13/cobra"
)

// AttachCommands - Attach default cli commands
func AttachCommands(cmd *cobra.Command) {
	cmd.AddCommand(command.NewRunServerCmd())      // Run Server Command
	cmd.AddCommand(command.NewCompileCommandCmd()) // Compile protobuf Command
}

package command

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Blocktunium/gonyx/internal/cache"
	"github.com/Blocktunium/gonyx/internal/grpc"
	"github.com/Blocktunium/gonyx/internal/http"
	"github.com/spf13/cobra"
)

const (
	RunServerInitMsg     = `Gonyx > Running Server ...`
	RunServerShutdownMsg = `Gonyx > Shutting Down Server ...`
)

func NewRunServerCmd() *cobra.Command {
	runServerCmd := &cobra.Command{
		Use:   "runserver",
		Short: "Run Restfull Server And Other Engine If Existed",
		Long:  ``,

		Run:  runServerCmdExecute,
		RunE: runServerCmdExecuteE,
	}

	// Add server-type flag with short name
	runServerCmd.Flags().StringP("server-type", "s", "", "Type of server to run (http, grpc)")

	return runServerCmd
}

func runServerCmdExecuteE(cmd *cobra.Command, args []string) error {
	runServerCmdExecute(cmd, args)
	return nil
}

func runServerCmdExecute(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), RunServerInitMsg)
	m := cache.GetManager()

	// Get the server-type flag value
	serverType, _ := cmd.Flags().GetString("server-type")
	serverType = strings.ToLower(serverType)

	// Start the appropriate servers based on the flag value
	if serverType == "" || serverType == "http" {
		http.GetManager().StartServers()
	}

	if serverType == "" || serverType == "grpc" {
		grpc.GetManager().StartServers()
	}

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Fprintf(cmd.OutOrStdout(), RunServerShutdownMsg)

	// Stop only the servers that were started
	if serverType == "" || serverType == "http" {
		http.GetManager().StopServers()
	}

	if serverType == "" || serverType == "grpc" {
		grpc.GetManager().StopServers()
	}

	err := m.Release()
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), err.Error())
	}

	//var wg sync.WaitGroup
	//wg.Add(1)
	//
	//go func() {
	//	quit := make(chan os.Signal)
	//	// kill (no param) default send syscall.SIGTERM
	//	// kill -2 is syscall.SIGINT
	//	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//	<-quit
	//
	//	wg.Done()
	//}()
	//
	//wg.Wait()
}

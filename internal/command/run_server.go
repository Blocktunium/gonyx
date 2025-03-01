package command

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/cache"
	"github.com/Blocktunium/gonyx/internal/grpc"
	"github.com/Blocktunium/gonyx/internal/http"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

const (
	RunServerInitMsg     = `Zhycan > Running Server ...`
	RunServerShutdownMsg = `Zhycan > Shutting Down Server ...`
)

func NewRunServerCmd() *cobra.Command {
	runServerCmd := &cobra.Command{
		Use:   "runserver",
		Short: "Run Restfull Server And Other Engine If Existed",
		Long:  ``,

		Run:  runServerCmdExecute,
		RunE: runServerCmdExecuteE,
	}
	return runServerCmd
}

func runServerCmdExecuteE(cmd *cobra.Command, args []string) error {
	runServerCmdExecute(cmd, args)
	return nil
}

func runServerCmdExecute(cmd *cobra.Command, args []string) {
	// TODO: in future 'args' must be considered
	fmt.Fprintf(cmd.OutOrStdout(), RunServerInitMsg)
	m := cache.GetManager()

	http.GetManager().StartServers()
	grpc.GetManager().StartServers()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Fprintf(cmd.OutOrStdout(), RunServerShutdownMsg)

	http.GetManager().StopServers()
	grpc.GetManager().StopServers()
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

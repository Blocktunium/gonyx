package command

import (
	"fmt"
	"github.com/abolfazlbeh/zhycan/internal/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	RunCompileCommandInitMsg      = `Zhycan > Compiling the protobuf file ...`
	RunCompileCommandFileName     = `Zhycan > "%s.proto" file is going to be compiled ...`
	RunCompileCommandError        = `Zhycan > Compiling encountered the error: %s`
	RunCompileCommandFileNotExist = `Zhycan > %s.proto file does not exist in "%s"`
)

func NewCompileCommandCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "compile [protobuf_name]",
		Short: `Create a protobuf file with the ".proto" extension`,
		Long:  `This command create a directory with the name of protobuf filename, and then compile it to that folder`,
		Run:   runCompileCmdExecute,
		RunE:  runCompileCmdExecuteE,
	}
	return createCmd
}

func runCompileCmdExecuteE(cmd *cobra.Command, args []string) error {
	runCompileCmdExecute(cmd, args)
	return nil
}

func runCompileCmdExecute(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandInitMsg+"\n")
	fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandFileName+"\n", args[0])

	ex, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandError, err.Error())
		return
	}

	// Let's create a path for file and the compiled folder
	// Get the appropriate "proto" folder path from the config

	protoFolderPath, err := config.GetManager().Get("protobuf", "src_dir")
	if err != nil {
		protoFolderPath = "app/proto"
	}

	desiredProtoPath := filepath.Join(ex, protoFolderPath.(string))

	// Let's find the file
	protobufFilePath := filepath.Join(desiredProtoPath, args[0]+".proto")
	if _, err1 := os.Stat(protobufFilePath); os.IsNotExist(err1) {
		// file does not exist --> Show error and exit
		fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandFileNotExist, args[0]+".proto", desiredProtoPath)
		return
	}

	// File exist --> go to create the folder
	compiledFolderPath := filepath.Join(desiredProtoPath, args[0])
	err2 := os.Mkdir(compiledFolderPath, os.ModePerm)
	if err2 != nil {
		// Folder existed --> Let's empty it
	}

	// protoc --go_out=./exchange --go_opt=paths=source_relative --go-grpc_out=./exchange --go-grpc_opt=paths=source_relative exchange.proto
	//commandToBeExecuted := fmt.Sprintf("--go_out=%s --go_opt=paths=source_relative --go-grpc_out=%s --go-grpc_opt=paths=source_relative %s",
	//	compiledFolderPath, compiledFolderPath, protobufFilePath)

	// changing the directory to the protobuf folder
	err3 := os.Chdir(desiredProtoPath)
	if err3 != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandError, err3.Error())
		return
	}

	protoCmd := exec.Command("protoc")
	protoCmd.Args = append(protoCmd.Args, fmt.Sprintf("--go_out=./%s", args[0]))
	protoCmd.Args = append(protoCmd.Args, fmt.Sprintf("--go_opt=paths=source_relative"))
	protoCmd.Args = append(protoCmd.Args, fmt.Sprintf("--go-grpc_out=./%s", args[0]))
	protoCmd.Args = append(protoCmd.Args, fmt.Sprintf("--go-grpc_opt=paths=source_relative"))
	protoCmd.Args = append(protoCmd.Args, fmt.Sprintf("%s.proto", args[0]))

	protocErr := protoCmd.Start()
	if protocErr != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandError, protocErr)
		return
	}

	protocErr = protoCmd.Wait()
	if protocErr != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), RunCompileCommandError, protocErr)
		return
	}
}

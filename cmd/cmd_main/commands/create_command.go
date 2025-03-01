package commands

import (
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"text/template"
	"time"
)

const (
	InitializeCreateCommandMessage       = `Zhycan > Create A New Command`
	CannotGetCurrentDirectoryMessage     = `Zhycan > Cannot current working directory ... %v`
	CannotCreateCommandsDirectoryMessage = `Zhycan > Cannot create commands directory ... %v`
	CannotCreateCommandFileMessage       = `Zhycan > Cannot create command file: %v ... %v`
	CannotFillCommandFileMessage         = `Zhycan > Cannot fill command file: %v ... %v`
	CommandFileCreatedMessage            = `Zhycan > Command file is created: %v`
)

const (
	CommandDirectory = "commands"
)

func NewCreateCommandCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "command [command_name]",
		Short: `Create A New Command With Specified Name`,
		Long:  ``,
		Run:   createCommandCmdExecute,
		RunE:  createCommandCmdExecuteE,
	}
	return createCmd
}

func createCommandCmdExecuteE(cmd *cobra.Command, args []string) error {
	createCommandCmdExecute(cmd, args)
	return nil
}

func createCommandCmdExecute(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), InitializeCreateCommandMessage)

	myDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), CannotGetCurrentDirectoryMessage, err)
		return
	}

	commandDir := filepath.Join(myDir, CommandDirectory)
	if _, err := os.Stat(commandDir); errors.Is(err, os.ErrNotExist) {
		// create commands directory
		createErr := os.MkdirAll(CommandDirectory, os.ModePerm)
		if createErr != nil {
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintf(cmd.OutOrStdout(), CannotCreateCommandsDirectoryMessage, createErr)
			return
		}
	}

	// get command name from args
	commandName := args[0]

	currentUser, err := user.Current()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), UserNotExisted, err)
		return
	} else {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), UserExisted)
	}

	// get current year
	year := time.Now().Year()

	// `commands` directory is existed --> create command file first

	projectName := config.GetManager().GetName()
	createFileErr := createCommandFile(cmd, projectName, commandName,
		currentUser.Username, year)
	if createFileErr != nil {
		return
	}
}

func createCommandFile(cmd *cobra.Command, projectName, commandName, creatorUsername string, year int) error {
	filePath := filepath.Join(CommandDirectory, fmt.Sprintf("%s_command.go", commandName))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), CannotCreateCommandFileMessage, filePath, err)
		return err
	}
	defer file.Close()

	temp := template.Must(template.ParseFiles("./templates/root_command.gotmpl"))
	goModuleVars := struct {
		ProjectName     string
		CreatorUserName string
		Time            time.Time
		TimeFormat      string
		Year            int
		CommandName     string
	}{
		ProjectName:     projectName,
		CreatorUserName: creatorUsername,
		Time:            time.Now().Local(),
		TimeFormat:      time.RFC822,
		Year:            year,
		CommandName:     commandName,
	}
	err = temp.Execute(file, goModuleVars)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), CannotFillCommandFileMessage, filePath, err)
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), RootCommandGoFileIsCreated)

	return nil
}

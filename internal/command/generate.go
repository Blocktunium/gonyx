package command

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	// Generate command constants
	GenerateSwaggerStartMessage        = `Gonyx > Generating Swagger/OpenAPI documentation ...`
	GenerateSwaggerCompleteMessage     = `Gonyx > Swagger documentation generated successfully in "%s" ...`
	GenerateSwaggerErrorMessage        = `Gonyx > Error generating Swagger documentation ... %v`
	GenerateSwaggerOutputDirMessage    = `Gonyx > Output directory: %s`
	GenerateSwaggerCheckingMessage     = `Gonyx > Checking swag installation...`
	GenerateSwaggerNotFoundMessage     = `Gonyx > Swag not found, installing...`
	GenerateSwaggerInstallSuccessMsg   = `Gonyx > Swag installed successfully`
	GenerateSwaggerAlreadyInstalledMsg = `Gonyx > Swag is already installed`
	GenerateSwaggerInitRunningMsg      = `Gonyx > Running 'swag init'...`
	GenerateSwaggerInitCompleteMsg     = `Gonyx > 'swag init' completed successfully`
	GenerateSwaggerFmtRunningMsg       = `Gonyx > Running 'swag fmt'...`
	GenerateSwaggerFmtCompleteMsg      = `Gonyx > 'swag fmt' completed successfully`

	// Error messages
	GenerateSwaggerInstallFailMsg      = `Gonyx > Failed to install swag: %v`
	GenerateSwaggerInitFailMsg         = `Gonyx > Failed to run 'swag init': %v`
	GenerateSwaggerFmtFailMsg          = `Gonyx > Failed to run 'swag fmt': %v`
	GenerateSwaggerInstallStartFailMsg = `Gonyx > Failed to start install command: %v`
	GenerateSwaggerInstallCmdFailMsg   = `Gonyx > Install command failed: %v, stderr: %s`
	GenerateSwaggerInitStartFailMsg    = `Gonyx > Failed to start 'swag init': %v`
	GenerateSwaggerInitCmdFailMsg      = `Gonyx > 'swag init' failed: %v, stderr: %s`
	GenerateSwaggerFmtStartFailMsg     = `Gonyx > Failed to start 'swag fmt': %v`
	GenerateSwaggerFmtCmdFailMsg       = `Gonyx > 'swag fmt' failed: %v, stderr: %s`
)

// NewGenerateCmd creates the main generate command
func NewGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code scaffolds and documentation for your application",
		Long: `The generate command provides various sub-commands to automatically generate
code scaffolds, documentation, and configuration files for your Gonyx application.`,

		// This command requires a subcommand
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add subcommands
	generateCmd.AddCommand(NewGenerateSwaggerCmd())

	return generateCmd
}

// NewGenerateSwaggerCmd creates the swagger subcommand
func NewGenerateSwaggerCmd() *cobra.Command {
	swaggerCmd := &cobra.Command{
		Use:   "swagger",
		Short: "Generate Swagger/OpenAPI documentation for your application",
		Long: `Generate Swagger/OpenAPI documentation files for your Gonyx application.
This command will scan your application code and generate comprehensive API documentation.`,

		Run:  generateSwaggerExecute,
		RunE: generateSwaggerExecuteE,
	}

	return swaggerCmd
}

// generateSwaggerExecuteE is the error wrapper for swagger command
func generateSwaggerExecuteE(cmd *cobra.Command, args []string) error {
	generateSwaggerExecute(cmd, args)
	return nil
}

// generateSwaggerExecute is the main execution function for swagger command
func generateSwaggerExecute(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerStartMessage)

	// Get output directory flag
	outputDir := "./docs"

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerOutputDirMessage, outputDir)

	// TODO: Implement actual swagger generation logic
	err := generateSwaggerDocumentation(cmd, outputDir)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerErrorMessage, err)
		return
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerCompleteMessage, outputDir)
}

// generateSwaggerDocumentation generates the actual swagger documentation
func generateSwaggerDocumentation(cmd *cobra.Command, outputDir string) error {
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerCheckingMessage)

	// Step 1: Check if swag is installed
	err := checkSwagInstallation(cmd)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerNotFoundMessage)

		// Install swag
		err = installSwag(cmd)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr())
			fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInstallFailMsg, err)
			return fmt.Errorf("failed to install swag: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerInstallSuccessMsg)
	} else {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerAlreadyInstalledMsg)
	}

	// Step 2: Run swag init
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerInitRunningMsg)

	err = runSwagInit(cmd)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInitFailMsg, err)
		return fmt.Errorf("failed to run 'swag init': %v", err)
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerInitCompleteMsg)

	// Step 3: Run swag fmt
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerFmtRunningMsg)

	err = runSwagFmt(cmd)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerFmtFailMsg, err)
		return fmt.Errorf("failed to run 'swag fmt': %v", err)
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintf(cmd.OutOrStdout(), GenerateSwaggerFmtCompleteMsg)

	return nil
}

// checkSwagInstallation checks if swag command is available
func checkSwagInstallation(cmd *cobra.Command) error {
	checkCmd := exec.Command("swag", "version")
	var stderr bytes.Buffer
	checkCmd.Stderr = &stderr

	err := checkCmd.Run()
	return err
}

// installSwag installs the swag command using go install
func installSwag(cmd *cobra.Command) error {
	installCmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
	var stderr bytes.Buffer
	installCmd.Stderr = &stderr

	err := installCmd.Start()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInstallStartFailMsg, err)
		return fmt.Errorf("failed to start install command: %v", err)
	}

	err = installCmd.Wait()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInstallCmdFailMsg, err, stderr.String())
		return fmt.Errorf("install command failed: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// runSwagInit runs the swag init command
func runSwagInit(cmd *cobra.Command) error {
	swagCmd := exec.Command("swag", "init")
	var stderr bytes.Buffer
	swagCmd.Stderr = &stderr

	err := swagCmd.Start()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInitStartFailMsg, err)
		return fmt.Errorf("failed to start 'swag init': %v", err)
	}

	err = swagCmd.Wait()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerInitCmdFailMsg, err, stderr.String())
		return fmt.Errorf("'swag init' failed: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// runSwagFmt runs the swag fmt command
func runSwagFmt(cmd *cobra.Command) error {
	swagCmd := exec.Command("swag", "fmt")
	var stderr bytes.Buffer
	swagCmd.Stderr = &stderr

	err := swagCmd.Start()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerFmtStartFailMsg, err)
		return fmt.Errorf("failed to start 'swag fmt': %v", err)
	}

	err = swagCmd.Wait()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr())
		fmt.Fprintf(cmd.OutOrStderr(), GenerateSwaggerFmtCmdFailMsg, err, stderr.String())
		return fmt.Errorf("'swag fmt' failed: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

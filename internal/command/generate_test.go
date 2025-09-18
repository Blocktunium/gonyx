package command

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerateCmd(t *testing.T) {
	cmd := NewGenerateCmd()

	assert.Equal(t, "generate", cmd.Use)
	assert.Equal(t, "Generate code scaffolds and documentation for your application", cmd.Short)
	assert.Contains(t, cmd.Long, "The generate command provides various sub-commands")
	assert.True(t, len(cmd.Commands()) > 0, "Generate command should have subcommands")
}

func TestNewGenerateSwaggerCmd(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	assert.Equal(t, "swagger", cmd.Use)
	assert.Equal(t, "Generate Swagger/OpenAPI documentation for your application", cmd.Short)
	assert.Contains(t, cmd.Long, "Generate Swagger/OpenAPI documentation files")
}

func TestGenerateCommand_ShowsHelpWhenNoSubcommand(t *testing.T) {
	cmd := NewGenerateCmd()

	// Capture output
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	// Execute without subcommand - should show help
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	// Should not return error as it shows help
	assert.NoError(t, err)

	// Read the output
	out, err := io.ReadAll(b)
	assert.NoError(t, err)

	// Should contain help text
	output := string(out)
	assert.Contains(t, output, "generate command provides various sub-commands")
	assert.Contains(t, output, "Available Commands:")
	assert.Contains(t, output, "swagger")
}

func TestGenerateSwaggerExecute_BasicOutput(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	// Capture both stdout and stderr
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute the command - this will try to run real commands but we'll capture the initial output
	cmd.SetArgs([]string{})
	cmd.Execute() // Ignore error as external commands might fail in test environment

	// Check that the initial messages are printed correctly
	stdoutOutput, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	output := string(stdoutOutput)
	// Just check that we get the basic messages - don't worry about the actual swagger execution
	assert.Contains(t, output, "Generating Swagger/OpenAPI documentation")
	assert.Contains(t, output, "Output directory")
}

func TestGenerateSwaggerExecuteE_ReturnsNoError(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	// Test that the ExecuteE wrapper function doesn't return errors
	err := generateSwaggerExecuteE(cmd, []string{})
	assert.NoError(t, err, "generateSwaggerExecuteE should always return nil")
}

func TestSwagCommandsExistence(t *testing.T) {
	// Test that the swag command functions don't panic when called
	cmd := &cobra.Command{}

	// Test checkSwagInstallation doesn't panic
	assert.NotPanics(t, func() {
		checkSwagInstallation(cmd)
	}, "checkSwagInstallation should not panic")

	// We can't test the actual functionality without swag being installed
	// but we can ensure the functions are callable
}

func TestGenerateSwaggerExecute_OutputDirectory(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	// Capture stdout
	stdout := bytes.NewBufferString("")
	cmd.SetOut(stdout)

	// Execute the command (it will fail at some point but should print the output dir message)
	cmd.SetArgs([]string{})
	cmd.Execute() // Ignore error as external commands will fail in test environment

	// Check that output directory message is printed
	stdoutOutput, err := io.ReadAll(stdout)
	assert.NoError(t, err)

	output := string(stdoutOutput)
	assert.Contains(t, output, "./docs", "Should contain default output directory path")
}

func TestGenerateSwagger_WithInvalidArgs(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	// Test with invalid arguments - should still work as swagger command takes no args
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute with extra arguments
	cmd.SetArgs([]string{"invalid", "args"})
	cmd.Execute() // Should handle gracefully

	// Should still print the start message
	stdoutOutput, _ := io.ReadAll(stdout)
	output := string(stdoutOutput)
	assert.Contains(t, output, GenerateSwaggerStartMessage)
}

func TestGenerateDocumentation_ErrorReporting(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	// Capture both outputs
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute - will likely fail at swag installation/execution
	cmd.Execute()

	// Should have some output (either success messages or error messages)
	stdoutOutput, _ := io.ReadAll(stdout)
	stderrOutput, _ := io.ReadAll(stderr)

	totalOutput := string(stdoutOutput) + string(stderrOutput)

	// Should contain at least the start message
	assert.Contains(t, totalOutput, GenerateSwaggerStartMessage)

	// If there are errors, they should be properly formatted with "Gonyx >"
	if len(stderrOutput) > 0 {
		errorOutput := string(stderrOutput)
		lines := strings.Split(errorOutput, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				assert.Contains(t, line, "Gonyx >", "Error messages should be properly formatted")
			}
		}
	}
}

func TestIndividualSwagFunctions_ErrorHandling(t *testing.T) {
	cmd := &cobra.Command{}

	// Test checkSwagInstallation
	t.Run("checkSwagInstallation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			checkSwagInstallation(cmd)
		})
	})

	// Test installSwag
	t.Run("installSwag", func(t *testing.T) {
		assert.NotPanics(t, func() {
			installSwag(cmd)
		})
	})

	// Test runSwagInit
	t.Run("runSwagInit", func(t *testing.T) {
		assert.NotPanics(t, func() {
			runSwagInit(cmd)
		})
	})

	// Test runSwagFmt
	t.Run("runSwagFmt", func(t *testing.T) {
		assert.NotPanics(t, func() {
			runSwagFmt(cmd)
		})
	})
}

func TestGenerateSwagger_MessagesSequence(t *testing.T) {
	cmd := NewGenerateSwaggerCmd()

	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute the command
	cmd.Execute()

	stdoutOutput, _ := io.ReadAll(stdout)
	output := string(stdoutOutput)

	// Verify the sequence of messages appears in correct order
	expectedSequence := []string{
		GenerateSwaggerStartMessage,
		GenerateSwaggerOutputDirMessage,
		GenerateSwaggerCheckingMessage,
	}

	lastIndex := 0
	for i, expectedMsg := range expectedSequence {
		// Extract just the message part for searching (remove format specifiers)
		msgToSearch := strings.Split(expectedMsg, "%")[0]
		index := strings.Index(output[lastIndex:], msgToSearch)

		assert.True(t, index >= 0, "Message %d ('%s') should appear in output after previous messages", i, msgToSearch)
		if index >= 0 {
			lastIndex += index + len(msgToSearch)
		}
	}
}

func TestConstants_MessageFormatting(t *testing.T) {
	// Test specific formatting requirements
	testCases := []struct {
		name     string
		constant string
		checks   []string
	}{
		{
			name:     "Start Message",
			constant: GenerateSwaggerStartMessage,
			checks:   []string{"Gonyx >", "Generating", "Swagger"},
		},
		{
			name:     "Complete Message",
			constant: GenerateSwaggerCompleteMessage,
			checks:   []string{"Gonyx >", "successfully", "%s"},
		},
		{
			name:     "Error Message",
			constant: GenerateSwaggerErrorMessage,
			checks:   []string{"Gonyx >", "Error", "%v"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, check := range tc.checks {
				assert.Contains(t, tc.constant, check,
					"Constant %s should contain '%s'", tc.name, check)
			}
		})
	}
}

func TestGenerateCommand_SubcommandIntegration(t *testing.T) {
	generateCmd := NewGenerateCmd()

	// Verify swagger subcommand is properly registered
	swaggerCmd, _, err := generateCmd.Find([]string{"swagger"})
	assert.NoError(t, err, "Should find swagger subcommand")
	assert.Equal(t, "swagger", swaggerCmd.Name(), "Should return correct subcommand")

	// Test that invalid subcommand returns error (this is expected behavior)
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	generateCmd.SetOut(stdout)
	generateCmd.SetErr(stderr)
	generateCmd.SetArgs([]string{"invalid-subcommand"})

	err = generateCmd.Execute()
	// Invalid subcommand should return an error
	assert.Error(t, err, "Invalid subcommand should return error")
	assert.Contains(t, err.Error(), "unknown command", "Error should mention unknown command")
}

func TestConstants_AreProperlyDefined(t *testing.T) {
	// Test that all constants are defined and contain expected content
	constants := map[string]string{
		"GenerateSwaggerStartMessage":        GenerateSwaggerStartMessage,
		"GenerateSwaggerCompleteMessage":     GenerateSwaggerCompleteMessage,
		"GenerateSwaggerErrorMessage":        GenerateSwaggerErrorMessage,
		"GenerateSwaggerOutputDirMessage":    GenerateSwaggerOutputDirMessage,
		"GenerateSwaggerCheckingMessage":     GenerateSwaggerCheckingMessage,
		"GenerateSwaggerNotFoundMessage":     GenerateSwaggerNotFoundMessage,
		"GenerateSwaggerInstallSuccessMsg":   GenerateSwaggerInstallSuccessMsg,
		"GenerateSwaggerAlreadyInstalledMsg": GenerateSwaggerAlreadyInstalledMsg,
		"GenerateSwaggerInitRunningMsg":      GenerateSwaggerInitRunningMsg,
		"GenerateSwaggerInitCompleteMsg":     GenerateSwaggerInitCompleteMsg,
		"GenerateSwaggerFmtRunningMsg":       GenerateSwaggerFmtRunningMsg,
		"GenerateSwaggerFmtCompleteMsg":      GenerateSwaggerFmtCompleteMsg,
	}

	for name, constant := range constants {
		t.Run(name, func(t *testing.T) {
			assert.NotEmpty(t, constant, "%s should not be empty", name)
			assert.Contains(t, constant, "Gonyx >", "%s should contain 'Gonyx >' prefix", name)
		})
	}

	// Test error message constants
	errorConstants := map[string]string{
		"GenerateSwaggerInstallFailMsg":      GenerateSwaggerInstallFailMsg,
		"GenerateSwaggerInitFailMsg":         GenerateSwaggerInitFailMsg,
		"GenerateSwaggerFmtFailMsg":          GenerateSwaggerFmtFailMsg,
		"GenerateSwaggerInstallStartFailMsg": GenerateSwaggerInstallStartFailMsg,
		"GenerateSwaggerInstallCmdFailMsg":   GenerateSwaggerInstallCmdFailMsg,
		"GenerateSwaggerInitStartFailMsg":    GenerateSwaggerInitStartFailMsg,
		"GenerateSwaggerInitCmdFailMsg":      GenerateSwaggerInitCmdFailMsg,
		"GenerateSwaggerFmtStartFailMsg":     GenerateSwaggerFmtStartFailMsg,
		"GenerateSwaggerFmtCmdFailMsg":       GenerateSwaggerFmtCmdFailMsg,
	}

	for name, constant := range errorConstants {
		t.Run(name+"_Error", func(t *testing.T) {
			assert.NotEmpty(t, constant, "%s should not be empty", name)
			assert.Contains(t, constant, "Gonyx >", "%s should contain 'Gonyx >' prefix", name)
			assert.True(t, strings.Contains(constant, "Failed") || strings.Contains(constant, "failed"),
				"%s should contain 'Failed' or 'failed' as it's an error message", name)
		})
	}
}

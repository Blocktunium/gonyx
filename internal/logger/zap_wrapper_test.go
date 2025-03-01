package logger

import (
	"bytes"
	"encoding/json"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"os"
	"testing"
	"time"
)

func Test_ZapFileLogger(t *testing.T) {
	path := "../.."
	initialMode := "test"
	prefix := "Gonyx"

	err := config.CreateManager(path, initialMode, prefix)
	if err != nil {
		t.Errorf("Initializig Config without error, but got %v", err)
		return
	}

	logg := &ZapWrapper{}
	logg.Constructor("logger")

	expectedModule := "tester"
	expectedService := "Gonyx"
	expectedLogType := "FUNC_MAINT"
	expectedM := "TEST"
	expectedL := "\\u001b[35mDEBUG\\u001b[0m"

	logg.Log(types.NewLogObject(
		types.DEBUG, "tester", types.FuncMaintenanceType, time.Now().UTC(), "TEST", nil))
	logg.Sync()

	// read a log file to ensure it writes ok --> for now it must be written to /tmp folder
	file, err := os.OpenFile("/tmp/Gonyx.log", os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Opening the logs file in /tmp folder, but got %v", err)
		return
	}
	defer file.Close()

	if file.Name() != "/tmp/Gonyx.log" {
		t.Errorf("Expected the name of the log file be %v, but got %v", "/tmp/Gonyx.log", file.Name())
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contents := buf.String()

	// encode str to json
	var result map[string]any
	err = json.Unmarshal([]byte(contents), &result)
	if err != nil {
		// json deserialize error
		t.Errorf("Expected json content able to be deserialized, but got %v", err)
		return
	}

	if !(result["L"] != expectedL && result["M"] == expectedM && result["module"] == expectedModule &&
		result["service"] == expectedService && result["log_type"] == expectedLogType && result["additional"] == nil) {

		t.Errorf("Expected Value Is Not Correct: ")
		t.Errorf("Expected %v, but got %v --", result["L"], expectedL)
		t.Errorf("Expected %v, but got %v --", result["M"], expectedM)
		t.Errorf("Expected %v, but got %v --", result["module"], expectedModule)
		t.Errorf("Expected %v, but got %v --", result["service"], expectedService)
		t.Errorf("Expected %v, but got %v --", result["log_type"], expectedLogType)
		t.Errorf("Expected %v, but got %v --", result["additional"], nil)
		return
	}

}

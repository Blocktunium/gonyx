package logger

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_ZhycanConsoleLogger(t *testing.T) {
	path := "../.."
	initialMode := "test"
	prefix := "ZHYCAN"

	err := config.CreateManager(path, initialMode, prefix)
	if err != nil {
		t.Errorf("Initializig Config without error, but got %v", err)
		return
	}

	//old := os.Stdout // keep backup of the real stdout
	//r, w, _ := os.Pipe()
	//os.Stdout = w

	done := capture()

	logg := &LogMeWrapper{}
	err = logg.Constructor("logger")
	if err != nil {
		t.Errorf("Initializing the Zhycan Wrapper, Expected to don't have error, but got %v", err)
	}

	expectedFlag := true
	if !logg.IsInitialized() {
		t.Errorf("Zhycan Wrapper Must Be Initilaized, Expected to get: %v, but got %v", expectedFlag, logg.IsInitialized())
	}

	logTime := time.Now().UTC()
	logg.Log(types.NewLogObject(
		types.DEBUG, "tester", types.FuncMaintenanceType, logTime, "TEST", nil))
	logg.Sync()

	time.Sleep(time.Second * 5)
	out, err := done()
	//outC := make(chan string)
	//// copy the output in a separate goroutine so printing can't block indefinitely
	//go func() {
	//	var buf bytes.Buffer
	//	io.Copy(&buf, r)
	//	outC <- buf.String()
	//}()

	// back to normal state
	//w.Close()
	//os.Stdout = old // restoring the real stdout
	//out := <-outC

	// Check the output
	expectedLog := fmt.Sprintf("\\e[37mzhycan %v >>>   DEBUG >>> (FUNC_MAINT/tester)  - tester ... <nil>\\e[0m\n", logTime.UnixNano())
	if !reflect.DeepEqual(out, expectedLog) {
		t.Errorf("Expected Log must be: %v, but got: %v", expectedLog, out)
	}
}

// capture replaces os.Stdout with a writer that buffers any data written
// to os.Stdout. Call the returned function to cleanup and get the data
// as a string.
func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}

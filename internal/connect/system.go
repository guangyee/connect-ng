package connect

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

func execute(cmd []string, quiet bool, validExitCodes []int) ([]byte, error) {
	Debug.Printf("Executing: %s Quiet: %v\n", cmd, quiet)
	var stderr, stdout bytes.Buffer
	comm := execCommand(cmd[0], cmd[1:]...)
	comm.Stdout = &stdout
	comm.Stderr = &stderr
	// init env only if not set by (mocked) execCommand()
	if len(comm.Env) == 0 {
		comm.Env = os.Environ()
	}
	comm.Env = append(comm.Env, "LC_ALL=C")
	err := comm.Run()
	exitCode := comm.ProcessState.ExitCode()
	Debug.Printf("Return code: %d\n", exitCode)
	if stdout.Len() > 0 {
		Debug.Println("Output:", stdout.String())
	}
	if stderr.Len() > 0 {
		Debug.Println("Error:", stderr.String())
	}
	// TODO Ruby version also checks stderr for "ABORT request"
	if err != nil && !containsInt(validExitCodes, exitCode) {
		output := stderr.Bytes()
		// zypper with formatter option writes to stdout instead of stderr
		if len(output) == 0 {
			output = stdout.Bytes()
		}
		output = bytes.TrimSuffix(output, []byte("\n"))
		ee := ExecuteError{Commmand: cmd, ExitCode: exitCode, Output: output, Err: err}
		return nil, ee
	}
	if quiet {
		return nil, nil
	}
	out := stdout.Bytes()
	out = bytes.TrimSuffix(out, []byte("\n"))
	return out, nil
}

func containsInt(s []int, i int) bool {
	for _, e := range s {
		if e == i {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func removeFile(path string) error {
	Debug.Print("Removing file: ", path)
	if !fileExists(path) {
		return nil
	}
	return os.Remove(path)
}

func isRootFSWritable() bool {
	_, err := execute([]string{"test", "-w", "/"}, true, []int{zypperOK})
	return err == nil
}

// Cleanup removes system credentials and installed services
func Cleanup() error {
	err := removeSystemCredentials()
	if err != nil {
		return err
	}

	// remove all suse services from zypper
	services, err := installedServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		// NOTE: this check might not work correctly with SMT depending
		//       on the configuration (e.g. listen on https but API
		//       returns URL with http).
		if !strings.Contains(service.URL, CFG.BaseURL) {
			fmt.Printf("%s not in %s\n", CFG.BaseURL, service.URL)
			continue
		}
		if err := removeService(service.Name); err != nil {
			return err
		}

	}
	return nil
}

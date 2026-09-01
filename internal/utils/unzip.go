package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// UnzipCommandAvailable Check if unzip command available on this system
// if unavilable, program will exit with code 1
// Call this function before you start router
func UnzipCommandAvailable() {
	cmd := exec.Command("command", "-v", "unzip")

	if err := cmd.Run(); err != nil {
		// Exit the program because the unzip command is core of the program for unpack new post
		TextLogger.Error("error when check the unzip command, you should provide now, exiting system", "err", err)
		os.Exit(1)
	}
}

// Unzip exec the command unzip srcFilePath -d dstDirPath
// it will print output into Stdout
func Unzip(srcFilePath string, dstDirPath string) error {
	cmd := exec.Command("unzip", srcFilePath, "-d", dstDirPath)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	TextLogger.Info("begin unzip package", "exec_cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		// TODO: Clear the dstDirPath ?

		return fmt.Errorf("error when unzip package, exec cmd: %s, err: %w", cmd.String(), err)
	}

	TextLogger.Info("complete the unzip package, success", "src", srcFilePath, "dst", dstDirPath)

	return nil
}

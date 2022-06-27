package cmd

import (
	"fmt"
	"log"
	"os/exec"
)

type CmdFactory struct {
	dir string
}

func NewCmdFactory(dir string) *CmdFactory {
	return &CmdFactory{
		dir: dir,
	}
}

func (f *CmdFactory) ExecF(command string, params ...interface{}) (string, error) {
	cmdText := fmt.Sprintf(command, params...)
	cmd := exec.Command("bash", "-c", cmdText)
	cmd.Dir = f.dir

	out, err := cmd.CombinedOutput() //may need to separate stdout and stderr
	if err != nil {
		log.Printf("command: %s", cmd.String())
		log.Printf("output: %s", out)
		log.Printf("error: %s", err)
		return string(out), err
	}

	return string(out), nil
}

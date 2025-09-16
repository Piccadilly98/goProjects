package vpn

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
)

const (
	programmName  = "openvpn"
	flagConfig    = "--config"
	flagChangeDir = "--cd"
)

type buffer struct {
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func StartVPN(config *data_structs.VPNConfig, logs *data_structs.InitInfo) error {
	pwd, err := os.Getwd()
	if pwd != "" {

	}
	if err != nil {
		return fmt.Errorf("errrors in start vpn\nProblem in get word directory")
	}
	err = os.Chdir(filepath.Dir(logs.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("error change directory")
	}
	// buf := buffer{}
	pwd1, _ := os.Getwd()
	fmt.Println(filepath.Join(pwd1, logs.Filename))
	command := exec.Command(programmName, flagConfig, filepath.Join(pwd1, logs.Filename))
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Start()

	// fmt.Printf("%s\n%s\n%s", buf.stdin.String(), buf.stdout.String(), buf.stderr.String())
	return nil
}

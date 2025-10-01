package vpn

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
)

const (
	programmName  = "openvpn"
	flagConfig    = "--config"
	flagChangeDir = "--cd"
)

// func StartVPN(config *data_structs.VPNConfig) error {
// 	pwd, err := os.Getwd()
// 	defer os.Chdir(pwd)
// 	if err != nil {
// 		return fmt.Errorf("errrors in start vpn\nProblem in get work directory")
// 	}
// 	err = os.Chdir(filepath.Dir(config.Logs.ConfigFilePath))
// 	if err != nil {
// 		return fmt.Errorf("error change directory")
// 	}
// 	newRepository, _ := os.Getwd()
// 	// wg := sync.WaitGroup{}
// 	command := exec.Command(programmName, flagConfig, filepath.Join(newRepository, config.Logs.Filename))
// 	stdIn, _ := command.StdinPipe()
// 	stdOut, _ := command.StdoutPipe()
// 	stdErr, _ := command.StderrPipe()
// 	err = command.Start()
// 	if err != nil {
// 		return fmt.Errorf("cannot start command: %v", err)
// 	}
// 	go func() {
// 		readBuf(stdIn, stdOut, stdErr, config)
// 	}()
// 	go readStdout(stdOut)
// 	// time.Sleep(100 * time.Millisecond)
// 	// p, err := command.Stdin.Read([]byte(config.Logs.Name))
// 	// if p == 0 || err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Println(p, err)
// 	// // fmt.Printf("%s\n%s\n%s", buf.stdin.String(), buf.stdout.String(), buf.stderr.String())
// 	command.Wait()
// 	return nil
// }

// func readBuf(stdIn io.WriteCloser, stdOut io.ReadCloser, stdErr io.ReadCloser, config *data_structs.VPNConfig) {
// 	scanner := bufio.NewScanner(stdErr)
// 	// time.Sleep(100 * time.Millisecond)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		// fmt.Println(line)
// 		if strings.Contains(line, "Enter Auth Username:") {
// 			io.WriteString(stdIn, config.Logs.Name+"\n")
// 		} else if strings.Contains(line, "Enter Auth Password:") {
// 			io.WriteString(stdIn, config.Logs.Password+"\n")
// 			stdIn.Close()
// 		}
// 		// time.Sleep(100 * time.Millisecond)
// 	}
// }

// func readStdout(stdOut io.ReadCloser) {
// 	scanner := bufio.NewScanner(stdOut)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		fmt.Printf("LOG: %s\n", line)
// 	}
// }

func StartVPN(config *data_structs.VPNConfig) error {
	pwd, err := os.Getwd()
	defer os.Chdir(pwd)
	if err != nil {
		return fmt.Errorf("errors in start vpn: Problem in get work directory")
	}
	err = os.Chdir(filepath.Dir(config.Logs.ConfigFilePath))
	if err != nil {
		return fmt.Errorf("error change directory")
	}
	newRepository, _ := os.Getwd()

	command := exec.Command(programmName, flagConfig, filepath.Join(newRepository, config.Logs.Filename))
	stdIn, _ := command.StdinPipe()
	stdOut, _ := command.StdoutPipe()
	stdErr, _ := command.StderrPipe()

	err = command.Start()
	if err != nil {
		return fmt.Errorf("cannot start command: %v", err)
	}

	// Читаем из ОБОИХ потоков
	go readStderr(stdIn, stdErr, config)
	go readStdout(stdOut) // Важно: читаем stdout чтобы не блокировать процесс

	command.Wait()
	return nil
}

func readStderr(stdIn io.WriteCloser, stdErr io.ReadCloser, config *data_structs.VPNConfig) {
	defer stdIn.Close()

	scanner := bufio.NewScanner(stdErr)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("PROMPT: %s\n", line) // Отладочный вывод

		if strings.Contains(line, "Enter Auth Username:") {
			fmt.Printf("SENDING USERNAME: %s\n", config.Logs.Name)
			io.WriteString(stdIn, config.Logs.Name+"\n")
		} else if strings.Contains(line, "Enter Auth Password:") {
			fmt.Printf("SENDING PASSWORD: %s\n", config.Logs.Password)
			io.WriteString(stdIn, config.Logs.Password+"\n")
			// Не закрываем stdin сразу - может понадобиться для повторной аутентификации
		}
	}
}

func readStdout(stdOut io.ReadCloser) {
	scanner := bufio.NewScanner(stdOut)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("LOG: %s\n", line)
	}
}

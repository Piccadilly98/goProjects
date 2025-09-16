package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
)

func AuthenticationManager(config *data_structs.VPNConfig, logs *data_structs.InitInfo) error {
	var err error
	if logs.Password != "empty" && config.AuthUserPassFilename != "" {
		config.AuthUserPassFilename = ""
	}
	if config.AuthUserPass {
		if config.AuthUserPassFilename != "" {
			logs.Name, logs.Password, err = readFile(config.AuthUserPassFilename)
			if err != nil {
				return err
			}
		} else if logs.Name == "user" && logs.Password == "empty" {
			input := inputData()
			if len(input) == 2 {
				logs.Name, logs.Password = input[0], input[1]
			}
		}
	}
	return nil
}

func inputData() []string {
	result := make([]string, 2)
	fmt.Println("\nPlease, input your name for config file: ")
	fmt.Scan(&result[0])
	fmt.Println("\nPlease, input your password: ")
	fmt.Scan(&result[1])
	return result
}

func readFile(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	result := make([]string, 0)
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lineResult := strings.Split(strings.Split(line, "#")[0], " ")
		if len(lineResult) != 1 {
			return "", "", fmt.Errorf("file %s for auntefication incorrect", filename)
		}
		result = append(result, lineResult[0])
	}
	if len(result) >= 2 {
		return result[0], result[1], nil
	} else {
		return "", "", fmt.Errorf("incorrect file %s", filename)
	}
}

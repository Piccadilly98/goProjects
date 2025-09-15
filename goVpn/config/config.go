package parse_config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
	"github.com/Piccadilly98/goProjects/goVpn/validation"
)

func ParseConfig(filename string) error {
	if !strings.Contains(filename, ".ovpn") {
		return errors.New("filename is empty or not .ovpn")
	}
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("no such file")
	}
	defer file.Close()
	dir := filepath.Dir(filename)
	scan := bufio.NewScanner(file)
	config := data_structs.VPNConfig{}
	i := 0
	for scan.Scan() {

		if strings.HasPrefix(scan.Text(), "#") {
			continue
		}
		line := strings.TrimSpace(scan.Text())
		lineSplitComment := strings.Split(line, "#")
		lineSplitSpace := strings.Split(lineSplitComment[0], " ")
		sliceResult := make([]string, 0)
		for _, word := range lineSplitSpace {
			if word != "" && word != " " {
				sliceResult = append(sliceResult, strings.TrimSpace(word))
			}
		}
		if len(sliceResult) == 0 {
			continue
		}

		if strings.HasPrefix(sliceResult[0], "<") && strings.HasSuffix(sliceResult[0], ">") {
			if len(sliceResult) == 1 {
				tag := strings.Trim(sliceResult[0], "<>")
				err = ReadBlock(tag, scan, &config)
				if err != nil {
					return err
				}
			}
		}
		if i == 0 && sliceResult[0] != "client" {
			return errors.New("incorrect config, is not contains \"client\"")
		}
		tag := strings.TrimSpace(sliceResult[0])
		switch tag {
		case "remote":
			if len(sliceResult) == 3 {
				config.RemoteHost = sliceResult[1]
				config.RemotePort, err = strconv.Atoi(sliceResult[2])
				if err != nil {
					return errors.New("invalid port number")
				}
			}
		case "ca":
			if len(sliceResult) == 2 {

				path := filepath.Join(dir, sliceResult[1])
				config.CaFilename = path

			}
		case "cert":
			if len(sliceResult) == 2 {
				path := filepath.Join(dir, sliceResult[1])
				config.CertFilename = path
			}

		case "key":
			if len(sliceResult) == 2 {
				path := filepath.Join(dir, sliceResult[1])
				config.KeyFileName = path
			}
		case "secret":
			if len(lineSplitSpace) == 2 {
				path := filepath.Join(dir, sliceResult[1])
				config.SecretFilename = path
			}

		case "proto":
			if len(sliceResult) == 2 {
				if config.Proto == "" {
					config.Proto = sliceResult[1]
				}
			}

		case "auth-user-pass":
			config.AuthUserPass = true
			if len(sliceResult) == 2 {
				path := filepath.Join(dir, sliceResult[1])
				config.AuthUserPassFilename = path
			}
		}
		i++
	}
	// fmt.Println(config)
	_, err = validation.ValidateConfigInfo(&config)
	if err != nil {
		return err
	}
	return nil
}

func ReadBlock(tag string, scanner *bufio.Scanner, config *data_structs.VPNConfig) error {
	var content string
	endBlock := "</" + tag + ">"
	isContainsEnd := false
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == endBlock {
			isContainsEnd = true
			break
		}
		content += scanner.Text()

	}
	if !isContainsEnd {
		return fmt.Errorf("tag %s is not contains end tag(%s)", tag, endBlock)
	}
	content = strings.TrimSpace(content)
	switch tag {
	case "ca":
		if config.CaFilename == "" {
			config.CaInbuilt = content
		} else {
			return errors.New("invalid, block ca repeat")
		}
	case "key":
		if config.KeyFileName == "" {
			config.KeyInbuilt = content
		} else {
			return errors.New("invalid, block key repeat")
		}
	case "cert":
		if config.CertFilename == "" {
			config.CertInbuilt = content
		} else {
			return errors.New("invalid, block cert repeat")
		}
	case "tls-auth":
		if config.TlsAuth == "" {
			config.TlsAuth = content
		} else {
			return errors.New("invalid, block secret repeat")
		}
	}
	return nil
}

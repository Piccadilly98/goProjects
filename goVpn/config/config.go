package main

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
	for scan.Scan() {

		if strings.HasPrefix(scan.Text(), "#") {
			continue
		}
		line := strings.TrimSpace(scan.Text())
		lineSplitComment := strings.Split(line, "#")
		lineSplitSpace := strings.Split(lineSplitComment[0], " ")

		if strings.HasPrefix(lineSplitSpace[0], "<") && strings.HasSuffix(lineSplitSpace[0], ">") {
			if len(lineSplitSpace) == 1 {
				tag := strings.Trim(lineSplitSpace[0], "<>")
				err = ReadBlock(tag, scan, &config)
				if err != nil {
					return err
				}
			}
		}

		switch lineSplitSpace[0] {
		case "remote":
			if len(lineSplitSpace) == 3 {
				if config.RemoteHost == "" && config.RemotePort == 0 {
					config.RemoteHost = lineSplitSpace[1]
					config.RemotePort, err = strconv.Atoi(lineSplitSpace[2])
					if err != nil {
						return errors.New("invalid port number")
					}
				} else {
					return errors.New("invalid repeat port num or len remote < 3")
				}
			}
		case "ca":
			if len(lineSplitSpace) == 2 {
				if config.CaFilename == "" && config.CaInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.CaFilename = path
				} else {
					return errors.New("invalid repeat ca")
				}
			}
		case "cert":
			if len(lineSplitSpace) == 2 {
				if config.CertFilename == "" && config.CertInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.CertFilename = path
				} else {
					return errors.New("invalid repeat cert")
				}
			}

		case "key":
			if len(lineSplitSpace) == 2 {
				if config.KeyFileName == "" && config.KeyInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.KeyFileName = path
				} else {
					return errors.New("invalid repeat key")
				}
			}
		case "secret":
			if len(lineSplitSpace) == 2 {
				if config.SecretInbuilt == "" && config.SecretFilename == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.SecretFilename = path
				} else {
					return errors.New("invalid repeat secret")
				}
			}

		case "proto":
			if len(lineSplitSpace) == 2 {
				if config.Proto == "" {
					config.Proto = lineSplitSpace[1]
				} else {
					return errors.New("invalid repeat proto")
				}
			}

		case "auth-user-pass":
			config.AuthUserPass = true

		}
	}
	fmt.Println(config)
	_, err = validation.ValidateConfigInfo(&config)
	if err != nil {
		return err
	}
	return nil
}

func ReadBlock(tag string, scanner *bufio.Scanner, config *data_structs.VPNConfig) error {
	var content string
	if tag != "ca" && tag != "key" && tag != "cert" && tag != "Secret" {
		return errors.New("incorrect tag")
	}
	endBlock := "</" + tag + ">"
	for scanner.Scan() {
		if scanner.Text() == endBlock {
			break
		}
		content += scanner.Text()

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
	case "secret":
		if config.SecretInbuilt == "" {
			config.SecretInbuilt = content
		} else {
			return errors.New("invalid, block secret repeat")
		}
	}
	return nil
}

func main() {
	err := ParseConfig("../text.ovpn")
	fmt.Println(err)
}

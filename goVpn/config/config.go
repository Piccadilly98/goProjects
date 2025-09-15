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
	i := 0
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
		if i == 0 && lineSplitSpace[0] != "client" {
			return errors.New("incorrect config, is not contains \"client\"")
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
				}
			}
		case "ca":
			if len(lineSplitSpace) == 2 {
				if config.CaFilename == "" && config.CaInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.CaFilename = path
				}
			}
		case "cert":
			if len(lineSplitSpace) == 2 {
				if config.CertFilename == "" && config.CertInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.CertFilename = path
				}
			}

		case "key":
			if len(lineSplitSpace) == 2 {
				if config.KeyFileName == "" && config.KeyInbuilt == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.KeyFileName = path
				}
			}
		case "secret":
			if len(lineSplitSpace) == 2 {
				if config.TlsAuth == "" && config.SecretFilename == "" {
					path := filepath.Join(dir, lineSplitSpace[1])
					config.SecretFilename = path
				}
			}

		case "proto":
			if len(lineSplitSpace) == 2 {
				if config.Proto == "" {
					config.Proto = lineSplitSpace[1]
				}
			}

		case "auth-user-pass":
			config.AuthUserPass = true

		}
		i++
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
	// if tag != "ca" && tag != "key" && tag != "cert" && tag != "Secret" {
	// 	return errors.New("incorrect tag")
	// }
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

func main() {
	err := ParseConfig("../text.ovpn")
	fmt.Println(err)
}

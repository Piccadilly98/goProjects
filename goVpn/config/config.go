package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type VPNConfig struct {
	RemoteHost     string
	RemotePort     int
	ca_filename    string
	ca_inbuilt     string
	cert_filename  string
	cert_inbuilt   string
	key_filename   string
	key_inbuilt    string
	secret         string
	auth_user_pass bool
	proto          string
}

func (v VPNConfig) String() string {
	return fmt.Sprintf("Host: %s\nPort: %d\nca_file: %s\nca_inbuilt: %v\ncert_file: %s\ncert_inbuilt: %v\nkey_file: %s\nkey_inbuilt: %v\nsecret: %v\nauth: %v\nproto: %s", v.RemoteHost, v.RemotePort, v.ca_filename, v.ca_inbuilt, v.cert_filename, v.cert_inbuilt, v.key_filename, v.key_inbuilt, v.secret, v.auth_user_pass, v.proto)
}

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
	config := VPNConfig{}
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
						return errors.New("invalid")
					}
				} else {
					return errors.New("invalid ")
				}
			}
		case "ca":
			if len(lineSplitSpace) == 2 {
				if config.ca_filename == "" {
					config.ca_filename = filepath.Join(dir, lineSplitSpace[1])
				} else {
					return errors.New("invalid ")
				}
			}
		case "cert":
			if len(lineSplitSpace) == 2 {
				if config.cert_filename == "" {
					config.cert_filename = filepath.Join(dir, lineSplitSpace[1])
				} else {
					return errors.New("invalid ")
				}
			}

		case "key":
			if len(lineSplitSpace) == 2 {
				if config.key_filename == "" {
					config.key_filename = filepath.Join(dir, lineSplitSpace[1])
				} else {
					return errors.New("invalid ")
				}
			}
		case "secret":
			if len(lineSplitSpace) == 2 {
				if config.secret == "" {
					config.secret = filepath.Join(dir, lineSplitSpace[1])
				} else {
					return errors.New("invalid ")
				}
			}

		case "proto":
			if len(lineSplitSpace) == 2 {
				if config.proto == "" {
					config.proto = lineSplitSpace[1]
				} else {
					return errors.New("invalid ")
				}
			}

		case "auth-user-pass":
			config.auth_user_pass = true

		}
	}
	fmt.Println(config)
	return nil
}

func ReadBlock(tag string, scanner *bufio.Scanner, config *VPNConfig) error {
	var content string
	if tag != "ca" && tag != "key" && tag != "cert" && tag != "secret" {
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
		config.ca_inbuilt = content
	case "key":
		config.key_inbuilt = content
	case "cert":
		config.cert_inbuilt = content
	case "secret":
		config.secret = content
	}
	return nil
}

func main() {
	err := ParseConfig("../text.ovpn")
	if err != nil {
		fmt.Println(err)
	}
}

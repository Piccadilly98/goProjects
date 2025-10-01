package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Piccadilly98/goProjects/goVpn/auth"
	parsing_config "github.com/Piccadilly98/goProjects/goVpn/config"
	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
	"github.com/Piccadilly98/goProjects/goVpn/validation"
	"github.com/Piccadilly98/goProjects/goVpn/vpn"
)

func main() {
	logs := data_structs.NewInitInfo()
	configFile := flag.String("config", "no such file", "During startup, use -config <filepath> to enter the path to the configuration file (.ovpn). \nWithout this flag, the program will not start.")
	username := flag.String("username", "user", "If you want to write your nickname from the config for connect to vpn, use -username <nickname>")
	password := flag.String("password", "empty", "If you want to record your configuration password and/or have already recorded your login, use -password <password>. The program will not start if you have already specified a username and have not specified a password.If you do not use -password and -username, the program will check and get password and username later.")
	flag.Parse()

	logs.Name = *username
	logs.Password = *password
	logs.ConfigFilePath = *configFile
	err := processingFlags(logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Start Parsing config: %v ...\n", logs.ConfigFilePath)
	config, err := parsing_config.ParseConfig(logs)
	if err != nil {
		log.Fatalf("Error parsing config:\n%v", err)
	}
	fmt.Println("The configuration is checked, let's proceed to testing...")
	err = validation.ValidateConfigInfo(config)
	if err != nil {
		log.Fatalf("Error validation:\n%v", err)
	}
	fmt.Printf("Configuration finish parsing and validation\n")
	fmt.Println("Checking your input data and proccess your authentication..")
	err = auth.AuthenticationManager(config)
	if err != nil {
		log.Fatalf("Error Auntification:\n%v", err)
	}
	fmt.Println("Configuration preparation and processing complete.\nStarting VPN...")
	// fmt.Println(time.Since(logs.TimeInit))
	// fmt.Println(config)
	err = vpn.StartVPN(config)
	if err != nil {
		log.Fatal(err)
	}

	// err = vpn.StartVPN(config, logs)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func processingFlags(log *data_structs.InitInfo) error {
	if log.ConfigFilePath == "no such file" {
		return fmt.Errorf("incorrect input, please repeat start programm and write configfilename")
	}
	if (log.Name != "user" && log.Password == "empty") || (log.Password != "empty" && log.Name == "user") {
		return fmt.Errorf("incorrect input auth data.\nPasword without login or login without password")
	}
	return nil
}

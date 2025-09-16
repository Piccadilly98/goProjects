package validation

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/Piccadilly98/goProjects/goVpn/data_structs"
)

func ValidateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	} else if info.Size() == 0 {
		return errors.New("file is empty")
	}
	return nil
}

func ValidateConfigInfo(config *data_structs.VPNConfig, logs *data_structs.InitInfo) error {
	err := ValidationHost(config.RemoteHost)
	if err != nil {
		return err
	}
	err = validationPort(config.RemotePort)
	if err != nil {
		return err
	}
	if config.CaFilename != "" {
		err := ValidateFile(config.CaFilename)
		if err != nil {
			return err
		}
	}

	if (config.CaFilename == "" && config.CaInbuilt == "") && (config.TlsAuth == "" && config.SecretFilename == "") {
		return errors.New("invalid!\nNot contains Ca or Secret ")
	}

	//приоритет
	if config.SecretFilename == "" && config.TlsAuth == "" {
		if config.CertFilename != "" {
			err := ValidateFile(config.CertFilename)
			if err != nil {
				return err
			}
		}

		if config.KeyFileName != "" {
			err := ValidateFile(config.KeyFileName)
			if err != nil {
				return err
			}
		}
		if (config.CertFilename != "" || config.CertInbuilt != "") && (config.KeyFileName == "" && config.KeyInbuilt == "") {
			return errors.New("config contain cert and not contain key")
		} else if (config.KeyFileName != "" || config.KeyInbuilt != "") && (config.CertInbuilt == "" && config.CertFilename == "") {
			return errors.New("config contain key and not contain cert")
		}
	}
	if (config.CertFilename != "" || config.CertInbuilt != "") && (config.TlsAuth != "" || config.SecretFilename != "") {
		log.Println("config contain cert and secret\nSecret is priority")
	}
	if config.SecretFilename != "" {
		if config.TlsAuth != "" {
			return errors.New("invalid, block secret repeat")
		}
		err := ValidateFile(config.SecretFilename)
		if err != nil {
			return err
		}
	}
	if config.AuthUserPass && config.AuthUserPassFilename != "" {
		err := ValidateFile(config.AuthUserPassFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidationHost(remoteHost string) error {
	if remoteHost == "" {
		return errors.New("empty adress")
	} else if strings.HasPrefix(remoteHost, " ") || strings.HasSuffix(remoteHost, " ") {
		return errors.New("adress contains spacebar")
	}
	if ip := net.ParseIP(remoteHost); ip != nil {
		return nil
	}
	if looksLikeIPButInvalid(remoteHost) {
		return errors.New("incorrect ip")
	}

	err := validateDomainName(remoteHost)
	if err != nil {
		return err
	}
	return nil
}

func validationPort(remotePort int) error {
	if remotePort >= 1 && remotePort <= 65535 {
		return nil
	}
	return errors.New("invalid port number")
}

func validateDomainName(remoteHost string) error {
	if utf8.RuneCountInString(remoteHost) > 253 {
		return errors.New("invalid domain name, lenght > 253")
	}
	if strings.HasPrefix(remoteHost, "-") || strings.HasSuffix(remoteHost, "-") {
		return errors.New("'-' not contains in begin hostname")
	}
	if remoteHost == "localhost" {
		return nil
	}
	for i, symbol := range remoteHost {
		if !((symbol >= 'a' && symbol <= 'z') || (symbol >= 'A' && symbol <= 'Z') || (symbol >= '0' && symbol <= '9') || (symbol <= '.' && symbol >= '-')) {
			return errors.New("invalid symbols in domain name")
		}
		if (i == 0 || i == utf8.RuneCountInString(remoteHost)-1) && (symbol == '.' || symbol == '-') {
			return errors.New(". or - can't be at the beginning or at the end")
		}
	}
	splitDomain := strings.Split(remoteHost, ".")
	for _, segment := range splitDomain {
		if utf8.RuneCountInString(segment) > 63 {
			return errors.New("between '.' >63 symbols")
		} else if segment == "" {
			return errors.New("between '.' empty")
		}
		if segment[len(segment)-1] == '-' {
			return errors.New("ends with hyphen")
		}
	}
	return nil
}

func looksLikeIPButInvalid(host string) bool {
	if strings.Contains(host, ".") || strings.Contains(host, ":") {
		for _, symbol := range host {
			if !isIPCharacter(symbol) {
				return false
			}
		}
		return true
	}
	return false
}

func isIPCharacter(char rune) bool {
	return (char >= '0' && char <= '9') ||
		char == '.' ||
		char == ':' ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const configFile = ".keyrun.yml"

const helpMessage = `Commands:
  exec -- <command> <args>			- executes command with arguments
  encrypt <unencrypted source file>		- encrypts file read from source, writes encrypted into destination appending .enc extension
  decrypt <encrypted source file>		- decrypts file read from source, writes unencrypted into destination without .enc extension
  show <encrypted source file>			- decrypts file read from source and prints to stdout
  key <subcommand>:
	create					- create new key in keyring
	show					- show key stored in keyring
	delete					- delete key from keyring
  version					- show version
  help						- show help message
`

type ConfigType struct {
	Env map[string]string `yaml:"env"`
	Key string            `yaml:"key"`
}

// Parse unmarchals json into ConfigType struct
func (c *ConfigType) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// GetConf reads and parses config file
func GetConf(configFile string) (config ConfigType, err error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	config.Parse(data)
	return
}

func showHelpMessage() {
	fmt.Println(helpMessage)
	os.Exit(1)
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		showHelpMessage()
	}

	command := args[0]

	switch command {
	case "exec":
		if len(args) >= 3 {
			if args[1] == "--" {
				var secret string
				var files map[string]string

				conf, err := GetConf(configFile)
				if err != nil {
					log.Fatal("Error: ", err)
					os.Exit(1)
				}

				err = setEnv(conf.Env)
				if err != nil {
					log.Fatal("Error: ", err)
				}

				// retriving secret from keyring
				if conf.Key != "" {
					secret, err = getKey(conf.Key)
					if err != nil {
						log.Fatalln("Error retrieveing key", err)
						os.Exit(2)
					}
				}

				// unencrypting .enc files, recording md5 of unencrypted content
				if conf.Key != "" {
					files := findFiles()
					for fname := range files {
						filehash, err := decryptFile(fname+".enc", secret)
						if err != nil {
							log.Fatalln("Error decrypting file", fname)
							os.Exit(2)
						}
						files[fname] = filehash
					}
				}

				// executing command
				runit(args[2], args[3:])

				// encrypting state files
				if conf.Key != "" {
					for fname, filehash := range files {
						filehashNew, err := getFileHash(fname)
						if err != nil {
							log.Fatalln("Failed to read unencrypted", fname)
							os.Exit(2)
						}

						// skip encrypting file if md5 is the same
						if filehash == filehashNew {
							continue
						}

						err = encryptFile(fname, secret)
						if err != nil {
							log.Fatalln("Error encrypting file", fname)
							os.Exit(2)
						}
					}
				}

				// cleanup
				for fname := range files {
					err := os.Remove(fname)
					if err != nil {
						log.Fatalln("Error cleaning up,", err)
					}
				}
			}
		} else {
			showHelpMessage()
		}
	case "encrypt":
		if len(args) == 2 {
			conf, err := GetConf(configFile)
			if err != nil {
				log.Fatal("Error: ", err)
				os.Exit(1)
			}
			// retriving secret from keyring
			secret, err := getKey(conf.Key)
			if err != nil {
				log.Fatalln("Error: ", err)
				os.Exit(2)
			}
			// encrypting files
			err = encryptFile(args[1], secret)
			if err != nil {
				log.Fatalln("Error: ", args[1])
				os.Exit(2)
			}
			// cleanup
			os.Remove(args[1])
			fmt.Println("Encrypted ", args[1])
		} else {
			showHelpMessage()
		}
	case "decrypt":
		if len(args) == 2 {
			conf, err := GetConf(configFile)
			if err != nil {
				log.Fatal("Error: ", err)
				os.Exit(1)
			}
			// retriving secret from keyring
			secret, err := getKey(conf.Key)
			if err != nil {
				log.Fatalln("Error retrieveing key", err)
				os.Exit(2)
			}
			// decrypting files
			_, err = decryptFile(args[1], secret)
			if err != nil {
				log.Fatalln("Error decrypting file", args[1])
				os.Exit(2)
			}
			// cleanup
			os.Remove(args[1])
			fmt.Println("Decrypted", args[1])
		} else {
			showHelpMessage()
		}
	case "show":
		if len(args) == 2 {
			conf, err := GetConf(configFile)
			if err != nil {
				log.Fatal("Error: ", err)
				os.Exit(1)
			}
			// retriving secret from keyring
			secret, err := getKey(conf.Key)
			if err != nil {
				log.Fatalln("Error retrieveing key", err)
				os.Exit(2)
			}
			// decrypting file
			err = decryptFileString(args[1], secret)
			if err != nil {
				log.Fatalln("Error decrypting file", args[1])
				os.Exit(2)
			}
		} else {
			showHelpMessage()
		}
	case "version":
		fmt.Println("Version:", version)
		os.Exit(0)
	case "key":
		if len(args) == 2 {
			switch args[1] {
			case "show":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter a key name: ")
				keyName, _ := reader.ReadString('\n')

				key, err := getKey(keyName)
				if err != nil {
					log.Fatal("Error: ", err)
					os.Exit(2)
				}
				fmt.Println(key)
			case "delete":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter a key name: ")
				keyName, _ := reader.ReadString('\n')

				fmt.Print("Are you sure? Any encrypted files cannot be decrypted if key is deleted! (yes/No): ")
				input, _ := reader.ReadString('\n')

				if input == "yes\n" {
					err := deleteKey(keyName)
					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
					log.Print("Successfully deleted key")
					os.Exit(0)
				} else {
					log.Fatal("Skipping, only 'yes' is accepted as a confirmation")
					os.Exit(1)
				}
			case "create":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter a key name: ")
				keyName, _ := reader.ReadString('\n')
				keyName = strings.TrimSpace(keyName)

				fmt.Print("Enter a key: ")
				input1, _ := reader.ReadString('\n')

				fmt.Print("One more time please: ")
				input2, _ := reader.ReadString('\n')

				input1 = strings.TrimSpace(input1)
				input2 = strings.TrimSpace(input2)

				if input1 != input2 {
					log.Fatal("Error: entered keys do not match!")
					os.Exit(1)
				}

				err := setKey(keyName, input1)
				if err != nil {
					log.Fatal("Error: key creation failed,", err)
					os.Exit(1)
				}
			}
		} else {
			showHelpMessage()
		}
	default:
		showHelpMessage()
	}
}

func runit(executable string, args []string) {
	// fmt.Printf("Executing %s with args %s\n", executable, args)
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

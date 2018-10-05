package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/zalando/go-keyring"
)

const serviceNamePrefix = "keyrun"

func getKey(name string) (string, error) {

	name = strings.TrimSpace(name)

	user, err := getUserName()
	if err != nil {
		return "", err
	}

	secret, err := keyring.Get(serviceNamePrefix+"_"+name, user)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func setKey(name string, key string) error {

	name = strings.TrimSpace(name)

	user, err := getUserName()
	if err != nil {
		return err
	}

	fmt.Println("Creating key", serviceNamePrefix+"_"+name)

	err = keyring.Set(serviceNamePrefix+"_"+name, user, key)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func deleteKey(name string) error {

	name = strings.TrimSpace(name)

	user, err := getUserName()
	if err != nil {
		return err
	}

	err = keyring.Delete(serviceNamePrefix+"_"+name, user)
	if err != nil {
		return err
	}
	return nil
}

func getUserName() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

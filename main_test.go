package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestOpenAndReadFile(t *testing.T) {
	validFileName := "email-config.example.json"
	invalidFileName := "email-config.example.jsn"

	fileContents := []byte(`{
    "user": "someone@address.com",
    "password": "changeMe",
    "host": "mail.host.com",
    "port": "587",
    "sender-name": "Email Sender",
    "sender-email": "sender@address.com"
}`)

	validFile, err := openAndReadFile(validFileName)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	strippedFileContents := string(strings.Replace(string(validFile), "\t", "", -1))

	if string(fileContents) == strippedFileContents {
		t.Logf("openAndReadFile is successful with valid file")
	} else {
		t.Errorf("openAndReadFile is unsuccessful, expected %v, got %v", string(fileContents), strippedFileContents)
	}

	_, err = openAndReadFile(invalidFileName)

	if err != nil {
		t.Logf("openAndReadFile fails with inexisting file")
	}
}

func TestInitEmailConfig(t *testing.T) {
	validEmailConfigs := EmailConfig{
		User:        "someone@address.com",
		Password:    "changeMe",
		Host:        "mail.host.com",
		Port:        "587",
		SenderName:  "Email Sender",
		SenderEmail: "sender@address.com",
	}
	invalidEmailConfigs := EmailConfig{}

	validFileName := "email-config.example.json"
	validFile, err := openAndReadFile(validFileName)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	validFileData, err := initEmailConfig(validFile)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if validFileData == validEmailConfigs {
		t.Logf("initEmailConfig is successful with valid file")
	} else {
		t.Errorf("initEmailConfig is unsuccessful, expected %v, got %v", validEmailConfigs, validFileData)
	}

	inexistingFile, err := initEmailConfig(make([]byte, 0))

	if inexistingFile == invalidEmailConfigs && err != nil && err.Error() == "unexpected end of JSON input" {
		t.Logf("initEmailConfig fails with inexisting file")
	} else {
		t.Errorf("initEmailConfig is unsuccessful, expected %v, got %v", invalidEmailConfigs, inexistingFile)
	}
}
func TestInitEmailVariables(t *testing.T) {
	validEmailVariables := EmailVariables{
		Variables: []map[string]interface{}{0: {"RecipientName": "Email Receiver", "RecipientEmail": "receiver-1@another-address.com", "Variable1": "Value1", "Variable2": "Value2"}, 1: {"RecipientName": "Email Receiver", "RecipientEmail": "receiver-2@another-address.com", "Variable1": "Value1", "Variable2": "Value2"}, 2: {"RecipientName": "Email Receiver", "RecipientEmail": "receiver-3@another-address.com", "Variable1": "Value1", "Variable2": "Value2"}},
	}
	invalidEmailVariables := EmailVariables{}

	validFileName := "email-data.example.json"
	validFile, err := openAndReadFile(validFileName)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	validFileData, err := initEmailVariables(validFile)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if reflect.DeepEqual(validFileData, validEmailVariables) {
		t.Logf("initEmailVariables is successful with valid file")
	} else {
		t.Errorf("initEmailVariables is unsuccessful, expected %v, got %v", validEmailVariables, validFileData)
	}

	inexistingFile, err := initEmailVariables(make([]byte, 0))

	if reflect.DeepEqual(inexistingFile, invalidEmailVariables) && err != nil && err.Error() == "unexpected end of JSON input" {
		t.Logf("initEmailVariables fails with inexisting file")
	} else {
		t.Errorf("initEmailVariables is unsuccessful, expected %v, got %v", invalidEmailVariables, inexistingFile)
	}
}

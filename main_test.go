package main

import (
	"testing"
)

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

	validFile, err := initEmailConfig("email-config.example.json")

	if validFile == validEmailConfigs && err == nil {
		t.Logf("initEmailConfig config is successful with valid file")
	} else {
		t.Errorf("initEmailConfig config is unsuccessful, expected %v, got %v", validEmailConfigs, validFile)
	}

	inexistingFile, err := initEmailConfig("wrong-email-config.example.json")

	if inexistingFile == invalidEmailConfigs && err != nil && err.Error() == "open wrong-email-config.example.json: no such file or directory" {
		t.Logf("initEmailConfig config fails with inexisting file")
	} else {
		t.Errorf("initEmailConfig config is unsuccessful, expected %v, got %v", validEmailConfigs, inexistingFile)
	}
}

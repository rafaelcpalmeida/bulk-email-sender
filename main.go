package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	User           string `json:"user"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	SenderName     string `json:"sender-name"`
	SenderEmail    string `json:"sender-email"`
	RecipientName  string `json:"recipient-name"`
	RecipientEmail string `json:"recipient-email"`
}

type EmailVariables struct {
	Variables []map[string]interface{} `json:"variables"`
}

func main() {
	var emailConfig EmailConfig
	var emailVariables EmailVariables

	jsonFile, err := os.Open("email-config.json")

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &emailConfig)

	jsonFile, err = os.Open("email-data.json")

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &emailVariables)

	auth := smtp.PlainAuth("", emailConfig.User, emailConfig.Password, emailConfig.Host)

	for k := range emailVariables.Variables {

		vars := make(map[string]interface{})
		vars["RecipientName"] = emailConfig.RecipientName
		vars["RecipientEmail"] = emailConfig.RecipientEmail
		vars["SenderName"] = emailConfig.SenderName
		vars["SenderEmail"] = emailConfig.SenderEmail

		for _k := range emailVariables.Variables[k] {
			if str, ok := emailVariables.Variables[k][_k].(string); ok {
				fmt.Println()
				vars[_k] = str
			}
		}

		emailTemplate, err := template.ParseFiles("email.tmpl")

		if err != nil {
			fmt.Println(err)
		}

		var emailBytes bytes.Buffer
		err = emailTemplate.Execute(&emailBytes, vars)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Trying...")
		if err := smtp.SendMail(fmt.Sprintf("%s:%s", emailConfig.Host, emailConfig.Port), auth, emailConfig.User, []string{emailConfig.RecipientEmail}, []byte(emailBytes.String())); err != nil {
			fmt.Println("Error SendMail: ", err)
			os.Exit(1)
		}

		fmt.Println("Email Sent!")

		vars = nil
	}
}

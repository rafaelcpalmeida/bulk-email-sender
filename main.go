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
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	SenderName  string `json:"sender-name"`
	SenderEmail string `json:"sender-email"`
}

type EmailVariables struct {
	Variables []map[string]interface{} `json:"variables"`
}

func main() {
	var emailConfig EmailConfig
	var emailVariables EmailVariables

	setupVars(&emailConfig, &emailVariables)

	auth := configEmailAuth(emailConfig)

	for k := range emailVariables.Variables {
		vars := configureEmailVariables(emailConfig, k, emailVariables)

		emailTemplate := configureEmailTemplate(vars)

		fmt.Println("Sending email to: " + vars["RecipientEmail"].(string) + "...")
		if err := smtp.SendMail(fmt.Sprintf("%s:%s", emailConfig.Host, emailConfig.Port), auth, emailConfig.User, []string{vars["RecipientEmail"].(string)}, []byte(emailTemplate)); err != nil {
			printErrorAndDie(err)
		}

		fmt.Println("Email sent!")
		fmt.Println()

		vars = nil
	}
}

func printErrorAndDie(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func setupVars(emailConfig *EmailConfig, emailVariables *EmailVariables) {
	jsonFile, err := os.Open("email-config.json")

	if err != nil {
		printErrorAndDie(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &emailConfig)

	jsonFile, err = os.Open("email-data.json")

	if err != nil {
		printErrorAndDie(err)
	}

	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &emailVariables)
}

func configEmailAuth(emailConfig EmailConfig) smtp.Auth {
	return smtp.PlainAuth("", emailConfig.User, emailConfig.Password, emailConfig.Host)
}

func configureEmailVariables(emailConfig EmailConfig, k int, emailVariables EmailVariables) map[string]interface{} {
	vars := make(map[string]interface{})

	vars["SenderName"] = emailConfig.SenderName
	vars["SenderEmail"] = emailConfig.SenderEmail

	for _k := range emailVariables.Variables[k] {
		if str, ok := emailVariables.Variables[k][_k].(string); ok {
			vars[_k] = str
		}
	}

	return vars
}

func configureEmailTemplate(vars map[string]interface{}) string {
	template, err := template.New("email.tmpl").Funcs(template.FuncMap{
		"emailAddressStructure": func(str string) template.HTML {
			return template.HTML(fmt.Sprintf("<%s>", str))
		},
	}).ParseFiles("email.tmpl")

	if err != nil {
		printErrorAndDie(err)
	}

	var emailBytes bytes.Buffer
	err = template.Execute(&emailBytes, vars)

	if err != nil {
		printErrorAndDie(err)
	}

	return emailBytes.String()
}

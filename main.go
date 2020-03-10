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

// EmailConfig struct holds every data required to authenticate on an SMTP Server and also identifies the sender's name and email
type EmailConfig struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	SenderName  string `json:"sender-name"`
	SenderEmail string `json:"sender-email"`
}

// EmailVariables struct holds every variable name and its corresponding value it's meant to be replaced on every email body sent.
type EmailVariables struct {
	Variables []map[string]interface{} `json:"variables"`
}

func main() {
	emailConfigData, err := openAndReadFile("email-config.json")

	if err != nil {
		printErrorAndDie("Error reading email configuration: ", err)
	}

	emailVariablesData, err := openAndReadFile("email-data.json")
	if err != nil {
		printErrorAndDie("Error reading variable configuration: ", err)
	}

	emailConfig, err := initEmailConfig(emailConfigData)

	if err != nil {
		printErrorAndDie("Error initiating email configuration: ", err)
	}

	emailVariables, err := initEmailVariables(emailVariablesData)
	if err != nil {
		printErrorAndDie("Error initiating variable configuration: ", err)
	}

	smtpAuth := configEmailAuth(emailConfig)

	err = emailVariables.sendEmails(emailConfig, smtpAuth)

	if err != nil {
		printErrorAndDie("Error sending email: ", err)
	}
}

func printErrorAndDie(description string, err error) {
	fmt.Println(description + err.Error())
	os.Exit(1)
}

func openAndReadFile(filename string) ([]byte, error) {
	jsonFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}

func initEmailConfig(fileData []byte) (EmailConfig, error) {
	var emailConfig EmailConfig

	if err := json.Unmarshal(fileData, &emailConfig); err != nil {
		return EmailConfig{}, err
	}

	return emailConfig, nil
}

func initEmailVariables(fileData []byte) (EmailVariables, error) {
	var emailVariables EmailVariables

	if err := json.Unmarshal(fileData, &emailVariables); err != nil {
		return EmailVariables{}, err
	}

	return emailVariables, nil
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

func configureEmailTemplate(vars map[string]interface{}) (string, error) {
	template, err := template.New("email.tmpl").Funcs(template.FuncMap{
		"emailAddressStructure": func(str string) template.HTML {
			return template.HTML(fmt.Sprintf("<%s>", str))
		},
	}).ParseFiles("email.tmpl")

	if err != nil {
		return "", err
	}

	var emailBytes bytes.Buffer
	err = template.Execute(&emailBytes, vars)

	if err != nil {
		return "", err
	}

	return emailBytes.String(), nil
}

func (emailVariables *EmailVariables) sendEmails(emailConfig EmailConfig, smtpAuth smtp.Auth) error {
	for k := range emailVariables.Variables {
		vars := configureEmailVariables(emailConfig, k, *emailVariables)

		emailTemplate, err := configureEmailTemplate(vars)

		if err != nil {
			return err
		}

		fmt.Println("Sending email to: " + vars["RecipientEmail"].(string) + "...")
		if err := smtp.SendMail(fmt.Sprintf("%s:%s", emailConfig.Host, emailConfig.Port), smtpAuth, emailConfig.User, []string{vars["RecipientEmail"].(string)}, []byte(emailTemplate)); err != nil {
			return err
		}

		fmt.Println("Email sent!")
		fmt.Println()

		vars = nil
	}

	return nil
}

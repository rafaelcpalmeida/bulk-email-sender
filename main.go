package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
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
	emailConfig, err := initEmailConfig("email-config.json")
	if err != nil {
		printErrorAndDie("Error initiating email configuration: ", err)
	}

	emailVariables, err := initEmailVariables("email-data.json")
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

func initEmailConfig(filename string) (EmailConfig, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		return EmailConfig{}, err
	}

	defer jsonFile.Close()

	var emailConfig EmailConfig

	if err = json.NewDecoder(jsonFile).Decode(&emailConfig); err != nil {
		return EmailConfig{}, err
	}

	return emailConfig, nil
}

func initEmailVariables(filename string) (EmailVariables, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		return EmailVariables{}, err
	}

	defer jsonFile.Close()

	var emailVariables EmailVariables

	if err = json.NewDecoder(jsonFile).Decode(&emailVariables); err != nil {
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

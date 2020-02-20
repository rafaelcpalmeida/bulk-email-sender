[![Go Report Card](https://goreportcard.com/badge/github.com/rafaelcpalmeida/bulk-email-sender)](https://goreportcard.com/report/github.com/rafaelcpalmeida/bulk-email-sender)

# Bulk Email Sender

Bulk Email Sender is a small utility program that allows you to send bulk emails (duh!) using SMTP authentication.

## Installation

Build the binary using Go's compiler

## Windows
#### 32 bits

```bash
env GOOS=windows GOARCH=386 go build .
```
#### 64 bits

```bash
env GOOS=windows GOARCH=amd64 go build .
```
## Linux
#### 32 bits

```bash
env GOOS=linux GOARCH=386 go build .
```
#### 64 bits

```bash
env GOOS=linux GOARCH=amd64 go build .
```
## macOS
#### 32 bits

```bash
env GOOS=darwin GOARCH=386 go build .
```
#### 64 bits

```bash
env GOOS=darwin GOARCH=amd64 go build .
```

## Usage

Please keep in mind that the following file structure is required:

```bash
├── bulk-email-sender
├── email-config.json
├── email-data.json
└── email.tmpl
```

The program will panic if it fails to load any of the files above.

### Configuration instructions
Given the file structure above, the file:
* **email-config.json** is where you should configure your SMTP authentication details
* **email-data.json** is where you should define your variables and their respective values
* **email.tmpl** is where you should define your email text and reference the variables you wish to replace. Keep in mind that you can **ONLY** change from __subject__ onwards

### Example
Given the following files:
**email-data.json**
```json 
{
    "user": "someone@gmail.com",
    "password": "mySuperP@ssword",
    "host": "smtp.gmail.com",
    "port": "587",
    "sender-name": "John Appleseed",
    "sender-email": "someone@gmail.com"
}
```
**email-data.json**
```json 
{
    "variables": [
        {
            "RecipientName": "Rafael Almeida",
            "RecipientEmail": "rafael@gmail.com",
            "Variable1": "Value1",
            "Variable2": "Value2"
        }
    ]
}
```
**email-data.json**
```tmpl 
To: "{{.RecipientName}}" {{emailAddressStructure .RecipientEmail}}
From: "{{.SenderName}}" {{emailAddressStructure .SenderEmail}}
Subject: Testing Bulk Email Sender
Hey,

Please note that:

- Variable 1 is {{.Variable1}}.
- Variable 2 is {{.Variable2}}.
```

Will render the following email:
```text
Hey,

Please note that:

- Variable 1 is Value1.
- Variable 2 is Value2.
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)


**Made with :heart: in Portugal**

**Software livre c\*ralho! :v:**
package main

import (
    "./smtpd"
)

func main() {
    smtpd.StartSMTPServer(":465", "zencu.be");
}
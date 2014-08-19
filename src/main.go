package main

import (
    "./smtpd"
)

func main() {
    smtpd.StartSMTPServer(":25", "zencu.be");
}
package main

import (
    "flag"
    "./smtpd"
)

func main() {
    smtpAddr := flag.String("smtp_addr", "127.0.0.1:25", "address and port to listen for SMTP")
    domain := flag.String("domain", "local", "The domain this mail server will accept mails for")
    flag.Parse()

    smtpd.StartSMTPServer(*smtpAddr, *domain);
}
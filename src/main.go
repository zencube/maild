package main

import (
    "fmt"
    "log"
    "flag"
    "./smtpd"
    "crypto/tls"
)

func main() {
    smtpAddr := flag.String("smtp_addr", "127.0.0.1", "address (without the port) to listen for SMTP / SMTPS")
    domain := flag.String("domain", "local", "The domain this mail server will accept mails for")
    certPath := flag.String("cert", "./server.crt", "The SSL certificate (can be self-signed) for TLS")
    certKey := flag.String("certkey", "./server.key", "The SSL certificate key for the certificate")
    flag.Parse()

    cert, err := tls.LoadX509KeyPair(*certPath, *certKey)
    if err != nil {
        log.Printf("SMTP: No TLS support, error: %v", err)
    }
    tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.VerifyClientCertIfGiven, ServerName: *domain}

    var addr string
    addr = fmt.Sprintf("%s:25", *smtpAddr)
    smtpd.StartSMTPServer(addr, *domain, &tlsConfig)

}
package smtpd

import(
    "fmt"
    "net"
    "log"
    "time"
    "bufio"
    "errors"
    "strconv"
    "strings"
    "crypto/tls"
)

type Client struct {
    Connection net.Conn
    Reader *bufio.Reader
    Writer *bufio.Writer
    Addr string
}

func StartSMTPServer(addr string, domain string) {
    listener, err := net.Listen("tcp", addr);
    if err != nil {
        log.Fatalf("STMP: Cannot start server: %v", err)
    }
    
    for {
        conn, err := listener.Accept();
        if err != nil {
            log.Printf("SMTP: Client connection error: %v", err)
            continue
        }
        
        log.Printf("SMTP: Client %s connected.", conn.RemoteAddr().String());
        go handleClient(&Client{
            Connection: conn, 
            Reader:     bufio.NewReader(conn), 
            Writer:     bufio.NewWriter(conn),
            Addr:       conn.RemoteAddr().String(),
        }, domain)
    }
}

func handleClient(client *Client, domain string) {
    client.Writer.WriteString(fmt.Sprintf("220 %s ESMTP service ready\n", domain))
    client.Writer.Flush()
    loop:
    for {
        cmd, err := readClientCommand(client)
        if err != nil {
            log.Printf("SMTP: Connection with %s ended: %v", client.Addr, err)
            break
        }
        
        cmd = strings.ToUpper(strings.Trim(cmd, "\r\n"))
        log.Printf("SMTP: Received %s from %s", cmd, client.Addr)
        switch {
            case strings.Index(cmd, "HELO") == 0: // || strings.Index(cmd, "EHLO") == 0:
                var response string
                if len(cmd) > 6 {
                    response = fmt.Sprintf("250 %s Hello %s\n", domain, cmd[5:])
                } else {
                    response = fmt.Sprintf("250 %s Hello\n", domain)
                }
                client.Writer.WriteString(response)
                break
            /**/
            case strings.Index(cmd, "EHLO") == 0:
                var response string
                if len(cmd) > 6 {
                    response = fmt.Sprintf("250-%s %s Hello\n", cmd[5:], domain)
                } else {
                    response = fmt.Sprintf("250-%s Hello\n", domain)
                }
                response = response + "250-AUTH LOGIN\n250 OK\n"
                client.Writer.WriteString(response)
                break
            case strings.Index(cmd, "QUIT") == 0:
                client.Writer.WriteString("221 Good bye\n")
                client.Writer.Flush()
                break loop
            case strings.Index(cmd, "MAIL FROM:") == 0:
                if len(cmd) > 10 && strings.Index(cmd, "<") > 0 {
                    sender := cmd[strings.Index(cmd, "<")+1 : strings.Index(cmd, ">")]
                    log.Printf("SMTP: Sender %s", sender)
                    client.Writer.WriteString("250 OK\n")
                } else {
                    log.Printf("SMTP: Encountered invalid MAIL FROM: %s", cmd)
                    client.Writer.WriteString("500 Invalid sender\n")
                }
                break
            case strings.Index(cmd, "RCPT TO:") == 0:
                if len(cmd) > 8 && strings.Index(cmd, "<") > 0 {
                    rcpt := cmd[strings.Index(cmd, "<")+1 : strings.Index(cmd, ">")]
                    log.Printf("SMTP: Receipient %s", rcpt)
                    client.Writer.WriteString("250 OK\n")
                } else {
                    log.Printf("SMTP: Encountered invalid MAIL FROM: %s", cmd)
                    client.Writer.WriteString("500 Invalid receipient\n")
                }
                break
            case strings.Index(cmd, "DATA") == 0:
        		client.Connection.SetDeadline(time.Now().Add(120 * time.Second))
                log.Printf("SMTP: Receiving data...")
                client.Writer.WriteString("354 Input mail data\n")
                client.Writer.Flush()
                var headers string
                var body string
                inBody := false
                for {
            		line, err := client.Reader.ReadString('\n')
            		line = strings.Trim(line, "\r\n")

            		if err != nil {
            		    log.Printf("SMTP: Data error: %v", err)
            		    break
            		}

            		if line == "" {
            		    inBody = true;
            		    continue
            		}

            		if line == "." {
            		    client.Writer.WriteString("250 OK\n")
            		    log.Printf("Done.")
            		    break
            		}

            		if inBody == true {
            		    body = body + line + "\n"
            		} else {
            		    headers = headers + line + "\n"
            		}
                }
                
                log.Printf(fmt.Sprintf("SMTP: Headers:\n%s\nBody:\n%s", headers, body));
                break
            default:
                client.Writer.WriteString("500 unrecognized command\n")
                client.Writer.Flush()
                break loop
        }
        client.Writer.Flush()
    }
    client.Connection.Close()
}

func readClientCommand(client *Client) (string, error) {
    var command string
    var err error
    var line string
    suffix := "\r\n"

	for err == nil {
		client.Connection.SetDeadline(time.Now().Add(60 * time.Second))
		line, err = client.Reader.ReadString('\n')
		if line != "" {
			command = command + line
			if len(line) > 10240000 {
				err = errors.New("Maximum DATA size exceeded (" + strconv.Itoa(10240000) + ")")
				return command, err
			}
		}
		if err != nil {
			break
		}
		if strings.HasSuffix(command, suffix) {
			break
		}
	}
	return command, err
}
package twosshlibSsh

import (
	"fmt"

	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// executeSsh - подключится ао SSH (к cisco), выполнить команду и вернуть что она ответила.
//
// Пример использования:
//
// var (
//
//		User     string = "developer"
//		Password string = "C1sco12345"
//		hostname string = "192.168.1.11"
//	     port     string = "2222"
//		cmds            = []string{"show ip route | i 0.0.0.0/0", "show ip arp"}
//
// )
//
// results := make(chan string, 100)
//
//	config := &ssh.ClientConfig{
//		 User:            User,
//		 HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		 Auth: []ssh.AuthMethod {
//				ssh.Password(Password),
//		 },
//	}
//
//	results := executeCmd(hostname, port, cmds, config)
func ExecuteSsh(hostname string, port string, cmds []string, config *ssh.ClientConfig) (string, error) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Создадим подключение с нужными параметрами.
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		fmt.Println("Не могу подключится :-(")
		return "", &ssh.OpenChannelError{}
	}
	defer conn.Close()

	// Пробуем подключится
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		return "", &ssh.OpenChannelError{}
	}
	defer session.Close()

	// You can use session.Run() here but that only works
	// if you need a run a single command or you commands
	// are independent of each other.
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return "", fmt.Errorf("request for pseudo terminal failed: %s", err)
	}
	stdBuf, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("request for stdout pipe failed: %s", err)
	}
	stdinBuf, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("request for stdin pipe failed: %s", err)
	}
	err = session.Shell()
	if err != nil {
		return "", fmt.Errorf("failed to start shell: %s", err)
	}

	var cmd_output string

	for _, cmd := range cmds {
		stdinBuf.Write([]byte(cmd + "\n"))
		for {
			stdoutBuf := make([]byte, 1000000)
			time.Sleep(time.Millisecond * 700)
			byteCount, err := stdBuf.Read(stdoutBuf)
			if err != nil {
				log.Fatal(err)
			}
			cmd_output += string(stdoutBuf[:byteCount])
			if !(strings.Contains(string(stdoutBuf[:byteCount]), "More")) {
				break
			}
			stdinBuf.Write([]byte(" "))

		}
	}

	return cmd_output, nil
}

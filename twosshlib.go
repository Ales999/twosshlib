package twosshlib

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	myssh "github.com/ales999/twosshlib/internal"
)

// Выполнить подключение с первому, далее ко второму, на втором выполнить команду,
// отключится от второго, выполнить команду на первом и отключится.
func FirstTwoNextOne(oneHostIp []string, onePort, twoHostIp, twoPort, user, password, oneCommand, twoCommand string) (string, error) {

	//
	var sb strings.Builder
	// Предварительно выделяем сразу память
	sb.Grow(len(twoHostIp) + 1 + len(twoPort) + len(oneCommand) + len(twoCommand))
	// Собираем строку что будем выполнять на первом роутере
	sb.WriteString("ssh -l ")
	sb.WriteString(user)
	sb.WriteString(" ")
	sb.WriteString(password)
	sb.WriteString(" ")
	sb.WriteString(twoCommand)
	sb.WriteString(" ")
	sb.WriteString("exit ")
	sb.WriteString(oneCommand)
	sb.WriteString("exit")

	var cmdsstr = sb.String()

	fmt.Println("Отладка cmds:", cmdsstr)

	var cmds []string
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	return myssh.ExecuteSsh(oneHostIp[0], onePort, cmds, config)
}

// Выполнить подключение с первому, выполнить на первом команду, далее
// подключится ко второму, выполнить команду на втором.
// Далее последовательно отключится от обоих.
func FirstOneNextTwo() (string, error) {

	return "", nil
}

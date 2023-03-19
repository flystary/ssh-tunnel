package tunnel

import (
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

const Version = "1.0"

func forward(serverAddr, remoteAddr string, localConn net.Conn, sshConfig *ssh.ClientConfig) {
	sshClientConn, err := ssh.Dial("tcp", serverAddr, sshConfig)
	if err != nil {
		fmt.Printf("ssh.Dial failed: %s\n", err)
	}

	sshConn, err := sshClientConn.Dial("tcp", remoteAddr)

	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			fmt.Printf("io.Copy failed: %s\n", err)
		}
	}()

	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			fmt.Printf("io.Copy failed: %v\n", err)
		}
	}()
}

func Tunnel(username, password, serverAdrr, remoteAddr, localAddr string) {
	// 设置SSH配置
	fmt.Printf("%s，服务器：%s；远程：%s；本地：%s\n", "设置SSH配置", serverAdrr, remoteAddr, localAddr)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// 设置本地监听器
	localListener, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Printf("net.Listen failed: %v\n", err)
	}

	for {
		// 设置本地
		localConn, err := localListener.Accept()
		if err != nil {
			fmt.Printf("localListener.Accept failed: %v\n", err)
		}
		go forward(serverAdrr, remoteAddr, localConn, config)
	}
}

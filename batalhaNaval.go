package main

import (
	"fmt"
	"net"
	"os"
)

type Ship struct {
	Name string
	Size int8
}

type Board struct {
	Grid  [10][10]string
	Ships []Ship
}

func main() {

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Enviar mensagem
	message := []byte("hello")
	_, err = conn.Write(message)

	// Receber mensagem
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	fmt.Println(string(buffer[:n]))

}

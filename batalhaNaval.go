package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

type Ship struct {
	Name string
	Size int
}

type Board struct {
	Grid  [10][10]string
	Ships []Ship
}

type Player struct {
	Board Board
	Ships []Ship
}

func placeShipsRandomly(board *Board) {
	ships := []Ship{{Name: "Submarino", Size: 1}, {Name: "Submarino", Size: 1}, {Name: "Submarino", Size: 1}, {Name: "Cruzador", Size: 2}, {Name: "Cruzador", Size: 2}, {Name: "Porta-aviões", Size: 3}}
	for _, ship := range ships {
		for {
			row := rand.Intn(10)
			col := rand.Intn(10)
			direction := rand.Intn(2) // 0 para horizontal, 1 para vertical

			if direction == 0 {
				if col+ship.Size <= 10 {
					for k := 0; k < ship.Size; k++ {
						board.Grid[row][col+k] = ship.Name
					}
					break
				}
			} else {
				if row+ship.Size <= 10 {
					for k := 0; k < ship.Size; k++ {
						board.Grid[row+k][col] = ship.Name
					}
					break
				}
			}
		}
	}
}

func attack(attacker *Player, target *Player) {
	//row := rand.Intn(10)
	//col := rand.Intn(10)
	targetRow := rand.Intn(10)
	targetCol := rand.Intn(10)

	if attacker.Board.Grid[targetRow][targetCol] == "Água" {
		fmt.Println("Ataque na água!")
	} else {
		fmt.Println("Ataque em um navio!")
		if target.Board.Grid[targetRow][targetCol] == "Navio" {
			target.Board.Grid[targetRow][targetCol] = "Acertado"
			if target.Board.checkAllShipsSunk() {
				fmt.Println("Jogo terminou Todos os navios do oponente foram acertados.")
			}
		}
	}
}

func (b *Board) checkAllShipsSunk() bool {
	for _, ship := range b.Ships {
		for _, row := range b.Grid {
			for _, cell := range row {
				if strings.Contains(cell, ship.Name) && cell != "Acertado" {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	player1 := Player{Board: Board{}, Ships: []Ship{}}
	player2 := Player{Board: Board{}, Ships: []Ship{}}

	placeShipsRandomly(&player1.Board)
	placeShipsRandomly(&player2.Board)

	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	checkError(err)

	conn1, err := net.DialUDP("udp", nil, addr1)
	checkError(err)
	defer conn1.Close()

	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:8081")
	checkError(err)

	conn2, err := net.DialUDP("udp", nil, addr2)
	checkError(err)
	defer conn2.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		if player1Turn := rand.Intn(2) == 0; player1Turn {
			fmt.Print("Player 1, enter command: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			_, err = conn1.Write([]byte(text))
			checkError(err)

			buffer := make([]byte, 1024)
			n, _, err := conn1.ReadFromUDP(buffer)
			checkError(err)

			response := string(buffer[:n])
			fmt.Println("Player 1 received:", response)

			if strings.Contains(response, "acertou") {
				continue
			}

			if strings.Contains(response, "terminou") {
				break
			}
		} else {
			fmt.Print("Player 2, enter command: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			_, err = conn2.Write([]byte(text))
			checkError(err)

			buffer := make([]byte, 1024)
			n, _, err := conn2.ReadFromUDP(buffer)
			checkError(err)

			response := string(buffer[:n])
			fmt.Println("Player 2 received:", response)

			if strings.Contains(response, "acertou") {
				continue
			}

			if strings.Contains(response, "terminou") {
				break
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

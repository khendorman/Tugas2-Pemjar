package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients = make(map[net.Conn]Client)
	mu      sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Gagal menjalankan server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server TCP berjalan di port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error koneksi:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer func() {
		mu.Lock()
		client, exists := clients[conn]
		delete(clients, conn)
		mu.Unlock()

		if exists {
			broadcast("*** " + client.name + " keluar ***")
			fmt.Println(client.name, "terputus")
		}

		conn.Close()
	}()

	fmt.Fprintln(conn, "Masukkan nama Anda:")
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	mu.Lock()
	clients[conn] = Client{conn: conn, name: name}
	mu.Unlock()

	fmt.Printf("%s bergabung\n", name)
	fmt.Fprintf(conn, "Selamat datang, %s! Pesan Anda akan dilihat semua client lain.\n", name)

	broadcast("*** " + name + " bergabung ***")

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		fmt.Printf("%s: %s\n", name, message)
		broadcast(name + ": " + message)
	}
}

func broadcast(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, client := range clients {
		fmt.Fprintln(client.conn, message)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Gagal terhubung ke server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Terhubung ke server TCP!")
	fmt.Println("Ketik pesan lalu tekan Enter.")
	fmt.Println("Ketik /keluar untuk keluar.")
	fmt.Println("------------------------------")

	go func() {
		reader := bufio.NewScanner(conn)
		for reader.Scan() {
			fmt.Println(reader.Text())
		}
		os.Exit(0)
	}()

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		message := strings.TrimSpace(input.Text())

		if message == "/keluar" {
			fmt.Println("Keluar dari chat...")
			return
		}

		if message == "" {
			continue
		}

		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			fmt.Println("Gagal mengirim pesan:", err)
			return
		}
	}
}

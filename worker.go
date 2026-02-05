package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")

	if err != nil {
		fmt.Fprintf(os.Stderr, "A problem occured while connecting to the server")	
	}

	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')

		if msg == "quit\n" {
			break
		}

		_, err := conn.Write([]byte(msg))		

		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		}
		
		scanner := bufio.NewScanner(conn)
		ok := scanner.Scan()

		if !ok {
			fmt.Println("Scan error")
		}
		
		fmt.Println("> " + scanner.Text())
	}
}




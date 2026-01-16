package main

import "fmt"

func main() {
	server := NewAPIServer(":8080")
	err := server.Run()
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}
}

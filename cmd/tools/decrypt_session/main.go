package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/securecookie"
)

func init() {
	// Register common types often used in gin-contrib/sessions
	gob.Register(map[any]any{})
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("   Connectify Session Decryption Tool   ")
	fmt.Println("========================================")

	// 1. Get Cookie Value
	fmt.Print("➤ Please enter the Cookie string: ")
	rawCookie, _ := reader.ReadString('\n')
	rawCookie = cleanInput(rawCookie)

	// 2. Get Hash Key
	fmt.Print("➤ Please enter the Hash Key (default [Ktsoator]): ")
	hashKeyStr, _ := reader.ReadString('\n')
	hashKeyStr = strings.TrimSpace(hashKeyStr)
	if hashKeyStr == "" {
		hashKeyStr = "Ktsoator"
	}
	hashKey := []byte(hashKeyStr)

	// 3. Get Block Key
	fmt.Print("➤ Please enter the Block Key (default [np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B]): ")
	blockKeyStr, _ := reader.ReadString('\n')
	blockKeyStr = strings.TrimSpace(blockKeyStr)
	if blockKeyStr == "" {
		blockKeyStr = "np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B"
	}
	blockKey := []byte(blockKeyStr)

	data := make(map[any]any)
	s := securecookie.New(hashKey, blockKey)
	err := s.Decode("connectify", rawCookie, &data)

	if err != nil {
		fmt.Printf("\nDecryption failed: %v\n", err)
		return
	}

	fmt.Println("\nDecryption Successful!")
	fmt.Println("----------------------------------------")

	if len(data) == 0 {
		fmt.Println("Session is empty")
	} else {
		for k, v := range data {
			fmt.Printf("%v: %v (%T)\n", k, v, v)
		}
	}
	fmt.Println("----------------------------------------")
}

func cleanInput(input string) string {
	input = strings.TrimSpace(input)
	if idx := strings.Index(input, "connectify="); idx != -1 {
		input = input[idx+11:]
	}
	if idx := strings.Index(input, ";"); idx != -1 {
		input = input[:idx]
	}
	return strings.TrimSpace(input)
}

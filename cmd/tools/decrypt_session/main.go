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

	// In Redis mode, gin-contrib/sessions stores the Session ID as a GOB-encoded string.
	// In Cookie mode, it stores the data as a GOB-encoded map[interface{}]interface{}.
	// We try string first (Redis), then fallback to map (Cookie).

	s := securecookie.New(hashKey, blockKey)

	// Try String (Redis/Session ID)
	var sessionID string
	err := s.Decode("connectify", rawCookie, &sessionID)
	if err == nil {
		fmt.Println("\nDecryption Successful!")
		fmt.Println("----------------------------------------")
		fmt.Println("Store Type: Redis/Other (Data stored on Server)")
		fmt.Printf("Session ID: %s\n", sessionID)
		fmt.Println("\n[Note] In Redis mode, the cookie only contains the Session ID.")
		fmt.Println("To see actual data, you must query Redis using this ID.")
		fmt.Println("----------------------------------------")
		return
	}

	// Try Map (Cookie Store)
	data := make(map[interface{}]interface{})
	err = s.Decode("connectify", rawCookie, &data)
	if err == nil {
		fmt.Println("\nDecryption Successful!")
		fmt.Println("----------------------------------------")
		fmt.Println("Store Type: Cookie Store (Data stored in Cookie)")
		for k, v := range data {
			fmt.Printf("%v: %v (%T)\n", k, v, v)
		}
		fmt.Println("----------------------------------------")
		return
	}

	fmt.Printf("\nDecryption failed: %v\n", err)
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

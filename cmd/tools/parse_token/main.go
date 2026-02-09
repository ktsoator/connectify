package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ktsoator/connectify/internal/web/user"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("      Connectify JWT Parse Tool         ")
	fmt.Println("========================================")

	fmt.Print("Please enter JWT token: ")
	rawToken, _ := reader.ReadString('\n')
	tokenStr := cleanToken(rawToken)
	if tokenStr == "" {
		fmt.Println("empty token")
		os.Exit(1)
	}

	fmt.Print("Please enter JWT secret (default [Ktsoator]): ")
	secretStr, _ := reader.ReadString('\n')
	secretStr = strings.TrimSpace(secretStr)
	if secretStr == "" {
		secretStr = "Ktsoator"
	}

	claims := user.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		return []byte(secretStr), nil
	})

	fmt.Println("----------------------------------------")
	if err != nil {
		fmt.Printf("Parse result: failed (%v)\n", err)
	} else {
		fmt.Println("Parse result: success")
	}
	if token != nil {
		fmt.Printf("Signature valid: %v\n", token.Valid)
		fmt.Printf("Algorithm: %v\n", token.Method.Alg())
	}

	fmt.Printf("UserId: %d\n", claims.UserId)
	fmt.Printf("UserEmail: %s\n", claims.UserEmail)
	fmt.Printf("UserAgent: %s\n", claims.UserAgent)
	printTime("IssuedAt", claims.IssuedAt)
	printTime("NotBefore", claims.NotBefore)
	printTime("ExpiresAt", claims.ExpiresAt)

	if len(claims.Audience) > 0 {
		fmt.Printf("Audience: %v\n", []string(claims.Audience))
	}
	if claims.Issuer != "" {
		fmt.Printf("Issuer: %s\n", claims.Issuer)
	}
	if claims.Subject != "" {
		fmt.Printf("Subject: %s\n", claims.Subject)
	}
	if claims.ID != "" {
		fmt.Printf("JWT ID: %s\n", claims.ID)
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) == 3 {
		var header any
		if b, decodeErr := jwt.NewParser().DecodeSegment(parts[0]); decodeErr == nil {
			_ = json.Unmarshal(b, &header)
		}
		if header != nil {
			if headerJSON, marshalErr := json.MarshalIndent(header, "", "  "); marshalErr == nil {
				fmt.Println("Header:")
				fmt.Println(string(headerJSON))
			}
		}
	}

	fmt.Println("----------------------------------------")
}

func cleanToken(input string) string {
	v := strings.TrimSpace(input)
	if strings.HasPrefix(v, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(v, "Bearer "))
	}
	if strings.HasPrefix(v, "Jwt-Token:") {
		return strings.TrimSpace(strings.TrimPrefix(v, "Jwt-Token:"))
	}
	if strings.HasPrefix(v, "Jwt-Token=") {
		return strings.TrimSpace(strings.TrimPrefix(v, "Jwt-Token="))
	}
	return v
}

func printTime(name string, d *jwt.NumericDate) {
	if d == nil {
		return
	}
	t := d.Time
	fmt.Printf("%s: %s (in %s)\n", name, t.Format(time.RFC3339), time.Until(t).Round(time.Second))
}

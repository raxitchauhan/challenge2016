package util

import (
	"fmt"
)

func GetMainMenu() {
	fmt.Println("")
	fmt.Println("1. List all distributors")
	fmt.Println("2. Add distributor with required permission")
	fmt.Println("3. Check permission for a distributor")
	fmt.Println("4. Create a network of distributors")
	fmt.Println("5. Back to the main menu")
	fmt.Println("6. Exit")
	fmt.Println("")
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

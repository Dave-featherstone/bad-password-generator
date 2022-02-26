package main

import (
	"fmt"
)

func main() {

	fmt.Println("Bad Password Generator")

	pwList := readPasswordFile()

	for {
		var menuOption int
		fmt.Println("Welcome to your comamnd line password manager")
		fmt.Println("")
		fmt.Println("MENU")
		fmt.Println("1 - Show passwords")
		fmt.Println("2 - Add Manual Password")
		fmt.Println("3 - Add Generated Password")
		fmt.Println("4 - Quit")

		fmt.Scanln(&menuOption)

		if menuOption == 1 {
			displayPasswords(pwList)
			continue
		} else if menuOption == 2 {
			newEntry := newPasswordManual()
			pwList.PasswordList = append(pwList.PasswordList, newEntry)
			savePasswordFile(pwList)
		} else if menuOption == 3 {
			newEntry := newPasswordAuto()
			pwList.PasswordList = append(pwList.PasswordList, newEntry)
			savePasswordFile(pwList)
		} else if menuOption == 4 {
			break
		} else {
			fmt.Println("That wasn't an option")
		}
	}

}

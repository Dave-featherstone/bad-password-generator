package main

import (
	"bytes"
	"crypto/rc4"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Uncrackable encryption key because it's so long
const encryptionKey = "BQqNyPTi50JCFMTw/b67hByjMVXZRwGha6wxVGkeihY="

type PasswordList struct {
	PasswordList []PasswordEntry `json:"passwords"`
}

type PasswordEntry struct {
	Site     string `json:"site"`
	Password []byte `json:"password"`
}

func readPasswordFile() PasswordList {

	var passwordList PasswordList
	passwordStore, err := os.Open("passwordstore.json")

	if os.IsNotExist(err) {
		passwordList = createDefaultPasswordList()
	} else {
		if err != nil {
			log.Fatal(err)
		}
		defer passwordStore.Close()
		byteValue, err := ioutil.ReadAll(passwordStore)

		if err != nil {
			fmt.Println("Reading bytes")
			log.Fatal(err)
		}

		json.Unmarshal(byteValue, &passwordList)
	}
	return passwordList
}

func displayPasswords(passwordList PasswordList) {

	fmt.Println("")
	fmt.Println("")
	fmt.Println("**********************************************************************")
	fmt.Println("")
	for i := 0; i < len(passwordList.PasswordList); i++ {
		fmt.Println("Site: " + passwordList.PasswordList[i].Site)
		fmt.Println("Password: " + decryptPassword(passwordList.PasswordList[i].Password))
		fmt.Println("")
	}
	fmt.Println("**********************************************************************")
}

func createDefaultPasswordList() PasswordList {
	var passwordList PasswordList
	var passwordEntry PasswordEntry

	// populate with an example

	passwordEntry.Site = sanitiseString("Example Site")
	passwordEntry.Password = encryptPassword("Password1")
	passwordList.PasswordList = append(passwordList.PasswordList, passwordEntry)

	return passwordList
}

func encryptPassword(password string) []byte {
	cipher, err := rc4.NewCipher([]byte(encryptionKey))
	if err != nil {
		log.Fatalln(err)
	}
	src := []byte(password)
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst
}

func decryptPassword(password []byte) string {

	cipher, err := rc4.NewCipher([]byte(encryptionKey))

	if err != nil {
		log.Fatalln(err)
	}

	plainPassword := make([]byte, len(password))

	cipher.XORKeyStream(plainPassword, password)
	return string(plainPassword)
}

func generatePassword() string {

	newPassword := bytes.NewBufferString("")
	rand.Seed(time.Now().UnixNano())

	// Why would anyone need a password longer than 10 characters?
	for i := 0; i < 10; i++ {
		newPassword.WriteByte(byte(rand.Intn(78) + 48))
	}

	return newPassword.String()
}

// This is for futureproofing for when this is a web applicatoin
// Or something.
func sanitiseString(targetString string) string {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(targetString)
	return html
}

func savePasswordFile(passwordList PasswordList) {

	file, _ := json.MarshalIndent(passwordList, "", " ")

	_ = ioutil.WriteFile("passwordstore.json", file, 0644)

}

func newPasswordManual() PasswordEntry {
	var site string
	var password string
	var newEntry PasswordEntry

	fmt.Println("Enter the site name:")
	fmt.Scanln(&site)
	fmt.Println("Enter the password")
	fmt.Scanln(&password)

	newEntry.Site = sanitiseString(site)
	newEntry.Password = encryptPassword(password)
	return newEntry
}

func newPasswordAuto() PasswordEntry {
	var site string
	var password string
	var newEntry PasswordEntry

	fmt.Println("Enter the site name:")
	fmt.Scanln(&site)
	password = generatePassword()
	fmt.Println("The new password is " + password)

	newEntry.Site = sanitiseString(site)
	newEntry.Password = encryptPassword(password)
	return newEntry
}

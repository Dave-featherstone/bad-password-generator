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

const encryptionKey = "FruitFly2020"

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
	for i := 0; i < len(passwordList.PasswordList); i++ {
		fmt.Println("Site: " + passwordList.PasswordList[i].Site)
		fmt.Println("Password: " + decryptPassword(passwordList.PasswordList[i].Password))
		fmt.Println("")
	}
}

func createDefaultPasswordList() PasswordList {
	var passwordList PasswordList
	var passwordEntry PasswordEntry

	// populate this with my AWS secret first

	passwordEntry.Site = sanitiseString("AWS Token")
	passwordEntry.Password = encryptPassword("AKIAAGHO14951GHOGA91")
	passwordList.PasswordList = append(passwordList.PasswordList, passwordEntry)

	passwordEntry.Site = sanitiseString("AWS Secret")
	passwordEntry.Password = encryptPassword("NOWlgsn22+HoFAFglGAGAaGAGg29219goib1+GAW")
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

func sanitiseString(targetString string) string {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(targetString)
	return html
}

func savePasswordFile(passwordList PasswordList) {

	file, _ := json.MarshalIndent(passwordList, "", " ")

	_ = ioutil.WriteFile("passwordstore.json", file, 0644)

}

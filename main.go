package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
)

func main() {
	fmt.Print("Enter Chromium version: ")
	var chromiumVersion string
	fmt.Scanln(&chromiumVersion)

	fmt.Print("Enter Extension ID: ")
	var extensionID string
	fmt.Scanln(&extensionID)

	url := fmt.Sprintf("https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion=%s&x=id%%3D%s%%26installsource%%3Dondemand%%26uc", chromiumVersion, extensionID)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching extension:", err)
		return
	}
	defer response.Body.Close()

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		return
	}
	homeDir := usr.HomeDir

	downloadsFolder := fmt.Sprintf("%s/Downloads", homeDir)
	if _, err := os.Stat(downloadsFolder); os.IsNotExist(err) {
		os.Mkdir(downloadsFolder, 0755)
	}

	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]

	outputPath := fmt.Sprintf("%s/%s", downloadsFolder, filename)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	fmt.Println("Extension downloaded and saved to:", outputPath)
}


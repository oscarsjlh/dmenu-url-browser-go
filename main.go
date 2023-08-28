package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// List of bookmarks
	type Links struct {
		url  string
		name string
	}
	urls := []Links{
		{
			"https://www.example.com",
			"example",
		},
		{
			"https://www.github.com",
			"github",
		},

		// "https://www.google.com",
		// "https://www.github.com",
		// "https://console.cloud.google.com/artifacts?referrer=search&authuser=2&project=infrastructure-20220921-363208",
		// // Add more bookmarks here
	}

	// Create the rofi options
	rofiOptions := []string{
		"-dmenu",
		"-p", "Select Bookmark:",
	}
	var links []string
	// Command to run rofi
	rofiCmd := exec.Command("rofi", rofiOptions...)
	for _, link := range urls {
		links = append(links, link.name)
	}
	rofiCmd.Stdin = strings.NewReader(strings.Join(links, "\n"))
	rofiOutput, err := rofiCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Run rofi and get the selected bookmark

	selectedBookmark := strings.TrimSpace(string(rofiOutput))

	linkpos := findStringPosition(links, selectedBookmark)
	// Launch the browser with the selected bookmark
	browserurl := urls[linkpos]
	browserCmd := exec.Command("xdg-open", browserurl.url)
	browserCmd.Stdin = os.Stdin
	browserCmd.Stdout = os.Stdout
	browserCmd.Stderr = os.Stderr

	// Run the browser command
	err = browserCmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func findStringPosition(arr []string, target string) int {
	for i, s := range arr {
		if s == target {
			return i
		}
	}
	return -1 // Return -1 if the target string is not found
}

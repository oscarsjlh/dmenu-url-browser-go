package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Links struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

const (
	configFile     = "config.yaml"
	rofiPrompt     = "Select Bookmark:"
	rofiExecutable = "rofi"
	browserCmd     = "xdg-open"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	cFile := home + "/.config/browser/" + configFile
	urlInfo, err := parseYaml(cFile)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	selectedBookmark, err := runRofiAndGetSelection(urlInfo)
	if err != nil {
		log.Fatalf("Error running Rofi: %v", err)
	}

	linkIndex, err := findStringPosition(urlInfo, selectedBookmark)
	if err != nil {
		log.Fatalf("Error finding link: %v", err)
	}

	if err := launchBrowser(urlInfo[linkIndex].Url); err != nil {
		log.Fatalf("Error launching browser: %v", err)
	}
}

func (l Links) findStringPosition(target string) bool {
	return l.Name == target
}

func findStringPosition(urlInfo []Links, target string) (int, error) {
	for i, link := range urlInfo {
		if link.findStringPosition(target) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("link not found") // Return -1 if the target string is not found
}

func parseYaml(file string) ([]Links, error) {
	var urlInfo []Links
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &urlInfo); err != nil {
		return nil, err
	}
	return urlInfo, nil
}

func runRofiAndGetSelection(urlInfo []Links) (string, error) {
	var links []string
	for _, link := range urlInfo {
		links = append(links, link.Name)
	}
	rofiOptions := []string{
		"-dmenu",
		"-p", "Select Bookmark:",
	}
	rofiCmd := exec.Command("rofi", rofiOptions...)
	rofiCmd.Stdin = strings.NewReader(strings.Join(links, "\n"))
	rofiOutput, err := rofiCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return strings.TrimSpace(string(rofiOutput)), nil
}

func launchBrowser(url string) error {
	browserCmd := exec.Command(browserCmd, url)
	browserCmd.Stdin = os.Stdin
	browserCmd.Stdout = os.Stdout
	browserCmd.Stderr = os.Stderr

	return browserCmd.Run()
}

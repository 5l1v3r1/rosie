package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// readlines - Read lines of a text file into a slice
func readlines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	words := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

// getRandomAdjective - Get a random noun, not cryptographically secure
func getRandomAdjective() string {
	rand.Seed(time.Now().UnixNano())
	rosieDir := GetRosieDir()
	words := readlines(path.Join(rosieDir, "adjectives.txt"))
	word := words[rand.Intn(len(words)-1)]
	return strings.TrimSpace(word)
}

// getRandomNoun - Get a random noun, not cryptographically secure
func getRandomNoun() string {
	rand.Seed(time.Now().Unix())
	rosieDir := GetRosieDir()
	words := readlines(path.Join(rosieDir, "nouns.txt"))
	word := words[rand.Intn(len(words)-1)]
	return strings.TrimSpace(word)
}

// GetCodename - Returns a randomly generated 'codename'
func GetCodename() string {
	adjective := strings.ToUpper(getRandomAdjective())
	noun := strings.ToUpper(getRandomNoun())
	return fmt.Sprintf("%s_%s", adjective, noun)
}

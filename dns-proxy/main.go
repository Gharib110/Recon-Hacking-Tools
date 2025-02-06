package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parse(filename string) (map[string]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	records := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("invalid line: %s", line)
		}
		records[parts[0]] = parts[1]
	}
	log.Println("records set to: ")
	for k, v := range records {
		log.Printf("%s -> %s\n", k, v)
	}

	return records, scanner.Err()
}

func main() {

}

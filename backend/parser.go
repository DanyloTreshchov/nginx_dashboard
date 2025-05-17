package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var nginxLogPattern = regexp.MustCompile(`^(?P<IP>\S+) \S+ \S+ \[(?P<Time>[^\]]+)\] "(?:(?P<Method>\S+)(?: (?P<Path>\S+))?(?: (?P<HttpV>\S+))?)?" (?P<Status>\d{3}) (?P<Bytes>\d+) "(?P<Referrer>[^"]*)" "(?P<UserAgent>[^"]*)"`)

type LogEntry struct {
	IP        string
	Time      string
	Method    string
	Path      string
	HTTPV     string
	Status    string
	Bytes     string
	Referrer  string
	UserAgent string
	Flag	  string
}

func parseLine(line string) (*LogEntry, error) {
	matches := nginxLogPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("no match")
	}

	result := make(map[string]string)
	for i, name := range nginxLogPattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	entry := &LogEntry{
		IP:        result["IP"],
		Time:      result["Time"],
		Method:    result["Method"],
		Path:      result["Path"],
		Status:    result["Status"],
		HTTPV:	   result["HttpV"],
		Bytes:     result["Bytes"],
		Referrer:  result["Referrer"],
		UserAgent: result["UserAgent"],
	}

	return entry, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("../temp/access.log")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseLine(line)
		if err != nil {
			fmt.Printf("Error parsing line: %s\n", line)
			continue
		}

		fmt.Printf("IP: %s | Time: %s | Method: %s | Path: %s | HTTPV: %s | Status: %s | Bytes: %s | Flag: %s\n",
			entry.IP, entry.Time, entry.Method, entry.Path, entry.HTTPV, entry.Status, entry.Bytes, entry.Flag)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
	}
}


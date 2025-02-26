package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type RealHostsFile struct {
	mu sync.Mutex
}

func (r *RealHostsFile) ReadLines() ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open("/etc/hosts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (r *RealHostsFile) WriteLines(lines []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile("/etc/hosts", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (r *RealHostsFile) BlockedDomains() map[string]bool {
	lines, err := r.ReadLines()
	if err != nil {
		return make(map[string]bool)
	}

	blocked := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "127.0.0.2") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				blocked[fields[1]] = true
			}
		}
	}
	return blocked
}

func (r *RealHostsFile) Write(domains []string) error {
	lines, err := r.ReadLines()
	if err != nil {
		return err
	}

	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain] = true
	}

	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "127.0.0.2") {
			fields := strings.Fields(trimmed)
			if len(fields) >= 2 && domainSet[fields[1]] {
				continue
			}
		}
		newLines = append(newLines, line)
	}

	for _, domain := range domains {
		newLines = append(newLines, fmt.Sprintf("127.0.0.2\t%s", domain))
	}

	return r.WriteLines(newLines)
}

func verifyDNS(domain string) bool {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if ip.Equal(net.IPv4(127, 0, 0, 2)) {
			return true
		}
	}
	return false
}

func verifyAndReBlock(handler *RealHostsFile) {
	blocked := []string{"youtube.com", "www.youtube.com"}
	var toReBlock []string

	for _, domain := range blocked {
		if !verifyDNS(domain) {
			log.Printf("Domain %s not properly blocked", domain)
			toReBlock = append(toReBlock, domain)
		}
	}

	if len(toReBlock) > 0 {
		log.Printf("Re-blocking %d domains at startup", len(toReBlock))
		if err := handler.Write(toReBlock); err != nil {
			log.Printf("Error re-blocking: %v", err)
		}
	}
}

func watchWithPolling(handler *RealHostsFile, interval time.Duration) {
	var lastModTime time.Time

	for {
		fileInfo, err := os.Stat("/etc/hosts")
		if err != nil {
			log.Printf("Error stating /etc/hosts: %v", err)
			time.Sleep(interval)
			continue
		}

		if fileInfo.ModTime() != lastModTime {
			log.Println("File modification detected via polling")
			verifyAndReBlock(handler)
			lastModTime = fileInfo.ModTime()
		}

		time.Sleep(interval)
	}
}

func main() {
	if os.Geteuid() != 0 {
		log.Fatal("This program must be run as root")
	}

	handler := &RealHostsFile{}

	// Verify and re-block at startup
	verifyAndReBlock(handler)

	// Start polling watcher
	go watchWithPolling(handler, 5*time.Second)

	log.Println("Domain block daemon started (polling mode)")
	defer log.Println("Domain block daemon stopped")

	// Keep the main goroutine alive
	select {}
}

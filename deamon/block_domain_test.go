package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"strings"
)

// TestHostsFile handler for in-memory testing
type TestHostsFile struct {
	buffer *bytes.Buffer
}

func (t *TestHostsFile) Read() ([]string, error) {
	var domains []string
	scanner := bufio.NewScanner(t.buffer)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "127.0.0.1") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				domains = append(domains, parts[1])
			}
		}
	}
	return domains, scanner.Err()
}

func (t *TestHostsFile) Write(domains []string) error {
	for _, domain := range domains {
		fmt.Fprintf(t.buffer, "127.0.0.1 %s\n", domain)
	}
	return nil
}

func (t *TestHostsFile) BlockedDomains() map[string]bool {
	result := make(map[string]bool)
	scanner := bufio.NewScanner(t.buffer)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "127.0.0.1") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				result[parts[1]] = true
			}
		}
	}
	return result
}

type featureTest struct {
	blockList []string
	handler   HostsFileHandler
}

func (f *featureTest) theFollowingDomainListToBlock(table *godog.Table) error {
	for _, row := range table.Rows {
		f.blockList = append(f.blockList, row.Cells[0].Value)
	}
	return nil
}

func (f *featureTest) modifyHostsFile() error {
	return f.handler.Write(f.blockList)
}

func (f *featureTest) verifyDomainUnreachable(domain string) error {
	if !f.handler.BlockedDomains()[domain] {
		return fmt.Errorf("domain %s is not blocked", domain)
	}
	return nil
}

func (f *featureTest) iWantToBlockDomain(domain string) error {
	f.blockList = append(f.blockList, domain)
	return nil
}

func (f *featureTest) domainIsNotBlocked(domain string) error {
	if f.handler.BlockedDomains()[domain] {
		return fmt.Errorf("domain %s is already blocked", domain)
	}
	return nil
}

func (f *featureTest) verifyDomainsUnreachable(table *godog.Table) error {
	for _, row := range table.Rows {
		domain := row.Cells[0].Value
		if !f.handler.BlockedDomains()[domain] {
			return fmt.Errorf("domain %s is not blocked", domain)
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ft := &featureTest{
		handler: &TestHostsFile{buffer: new(bytes.Buffer)},
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ft.blockList = nil
		ft.handler = &TestHostsFile{buffer: new(bytes.Buffer)}
		return ctx, nil
	})

	ctx.Step(`^the list of block domain is:$`, ft.theFollowingDomainListToBlock)
	ctx.Step(`^the following domain list of domain to be block:$`, ft.theFollowingDomainListToBlock)
	ctx.Step(`^I block the list of domain$`, ft.modifyHostsFile)
	ctx.Step(`^the hosts file is modify$`, ft.modifyHostsFile)
	ctx.Step(`^"([^"]*)" can't be access anymore$`, ft.verifyDomainUnreachable)
	ctx.Step(`^verify that each domain is unreachable$`, ft.verifyDomainsUnreachable)
	ctx.Step(`^I want to block the domain "([^"]*)"$`, ft.iWantToBlockDomain)
	ctx.Step(`^"([^"]*)" is not already block$`, ft.domainIsNotBlocked)
	ctx.Step(`^the following domains can't be reached:$`, ft.verifyDomainsUnreachable)
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Gem struct {
	Remote       string
	IsLocal      bool
	IsRubyGems   bool
	IsTransitive bool
	Name         string
	Version      string
}

type RubyGemsResponse struct {
	Name      string `json:"name"`
	Downloads int64  `json:"downloads"`
	Version   string `json:"version"`
}

// RubyGemsLookup represents a collection of rubygems packages to be tested for dependency confusion.
type RubyGemsLookup struct {
	Packages []Gem
	Verbose  bool
}

// NewRubyGemsLookup constructs an `RubyGemsLookup` struct and returns it.
func NewRubyGemsLookup(verbose bool) PackageResolver {
	return &RubyGemsLookup{Packages: []Gem{}, Verbose: verbose}
}

// ReadPackagesFromFile reads package information from a Gemfile.lock file
//
// Returns any errors encountered
func (r *RubyGemsLookup) ReadPackagesFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var remote string
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "remote:") {
			remote = strings.TrimSpace(strings.SplitN(trimmedLine, ":", 2)[1])
		} else if trimmedLine == "revision:" {
			continue
		} else if trimmedLine == "branch:" {
			continue
		} else if trimmedLine == "GIT" {
			continue
		} else if trimmedLine == "GEM" {
			continue
		} else if trimmedLine == "PATH" {
			continue
		} else if trimmedLine == "PLATFORMS" {
			break
		} else if trimmedLine == "specs:" {
			continue
		} else if len(trimmedLine) > 0 {
			parts := strings.SplitN(trimmedLine, " ", 2)
			name := strings.TrimSpace(parts[0])
			var version string
			if len(parts) > 1 {
				version = strings.TrimRight(strings.TrimLeft(parts[1], "("), ")")
			} else {
				version = ""
			}
			r.Packages = append(r.Packages, Gem{
				Remote:       remote,
				IsLocal:      !strings.HasPrefix(remote, "http"),
				IsRubyGems:   strings.HasPrefix(remote, "https://rubygems.org"),
				IsTransitive: countLeadingSpaces(line) == 6,
				Name:         name,
				Version:      version,
			})
		} else {
			continue
		}
	}
	return nil
}

// PackagesNotInPublic determines if a rubygems package does not exist in the public rubygems package repository.
//
// Returns a slice of strings with any rubygem packages not in the public rubygems package repository
func (r *RubyGemsLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range r.Packages {
		if pkg.IsLocal || !pkg.IsRubyGems {
			continue
		}
		if !r.isAvailableInPublic(pkg.Name, 0) {
			notavail = append(notavail, pkg.Name)
		}
	}
	return notavail
}

// isAvailableInPublic determines if a rubygems package exists in the public rubygems.org package repository.
//
// Returns true if the package exists in the public rubygems package repository.
func (r *RubyGemsLookup) isAvailableInPublic(pkgname string, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package: %s\n", pkgname)
		return false
	}
	url := fmt.Sprintf("https://rubygems.org/api/v1/gems/%s.json", pkgname)
	if r.Verbose {
		fmt.Printf("Checking: %s : \n", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf(" [W] Error when trying to request %s: %s\n", url, err)
		return false
	}
	defer resp.Body.Close()
	if r.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		rubygemsResp := RubyGemsResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &rubygemsResp)
		if err != nil {
			fmt.Printf(" [W] Error when trying to unmarshal response from %s: %s\n", url, err)
			return false
		}
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		return r.isAvailableInPublic(pkgname, retry)
	}
	return false
}

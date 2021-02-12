package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ComposerJSON struct {
	Require map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

type ComposerLookup struct {
	Packages []string
	Verbose bool
}

func NewComposerLookup(verbose bool) PackageResolver {
	return &ComposerLookup{Packages: []string{}, Verbose: verbose}
}

func (c *ComposerLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data := ComposerJSON{}
	err = json.Unmarshal([]byte(rawfile), &data)
	if err != nil {
		return err
	}

	for pkgname := range data.Require {
		c.Packages = append(c.Packages, pkgname)
	}

	for pkgname := range data.RequireDev {
		c.Packages = append(c.Packages, pkgname)
	}

	return nil
}

func (c *ComposerLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range c.Packages {
		if pkg == "php" {
			continue
		}

		if !c.isAvailableInPublic(pkg, 0) {
			notavail = append(notavail, pkg)
		}
	}

	return notavail
}

func (c *ComposerLookup) isAvailableInPublic(pkgname string, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package %s\n", pkgname)

		return false
	}

	if c.Verbose {
		fmt.Printf("Checking: https://packagist.org/packages/%s : ", pkgname)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("https://packagist.org/packages/" + pkgname)
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://packagist.org/packages/%s : %s\n", pkgname, err)

		return false
	}
	defer resp.Body.Close()

	if c.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}

	if resp.StatusCode == http.StatusOK {
		return true
	}

	if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying..\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1

		c.isAvailableInPublic(pkgname, retry)
	}

	return false
}

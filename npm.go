package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// PackageJSON represents the dependencies of an npm package
type PackageJSON struct {
	Dependencies         map[string]string `json:"dependencies,omitempty"`
	DevDependencies      map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
	BundledDependencies  []string          `json:"bundledDependencies,omitempty"`
	BundleDependencies   []string          `json:"bundleDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
}

type NpmResponse struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Time struct {
		Unpublished NpmResponseUnpublished `json:"unpublished"`
	} `json:"time"`
}

type NpmResponseUnpublished struct {
		Maintainers []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"maintainers"`
		Name string `json:"name"`
		Tags struct {
			Latest string `json:"latest"`
		} `json:"tags"`
		Time     time.Time `json:"time"`
		Versions []string  `json:"versions"`
}

// NotAvailable returns true if the package has its all versions unpublished making it susceptible for takeover
func (n *NpmResponse) NotAvailable() bool {
	// Check if a known field has a value
	return len(n.Time.Unpublished.Name) > 0
}

// NPMLookup represents a collection of npm packages to be tested for dependency confusion.
type NPMLookup struct {
	Packages []string
	Verbose  bool
}

// NewNPMLookup constructs an `NPMLookup` struct and returns it.
func NewNPMLookup(verbose bool) PackageResolver {
	return &NPMLookup{Packages: []string{}, Verbose: verbose}
}

// ReadPackagesFromFile reads package information from an npm package.json file
//
// Returns any errors encountered
func (n *NPMLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := PackageJSON{}
	err = json.Unmarshal([]byte(rawfile), &data)
	if err != nil {
		fmt.Printf(" [W] Non-fatal issue encountered while reading %s : %s\n", filename, err)
	}
	for pkgname := range data.Dependencies {
		n.Packages = append(n.Packages, pkgname)
	}
	for pkgname := range data.DevDependencies {
		n.Packages = append(n.Packages, pkgname)
	}
	for pkgname := range data.PeerDependencies {
		n.Packages = append(n.Packages, pkgname)
	}
	for pkgname := range data.OptionalDependencies {
		n.Packages = append(n.Packages, pkgname)
	}
	n.Packages = append(n.Packages, data.BundledDependencies...)
	n.Packages = append(n.Packages, data.BundleDependencies...)
	return nil
}

// PackagesNotInPublic determines if an npm package does not exist in the public npm package repository.
//
// Returns a slice of strings with any npm packages not in the public npm package repository
func (n *NPMLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range n.Packages {
		if !n.isAvailableInPublic(pkg, 0) {
			notavail = append(notavail, pkg)
		}
	}
	return notavail
}

// isAvailableInPublic determines if an npm package exists in the public npm package repository.
//
// Returns true if the package exists in the public npm package repository.
func (n *NPMLookup) isAvailableInPublic(pkgname string, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package: %s\n", pkgname)
		return false
	}
	if n.Verbose {
		fmt.Print("Checking: https://registry.npmjs.org/" + pkgname + "/ : ")
	}
	resp, err := http.Get("https://registry.npmjs.org/" + pkgname + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://registry.npmjs.org/"+pkgname+"/ : %s\n", err)
		return false
	}
	defer resp.Body.Close()
	if n.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		npmResp := NpmResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &npmResp)
		if npmResp.NotAvailable() {
			if n.Verbose {
				fmt.Printf("[W] Package %s was found, but all its versions are unpublished, making anyone able to takeover the namespace.\n", pkgname)
			}
			return false
		}
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		n.isAvailableInPublic(pkgname, retry)
	}
	return false
}

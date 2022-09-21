package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// NPMLookup represents a collection of npm packages to be tested for dependency confusion.
type MVNLookup struct {
	Packages []MVNPackage
	Verbose  bool
}

type MVNPackage struct {
	Group    string
	Artifact string
	Version  string
}

// NewNPMLookup constructs an `MVNLookup` struct and returns it.
func NewMVNLookup(verbose bool) PackageResolver {
	return &MVNLookup{Packages: []MVNPackage{}, Verbose: verbose}
}

// ReadPackagesFromFile reads package information from an npm package.json file
//
// Returns any errors encountered
func (n *MVNLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	fmt.Print("Checking: filename: " + filename + "\n")

	var project MavenProject
	if err := xml.Unmarshal([]byte(rawfile), &project); err != nil {
		log.Fatalf("unable to unmarshal pom file. Reason: %s\n", err)
	}

	for _, dep := range project.Dependencies {
		n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
	}

	for _, dep := range project.Build.Plugins {
		n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
	}

	for _, build := range project.Profiles {
		for _, dep := range build.Build.Plugins {
			n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
		}
	}

	return nil
}

// PackagesNotInPublic determines if an npm package does not exist in the public npm package repository.
//
// Returns a slice of strings with any npm packages not in the public npm package repository
func (n *MVNLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range n.Packages {
		if !n.isAvailableInPublic(pkg, 0) {
			notavail = append(notavail, pkg.Group+"/"+pkg.Artifact)
		}
	}
	return notavail
}

// isAvailableInPublic determines if an npm package exists in the public npm package repository.
//
// Returns true if the package exists in the public npm package repository.
func (n *MVNLookup) isAvailableInPublic(pkg MVNPackage, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package: %s\n", pkg.Group)
		return false
	}
	if pkg.Group == "" {
		return true
	}

	group := strings.Replace(pkg.Group, ".", "/", -1)
	if n.Verbose {
		fmt.Print("Checking: https://repo1.maven.org/maven2/" + group + "/ ")
	}
	resp, err := http.Get("https://repo1.maven.org/maven2/" + group + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://repo1.maven.org/maven2/"+group+"/ : %s\n", err)
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
				fmt.Printf("[W] Package %s was found, but all its versions are unpublished, making anyone able to takeover the namespace.\n", pkg.Group)
			}
			return false
		}
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		n.isAvailableInPublic(pkg, retry)
	}
	return false
}

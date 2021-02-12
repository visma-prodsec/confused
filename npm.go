package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PackageJSON struct {
	Dependencies map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	BundledDependencies []string `json:"bundledDependencies"`
	BundleDependencies []string `json:"bundleDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

type NPMLookup struct{
	Packages []string
	Verbose bool
}

func NewNPMLookup(verbose bool) PackageResolver {
	return &NPMLookup{Packages: []string{}, Verbose: verbose}
}

func (n *NPMLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := PackageJSON{}
	err = json.Unmarshal([]byte(rawfile), &data)
	if err != nil {
		return err
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

func (n *NPMLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range n.Packages {
		if !n.isAvailableInPublic(pkg, 0) {
			notavail = append(notavail, pkg)
		}
	}
	return notavail
}

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
		fmt.Printf(" [W] Error when trying to request https://registry.npmjs.org/" + pkgname + "/ : %s\n", err)
		return false
	}
	if n.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		n.isAvailableInPublic(pkgname, retry)
	}
	return false
}
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PackageJSON struct {
	Dependencies map[string]string `json:"dependencies"`
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
	for pkgname, _ := range data.Dependencies {
		n.Packages = append(n.Packages, pkgname)
	}
	return nil
}

func (n *NPMLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range n.Packages {
		if !n.isAvailableInPublic(pkg) {
			notavail = append(notavail, pkg)
		}
	}
	return notavail
}

func (n *NPMLookup) isAvailableInPublic(pkgname string) bool {
	if n.Verbose {
		fmt.Print("Checking: https://www.npmjs.com/package/" + pkgname + " : ")
	}
	resp, _ := http.Get("https://www.npmjs.com/package/" + pkgname)
	if n.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
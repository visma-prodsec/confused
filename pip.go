package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PythonLookup represents a collection of python packages to be tested for dependency confusion.
type PythonLookup struct {
	Packages []string
	Verbose  bool
}

// NewPythonLookup constructs a `PythonLookup` struct and returns it
func NewPythonLookup(verbose bool) PackageResolver {
	return &PythonLookup{Packages: []string{}, Verbose: verbose}
}

// ReadPackagesFromFile reads package information from a python `requirements.txt` file
//
// Returns any errors encountered
func (p *PythonLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	line := ""
	for _, l := range strings.Split(string(rawfile), "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			continue
		}
		if len(l) > 0 {
			// Support line continuation
			if strings.HasSuffix(l, "\\") {
				line += l[:len(l) - 1]
				continue
			}
			line += l
			pkgrow := strings.FieldsFunc(line, p.pipSplit)
			if len(pkgrow) > 0 {
				p.Packages = append(p.Packages, strings.TrimSpace(pkgrow[0]))
			}
			// reset the line variable
			line = ""
		}
	}
	return nil
}

// PackagesNotInPublic determines if a python package does not exist in the pypi package repository.
//
// Returns a slice of strings with any python packages not in the pypi package repository
func (p *PythonLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range p.Packages {
		if !p.isAvailableInPublic(pkg) {
			notavail = append(notavail, pkg)
		}
	}
	return notavail
}

func (p *PythonLookup) pipSplit(r rune) bool {
	delims := []rune{
		'=',
		'<',
		'>',
		'!',
		' ',
		'#',
		'[',
	}
	return inSlice(r, delims)
}

// isAvailableInPublic determines if a python package exists in the pypi package repository.
//
// Returns true if the package exists in the pypi package repository.
func (p *PythonLookup) isAvailableInPublic(pkgname string) bool {
	if p.Verbose {
		fmt.Print("Checking: https://pypi.org/project/" + pkgname + "/ : ")
	}
	resp, err := http.Get("https://pypi.org/project/" + pkgname + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://pypi.org/project/"+pkgname+"/ : %s\n", err)
		return false
	}
	if p.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

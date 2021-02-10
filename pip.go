package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PythonLookup struct{
	Packages []string
	Verbose bool
}

func NewPythonLookup(verbose bool) PackageResolver {
	return &PythonLookup{Packages: []string{}, Verbose: verbose}
}

func (p *PythonLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	for _, l := range strings.Split(string(rawfile), "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			continue
		}
		if len(l) > 0 {
			pkgrow := strings.FieldsFunc(l, p.pipSplit)
			if len(pkgrow) > 0 {
				p.Packages = append(p.Packages, strings.TrimSpace(pkgrow[0]))
			}
		}
	}
	return nil
}

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
	}
	return inSlice(r, delims)
}

func (p *PythonLookup) isAvailableInPublic(pkgname string) bool {
	if p.Verbose {
		fmt.Print("Checking: https://pypi.org/project/" + pkgname + "/ : ")
	}
	resp, err := http.Get("https://pypi.org/project/" + pkgname + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://pypi.org/project/" + pkgname + "/ : %s\n", err)
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
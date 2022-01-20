package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"net/http"
	"strings"
)

// PythonLookup represents a collection of python packages to be tested for dependency confusion.
type PythonLookup struct {
	Packages       []string
	Verbose        bool
	PackageManager string
}

// NewPythonLookup constructs a `PythonLookup` struct and returns it
func NewPythonLookup(verbose bool, packageManager string) PackageResolver {
	return &PythonLookup{Packages: []string{}, Verbose: verbose, PackageManager: packageManager}
}

// ReadPackagesFromFile chooses a file parser based on the user-supplied python package manager.
//
// Returns any errors encountered
func (p *PythonLookup) ReadPackagesFromFile(filename string) error {
	switch p.PackageManager {
	case "pip":
		return p.ReadPackagesFromRequirementsTxt(filename)
	case "pipenv":
		return p.ReadPackagesFromPipfile(filename)
	default:
		return fmt.Errorf("Python package manager not implemented: %s", p.PackageManager)
	}
}

// ReadPackagesFromRequirementsTxt reads package information from a python `requirements.txt`.
//
// Returns any errors encountered
func (p *PythonLookup) ReadPackagesFromRequirementsTxt(filename string) error {
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

// ReadPackagesFromPipfile reads package information from a python `Pipfile`.
//
// Returns any errors encountered
func (p *PythonLookup) ReadPackagesFromPipfile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	config, err := toml.Load(string(rawfile))
	if err != nil {
		return err
	}
	packages := config.Get("packages")
	if packages != nil {
		p.Packages = append(p.Packages, packages.(*toml.Tree).Keys()...)
	}
	dev_packages := config.Get("dev-packages")
	if dev_packages != nil {
		p.Packages = append(p.Packages, dev_packages.(*toml.Tree).Keys()...)
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
		'~',
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

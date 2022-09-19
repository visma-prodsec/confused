/*
Package main implements an automated Dependency Confusion scanner.

Original research provided by Alex Birsan.

Original blog post detailing Dependency Confusion : https://medium.com/@alex.birsan/dependency-confusion-4a5d60fec610 .
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var resolver PackageResolver
	lang := ""
	verbose := false
	filename := ""
	safespaces := ""
	flag.StringVar(&lang, "l", "npm", "Package repository system. Possible values: \"pip\", \"npm\", \"composer\", \"mvn\", \"rubygems\"")
	flag.StringVar(&safespaces, "s", "", "Comma-separated list of known-secure namespaces. Supports wildcards")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	// Check that we have a filename
	if flag.NArg() == 0 {
		Help()
		flag.Usage()
		os.Exit(1)
	}

	filename = flag.Args()[0]

	switch lang {
	case "pip":
		resolver = NewPythonLookup(verbose)
	case "npm":
		resolver = NewNPMLookup(verbose)
	case "composer":
		resolver = NewComposerLookup(verbose)
	case "mvn":
		resolver = NewMVNLookup(verbose)
	case "rubygems":
		resolver = NewRubyGemsLookup(verbose)
	default:
		fmt.Printf("Unknown package repository system: %s\n", lang)
		os.Exit(1)
	}

	err := resolver.ReadPackagesFromFile(filename)
	if err != nil {
		fmt.Printf("Encountered an error while trying to read packages from file: %s\n", err)
		os.Exit(1)
	}
	outputPackages := removeSafe(resolver.PackagesNotInPublic(), safespaces)
	PrintResult(outputPackages)
}

// Help outputs tool usage and help
func Help() {
	fmt.Printf("Usage:\n %s [-l LANGUAGENAME] depfilename.ext\n", os.Args[0])
}

// PrintResult outputs the result of the scanner
func PrintResult(notavail []string) {
	if len(notavail) == 0 {
		fmt.Printf("[*] All packages seem to be available in the public repositories. \n\n" +
			"In case your application uses private repositories please make sure that those namespaces in \n" +
			"public repositories are controlled by a trusted party.\n\n")
		return
	}
	fmt.Printf("Issues found, the following packages are not available in public package repositories:\n")
	for _, n := range notavail {
		fmt.Printf(" [!] %s\n", n)
	}
	os.Exit(1)
}

// removeSafe removes known-safe package names from the slice
func removeSafe(packages []string, safespaces string) []string {
	retSlice := []string{}
	safeNamespaces := []string{}
	var ignored bool
	safeTmp := strings.Split(safespaces, ",")
	for _, s := range safeTmp {
		safeNamespaces = append(safeNamespaces, strings.TrimSpace(s))
	}
	for _, p := range packages {
		ignored = false
		for _, s := range safeNamespaces {
			ok, err := filepath.Match(s, p)
			if err != nil {
				fmt.Printf(" [W] Encountered an error while trying to match a known-safe namespace %s : %s\n", s, err)
				continue
			}
			if ok {
				ignored = true
			}
		}
		if !ignored {
			retSlice = append(retSlice, p)
		}
	}
	return retSlice
}

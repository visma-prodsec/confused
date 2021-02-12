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
)

func main() {
	var resolver PackageResolver
	lang := ""
	verbose := false
	filename := ""
	flag.StringVar(&lang, "l", "npm", "Package repository system. Possible values: \"pip\", \"npm\"")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	// Check that we have a filename
	if flag.NArg() == 0 {
		Help()
		flag.Usage()
		os.Exit(1)
	}

	filename = flag.Args()[0]
	if lang == "pip" {
		resolver = NewPythonLookup(verbose)
	} else if lang == "npm" {
		resolver = NewNPMLookup(verbose)
	} else {
		fmt.Printf("Unknown package repository system: %s\n", lang)
		os.Exit(1)
	}
	err := resolver.ReadPackagesFromFile(filename)
	if err != nil {
		fmt.Printf("Encountered an error while trying to read packages from file: %s\n", err)
		os.Exit(1)
	}
	PrintResult(resolver.PackagesNotInPublic())
}

// Help outputs tool usage and help
func Help() {
	fmt.Println(fmt.Sprintf(`Usage:
 %s [-l LANGUAGENAME] depfilename.ext
`, os.Args[0]))
}

// PrintResult outputs the result of the scanner
func PrintResult(notavail []string) {
	if len(notavail) == 0 {
		fmt.Printf(" [*] All packages seem to be available in the public repositories. Dependency confusion should not be possible.\n")
		return
	}
	fmt.Printf("Issues found, the following packages are not available in public package repositories:\n")
	for _, n := range notavail {
		fmt.Printf(" [!] %s\n", n)
	}
}

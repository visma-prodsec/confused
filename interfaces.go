package main

// A PackageResolver resolves package information from a file
//
// ReadPackagesFromFile should take a filepath as input and parse relevant package information into a struct and then return any errors encountered while reading or unmarshalling the file.
//
// PackagesNotInPublic should determine whether or not a package is not available in a public package repository and return a slice of all packages not available in a public package repository.
type PackageResolver interface {
	ReadPackagesFromFile(string) error
	PackagesNotInPublic() []string
}

package main

type PackageResolver interface {
	ReadPackagesFromFile(string) error
	PackagesNotInPublic() []string
}
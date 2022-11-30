# Confused

A tool for checking for lingering free namespaces for private package names referenced in dependency configuration
for Python (pypi) `requirements.txt`, JavaScript (npm) `package.json`, PHP (composer) `composer.json`, MVN (maven) `pom.xml` or Ruby (RubyGems) `Gemfile`.

## What is this all about?

On 9th of February 2021, a security researcher Alex Birsan [published an article](https://medium.com/@alex.birsan/dependency-confusion-4a5d60fec610)
that touched different resolve order flaws in dependency management tools present in multiple programming language ecosystems.

Microsoft [released a whitepaper](https://azure.microsoft.com/en-gb/resources/3-ways-to-mitigate-risk-using-private-package-feeds/)
describing ways to mitigate the impact, while the root cause still remains.

## Interpreting the tool output

`confused` simply reads through a dependency definition file of an application and checks the public package repositories
for each dependency entry in that file. It will proceed to report all the package names that are not found in the public
repositories - a state that implies that a package might be vulnerable to this kind of attack, while this vector has not
yet been exploited.

This however doesn't mean that an application isn't already being actively exploited. If you know your software is using
private package repositories, you should ensure that the namespaces for your private packages have been claimed by a
trusted party (typically yourself or your company).

### Known false positives

Some packaging ecosystems like npm have a concept called "scopes" that can be either private or public. In short it means
a namespace that has an upper level - the scope. The scopes are not inherently visible publicly, which means that `confused`
cannot reliably detect if it has been claimed. If your application uses scoped package names, you should ensure that a
trusted party has claimed the scope name in the public repositories.

## Installation

- [Download](https://github.com/visma-prodsec/confused/releases/latest) a prebuilt binary from [releases page](https://github.com/visma-prodsec/confused/releases/latest), unpack and run!

  _or_
- If you have recent go compiler installed: `go get -u github.com/visma-prodsec/confused` (the same command works for updating)

  _or_
- git clone https://github.com/visma-prodsec/confused ; cd confused ; go get ; go build

## Usage
```
Usage:
 confused [-l LANGUAGENAME] depfilename.ext

Usage of confused:
  -l string
        Package repository system. Possible values: "pip", "npm", "composer", "mvn", "rubygems" (default "npm")
  -s string
        Comma-separated list of known-secure namespaces. Supports wildcards
  -v    Verbose output

```

## Example

### Python (PyPI)
```
./confused -l pip requirements.txt

Issues found, the following packages are not available in public package repositories:
 [!] internal_package1

```

### JavaScript (npm)
```
./confused -l npm package.json

Issues found, the following packages are not available in public package repositories:
 [!] internal_package1
 [!] @mycompany/internal_package1
 [!] @mycompany/internal_package2

# Example when @mycompany private scope has been registered in npm, using -s
./confused -l npm -s '@mycompany/*' package.json

Issues found, the following packages are not available in public package repositories:
 [!] internal_package1
```

### Maven (mvn)
```
./confused -l mvn pom.xml

Issues found, the following packages are not available in public package repositories:
 [!] internal
 [!] internal/package1
 [!] internal/_package2

```

### Ruby (rubygems)
```
./confused -l rubygems Gemfile.lock

Issues found, the following packages are not available in public package repositories:
 [!] internal
 [!] internal/package1
 [!] internal/_package2
 
```

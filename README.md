# Confused

A tool for checking for lingering free namespaces for private package names referenced in dependency configuration 
for Python (pypi) `requirements.txt` or JavaScript (npm) `package.json`


## Installation 

- [Download](https://github.com/visma-prodsec/confused/releases/latest) a prebuilt binary from [releases page](https://github.com/visma-prodsec/confused/releases/latest), unpack and run!

  _or_
- If you have recent go compiler installed: `go get -u github.com/visma-prodsec/confused` (the same command works for updating)

  _or_
- git clone https://github.com/visma-prodsec/confused ; cd confused ; go get ; go build

## Usage
```
Usage:
 ./confused [-l LANGUAGENAME] depfilename.ext

Usage of ./confused:
  -l string
        Package repository system. Possible values: "pip", "npm", "ruby" (default "pip")
  -v    Verbose output

```

## Example
```
./confused -l pip requirements.txt

Issues found, the following packages are not available in public package repositories:
 [!] internal_package1
```
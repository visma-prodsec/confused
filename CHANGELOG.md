# Changelog

- main
    - New
        - npm: In case package was found, also check if all the package versions have been unpublished. This makes the package vulnerable to takeover
        - npm: Check for http & https and GitHub version references
    - Changed
        - Fixed a bug where the pip requirements.txt parser processes a 'tilde equals' sign.

- v0.3
    - New
        - PHP (composer) support
        - Command line parameter to let the user to flag namespaces as known-safe
    - Changed
        - Python (pypi) dependency definition files that use line continuation are now parsed correctly
        - Revised the output to clarify the usage
        - Fixed npm package.json file parsing issues when the source file is not following the specification

- v0.2
    - Changed
        - npm registry checkup url
        - Throttle the rate of requests in case of 429 (Too many requests) responses

- v0.1
   - Initial release

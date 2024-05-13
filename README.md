# xc-strings

`xc-strings` is a command-line tool designed to help Swift developers manage and optimize their localization files. 
It provides functionalities to find unused localization keys, clean them from `.strings` files, and optionally sort the keys in these files for better organization.

## Features

- **Find Unused Keys**: Scans Swift files to detect any localization keys that are no longer used.
- **Find Duplicate Keys**: Scans `.strings` files to detect any duplicate keys within the same file.
- **Sort `.strings` Files**: Sorts keys in `.strings` files to maintain a consistent order.

## Installation

### Prerequisites

- Go 1.15 or later

### Building from Source

Clone the repository and build the executable:

```bash
git clone git@github.com:phillippbertram/go-xc-strings.git
cd go-xc-strings

# run directly
go run main.go help

# build and run the executable (macOS)
make build
./dist/go-xc-strings_darwin_arm64/xcs help
```

### Setup Development Environment

- Install [Go](https://golang.org/doc/install)
- Optional: Install golangci-lint: `brew install golangci-lint`
- Optional: Install goreleaser: `brew install goreleaser`

## Usage

```bash
# get help and list all available commands
xcs help

# list unused localization keys
# -b: path to the base localization file
# args: path to the directory containing the Swift files
# --strings: path to the directory containing the .strings files
# -i: optional glob pattern to exclude files (useful to ignore R.string generated files)
xcs unused -b path/to/Localizable.strings -d path/to/swift/files -i "*.generated.swift" App/Resources --remove

# sort strings files
xcs sort App/Resources

# find and remove specific keys from all strings files that are not used in the Swift files
xcs keys "this_is_a_key" "another_key" AppIOS/Resources --remove

# open github repository or release page
xcs gh [--releases]
```

## Configuration

No additional configuration is needed to run `xc-strings`.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests with any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, you can open an issue in the GitHub issue tracker.

## Authors

- **Phillipp Bertram** - *Initial work* - [phillippbertram](https://github.com/phillippbertram)
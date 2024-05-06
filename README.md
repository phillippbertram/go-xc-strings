# xc-strings

`xc-strings` is a command-line tool designed to help Swift developers manage and optimize their localization files. 
It provides functionalities to find unused localization keys, clean them from `.strings` files, and optionally sort the keys in these files for better organization.

## Features

- **Find Unused Keys**: Scans Swift files to detect any localization keys that are no longer used.
- **Clean `.strings` Files**: Removes unused keys from `.strings` files.
- **Sort `.strings` Files**: Optionally sorts keys in `.strings` files to maintain a consistent order.

## Installation

### Prerequisites

- Go 1.15 or later

### Building from Source

Clone the repository and build the executable:

```bash
git clone git@github.com:phillippbertram/go-xc-strings.git
cd go-xc-strings
make build
./bin/xc-strings help
```

## Usage

```bash

# get help and list all available commands
xc-strings help

# list unused localization keys
# -b: path to the base localization file
# -d: optional path to the directory containing the Swift files
# --strings: path to the directory containing the .strings files
# -i: optional glob pattern to exclude files
xc-strings unused -b path/to/Localizable.strings -d path/to/swift/files --strings App/Resources -i "*.generated.swift"

# sort strings files
xc-strings sort App/Resources

# remove keys from all strings files that are not used in the Swift files
xc-strings remove "this_is_a_key" AppIOS/Resources

# clean: aggregated command to remove unused keys and sort the strings files
xc-strings clean -b path/to/Localizable.strings -d path/to/swift/files --strings App/Resources -i "*.generated.swift"
```

## Configuration

No additional configuration is needed to run xc-strings.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests with any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, you can open an issue in the GitHub issue tracker.

## Authors

- **Phillipp Bertram** - *Initial work* - [phillippbertram](https://github.com/phillippbertram)
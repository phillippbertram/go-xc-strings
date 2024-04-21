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
go build -o xc-strings
```

## Usage

### Finding Unused Keys

To find and report unused localization keys:

```bash
./xc-strings find -r path/to/Localizable.strings -d path/to/swift/files
```

### Cleaning `.strings` Files

To remove unused keys from `.strings` files:

```bash
./xc-strings clean /path/to/project --sort
```

This command will remove unused keys and sort the `.strings` files if the `--sort` flag is provided.

### Dry-Run Mode

To simulate changes without making actual modifications:

```bash
./xc-strings clean /path/to/project --dry-run
```

### Sorting `.strings` Files

To sort `.strings` files without removing keys:

```bash
./xc-strings sort /path/to/strings/files
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
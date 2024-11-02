# My-ls-1


**My-ls-1** is a custom implementation of the Unix `ls` command written in Go. This project aims to replicate the functionality of the original `ls` command with certain variations and additional features.

## File Structure

```sh
my-ls-1/
├─ LICENSE
├─ README.md
├─ cmd/
│  ├─ ls/
│  │  └─ main.go
│  └─ terminal/
│     ├─ lsOptions/
│     │  ├─ read-dir.go
│     │  └─ read-dir_test.go
│     ├─ width.go
│     └─ width_test.go
├─ go.mod
├─ internal/
│  ├─ list/
│  │  ├─ list.go
│  │  └─ list_test.go
│  └─ sort/
│     ├─ sort.go
│     └─ sort_test.go
├─ pkg/
│  ├─ fileinfo/
│  │  ├─ fileinfo.go
│  │  └─ fileinfo_test.go
│  ├─ options/
│  │  ├─ options.go
│  │  └─ options_test.go
│  └─ utils/
│     ├─ color/
│     │  ├─ color.go
│     │  └─ color_test.go
│     ├─ determinePath.go
│     ├─ determinePath_test.go
│     ├─ formatfiles.go
│     ├─ formatfiles_test.go
│     ├─ math.go
│     ├─ math_test.go
│     ├─ printfiles.go
│     └─ printfiles_test.go
└─ run.sh
```

## Table of Contents

- [Objective](#objectives)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Test](#testing)
- [Contribution](#contribution)
- [Authors](#authors)
- [License](#license)
- [Acknowledgement](#acknowledgements)


## Objectives

The main objective of this project is to create a command that lists files and directories within a specified directory. By default, if no directory is specified, it lists files and directories in the current directory.

## Features

- Mimics the behavior of the original `ls` command.
- Supports the following flags:
  - `-l`: Long format listing
  - `-R`: Recursively list subdirectories
  - `-a`: Show all files, including hidden ones
  - `-r`: Reverse order of sorting
  - `-t`: Sort by modification time


## Installation

To use `my-ls-1`, you need to have Go installed on your machine. If you haven't installed Go yet, please follow the instructions on the [Go installation page](https://golang.org/doc/install).

1. Clone this repository:
    ```bash
    git clone https://learn.zone01kisumu.ke/git/somotto/my-ls-1.git

    cd my-ls-1
    ```

2. Build the project:
When you're at cmd/ls directory
    ```bash
    go build -o my-ls
    ```
    or run our script
    when you're at the root directory
    ```bash
    ./run.sh
    ```
   You can now run `my-ls` in your terminal.

## Usage

To list files in the current directory:
```bash
./my-ls
```
To list files in a specific directory
```bash
./my-ls /path/to/directory
```
To use the available flags
```bash
./my-ls -l

./my-ls -a

./my-ls -R

./my-ls -r

./my-ls -t

./my-ls -lR

./my-ls -S

./my-ls -1
```
## Testing
```bash
go test ./...

go test -cover ./...
```
## Contribution

If you find a bug or have suggestions for improvement, please submit an issue or a pull request. Contributions are welcome!

## Authors
- [Stephen Omotto](https://github.com/somotto)
- [Raymond Ogwel](https://github.com/anxielray)
- [Stephen Kisengese](https://learn.zone01kisumu.ke/git/skisenge)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgements
This project is inspired by the standard ls command in Unix/Linux systems. We recommend consulting the ls command manual(`man ls`) for more details on its options and usage.
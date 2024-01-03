![build](https://github.com/thejezzi/watchmego/actions/workflows/build.yml/badge.svg)

# Watch Me Go

Watchmego checks a directory for any changes in go files and runs a make
command.

## Installation

Make sure you added your GOPATH to PATH otherwise you won't be able to
run it globally.

```bash
go install github.com/thejezzi/watchmego@latest
```

## Usage

```bash
# To get help run it with the -h flag
watchmego -h

#      &*(&   (&#**********#&(    &&&
#      %&&&****/&#*********%&/***&(&&&
#    & #&.       &*&*****&*&..      &.%/
#   & &.          ,*&***&*,           & %
#   &..    &&%     &&&(&#&.    (&&      &
#   & %.          **&%#&&/.           & &
#    &&&..       &#%/,,,,&/&.        & &
#     %&/**//**&&***(,%****/&(//#%/*(&/
#     /*******************************%
#     ***************,,,**************&
#     **********............**********&
#     ********................********&
#     *******.................********&
#     /*****,..................*******&
#
# Usage: watchmego [flags][directory]
# Flags:
#   -h, --help		Print this help message
#   -v, --version		Print version information
#   -c, --check		Check makefile for compatibility
#   -C, --create		Create makefile for directory
#   -d, --debug		Print debug information

# watchmego will use the current directory if none is specified
# note that the directory is the only argument every other option
# is specified by flags

watchmego ./subdir

# You can also log all debug messages if you ever have any trouble
watchmego -d ./subdir

# If you have no makefile in the directory, create it by using the
# capital -C or --create flag
watchmego -C ./subdir

# If you want to check if your makefile is compatible use the
# lowercase -c
watchmego -c ./subdir

# [ERROR] Makefile not found
# Create Makefile? (y/n): y
# Makefile created
# [INFO] Please edit your makefile and add a watch target.
# [INFO] Exiting ...

```

## Makefile

```makefile
# Binary or command you want to run
WATCHMEGO=./runme

# You can use multiple commands too
WATCHMEGO=./runme ; ./andme ; ./ormeandthiscmdnexttome ; ls

# define a target called watch which will be called as a prebuild task
watch-build:
        @echo "I am running first :P"

# this directive will be generated when you use the cli to create your Makefile
# and is used for a more convenient way of calling watchmego
watch: watch-build
    @watchmego .
```

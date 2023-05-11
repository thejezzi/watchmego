package cli

import (
	"flag"
	"fmt"
	"github.com/thejezzi/watchmego/lib/logger"
	"os"
)

type Args struct {
	Dir     string
	Check   bool
	Create  bool
	Debug   bool
	Version bool
	Help    bool
}

var (
	dir     string
	check   bool
	create  bool
	debug   bool
	version bool
	help    bool
)

func PrintHelp() {
	fmt.Println("Usage: wmg [flags][directory]")
	fmt.Println("Flags:")
	fmt.Println("  -h, --help\t\tPrint this help message")
	fmt.Println("  -v, --version\t\tPrint version information")
	fmt.Println("  -c, --check\t\tCheck makefile for compatibility")
	fmt.Println("  -C, --create\t\tCreate makefile for directory")
	fmt.Println("  -d, --debug\t\tPrint debug information")
}

func ParseArgs(args []string) Args {
	flag.BoolVar(&check, "check", false, "Check makefile for compatibility")
	flag.BoolVar(&check, "c", false, "Check makefile for compatibility")
	flag.BoolVar(&create, "create", false, "Create makefile for directory")
	flag.BoolVar(&create, "C", false, "Create makefile for directory")
	flag.BoolVar(&debug, "debug", false, "Print debug information")
	flag.BoolVar(&debug, "d", false, "Print debug information")
	flag.BoolVar(&version, "version", false, "Print version information")
	flag.BoolVar(&version, "v", false, "Print version information")
	flag.BoolVar(&help, "help", false, "Print help message")
	flag.BoolVar(&help, "h", false, "Print help message")
	flag.Parse()

	if help {
		PrintHelp()
		os.Exit(0)
	}

	if debug {
		logger.Debug("Debug mode enabled")
	}

	if len(flag.Args()) == 0 {
		dir = "."
	} else {
		dir = flag.Arg(0)
	}

	// Get directory from args which is the only non-flag argument

	return Args{dir, check, create, debug, version, help}

}

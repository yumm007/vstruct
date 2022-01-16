package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/afero"

	"vstruct/internal/generator"
	"vstruct/internal/parser"
)

// Version is the version of `accessory`, injected at build time.
var Version = ""

// newUsage returns a function to replace default usage function of FlagSet.
func newUsage(flags *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage of vstruct:\n")
		fmt.Fprintf(os.Stderr, "\tvstruct [flags] [directory]\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flags.PrintDefaults()
	}
}

// Execute executes a whole process of generating accessor codes.
func Execute(fs afero.Fs, args []string) {
	log.SetFlags(0 | log.Lshortfile)
	log.SetPrefix("vstruct: ")

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = newUsage(flags)
	version := flags.Bool("version", false, "show the version of vstruct")
	structName := flags.String("struct", "", "struct name; must be set")
	receiver := flags.String("receiver", "", "receiver name; default first letter of struct name")
	output := flags.String("output", "", "output file name; default <struct_name>_vstruct.go")

	if err := flags.Parse(args[1:]); err != nil {
		flags.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Fprintf(os.Stdout, "accessory version: %s\n", getVersion())
		os.Exit(0)
	}

	if structName == nil || len(*structName) == 0 {
		flags.Usage()
		os.Exit(1)
	}

	var dir string
	if cliArgs := flags.Args(); len(cliArgs) > 0 {
		dir = cliArgs[0]
	} else {
		// Default: process whole package in current directory.
		dir = "."
	}

	if !isDir(dir) {
		fmt.Fprintln(os.Stderr, "Specified argument is not a directory.")
		flags.Usage()
		os.Exit(1)
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		os.Exit(1)
	}

	if err = generator.Generate(fs, pkg, *structName, *output, *receiver); err != nil {
		log.Fatal(err)
	}
}

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func getVersion() string {
	if Version != "" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "unknown"
	}

	return info.Main.Version
}
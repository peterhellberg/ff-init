package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

//go:embed all:content
var content embed.FS

type config struct {
	dir        string
	appID      string
	appName    string
	authorID   string
	authorName string
}

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stderr io.Writer) error {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.SetOutput(stderr)

	flags.Usage = func() {
		format := "Usage: %s [OPTION]... DIRECTORY\n\nOptions:\n"

		fmt.Fprintf(flags.Output(), format, os.Args[0])

		flags.PrintDefaults()
	}

	current, err := user.Current()
	if err != nil {
		return err
	}

	flags.StringVar(&cfg.appID, "app-id", "", "ID of the Firefly Zero app")
	flags.StringVar(&cfg.appName, "app-name", "", "Name of the Firefly Zero app")
	flags.StringVar(&cfg.authorID, "author-id", current.Username, "ID of the Firefly Zero author")
	flags.StringVar(&cfg.authorName, "author-name", current.Name, "Name of the Firefly Zero author")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return fmt.Errorf("no name given as the first argument")
	}

	cfg.dir = rest[0]

	{
		fallback := strings.TrimPrefix(cfg.dir, "ff-")

		if cfg.appID == "" {
			cfg.appID = fallback
		}

		if cfg.appName == "" {
			cfg.appName = fallback
		}
	}

	// Make sure that dir does not already exist
	if _, err := os.Stat(cfg.dir); !os.IsNotExist(err) {
		return fmt.Errorf("%q already exists", cfg.dir)
	}

	// Create the dir and dir/src
	if err := os.MkdirAll(cfg.dir+"/src", os.ModePerm); err != nil {
		return err
	}

	// Enter the new directory
	if err := os.Chdir(cfg.dir); err != nil {
		return err
	}

	entries, err := content.ReadDir("content")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			if err := writeFile(cfg, e.Name(), replacer); err != nil {
				return err
			}
		} else {
			if e.Name() == "src" {
				srcEntries, err := content.ReadDir("content/src")
				if err != nil {
					return err
				}

				for _, e := range srcEntries {
					if !e.IsDir() {
						if err := writeFile(cfg, "src/"+e.Name(), replacer); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return os.Chmod("spy.sh", 0o755)
}

func writeFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := content.ReadFile("content/" + name)
	if err != nil {
		return fmt.Errorf("writeFile: %w", err)
	}

	for i := range dataFuncs {
		data = dataFuncs[i](cfg, name, data)
	}

	return os.WriteFile(name, data, 0o644)
}

type dataFunc func(config, string, []byte) []byte

func replacer(cfg config, name string, data []byte) []byte {
	switch name {
	case "firefly.toml":
		data = replaceOne(data, "ff-app-id", cfg.appID)
		data = replaceOne(data, "ff-app-name", cfg.appName)
		data = replaceOne(data, "ff-author-id", cfg.authorID)
		data = replaceOne(data, "ff-author-name", cfg.authorName)

		return data
	case "build.zig":
		data = replaceOne(data, "ff-app-id", cfg.appID)
		data = replaceOne(data, "ff-author-id", cfg.authorID)

		return data
	case "build.zig.zon", "README.md":
		data = replaceOne(data, "ff-app-id", cfg.appID)
		data = replaceOne(data, "ff-author-id", cfg.authorID)

		return data
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	return bytes.Replace(data, []byte(old), []byte(new), 1)
}

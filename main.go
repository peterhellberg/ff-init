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

const (
	defaultHostname   = "play.c7.se"
	defaultServerRoot = "/var/www/play.c7.se"

	maxIDLength = 16
)

type config struct {
	dir        string
	appID      string
	appName    string
	authorID   string
	authorName string
	hostname   string
	serverRoot string
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
	flags.StringVar(&cfg.hostname, "hostname", defaultHostname, "The hostname to deploy the Firefly Zero app to")
	flags.StringVar(&cfg.serverRoot, "server-root", defaultServerRoot, "The root path on the server the app should be uploaded to")

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

	if err := validateID("app", cfg.appID); err != nil {
		return err
	}

	if err := validateID("author", cfg.authorID); err != nil {
		return err
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

func validateID(typ, id string) error {
	if id == "" {
		return fmt.Errorf("%s-id is empty", typ)
	}

	if len(id) > maxIDLength {
		return fmt.Errorf("%s-id too long", typ)
	}

	if strings.HasPrefix(id, "-") {
		return fmt.Errorf("%s-id has - prefix", typ)
	}

	if strings.HasSuffix(id, "-") {
		return fmt.Errorf("%s-id has - suffix", typ)
	}

	if strings.Contains(id, "--") {
		return fmt.Errorf("%s-id contains --", typ)
	}

	if !onlyAllowedRunes(id) {
		return fmt.Errorf("%s-id contains something not allowed: %q", typ, id)
	}

	return nil
}

func onlyAllowedRunes(s string) bool {
	for _, c := range s {
		if c != '-' && (c < '0' || c > '9') && (c < 'a' || c > 'z') {
			return false
		}
	}
	return true
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
	case "Makefile":
		data = replaceOne(data, "ff-app-id", cfg.appID)
		data = replaceOne(data, "ff-author-id", cfg.authorID)
		data = replaceOne(data, "localhost", cfg.hostname)
		data = replaceOne(data, "/tmp", cfg.serverRoot)

		return data
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

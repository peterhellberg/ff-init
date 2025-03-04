package main

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//go:embed all:content
var content embed.FS

//go:embed all:minimal
var minimal embed.FS

const (
	defaultHostname   = "play.c7.se"
	defaultServerRoot = "/var/www/play.c7.se"
	defaultMinimal    = false

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
	minimal    bool
	zon        ZON
}

type ZON struct {
	name        string
	fingerprint string
}

func parse(args []string, stderr io.Writer) (config, error) {
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
		return cfg, err
	}

	flags.StringVar(&cfg.appID, "app-id", "", "ID of the Firefly Zero app")
	flags.StringVar(&cfg.appName, "app-name", "", "Name of the Firefly Zero app")
	flags.StringVar(&cfg.authorID, "author-id", current.Username, "ID of the Firefly Zero author")
	flags.StringVar(&cfg.authorName, "author-name", current.Name, "Name of the Firefly Zero author")
	flags.StringVar(&cfg.hostname, "hostname", defaultHostname, "The hostname to deploy the Firefly Zero app to")
	flags.StringVar(&cfg.serverRoot, "server-root", defaultServerRoot, "The root path on the server the app should be uploaded to")
	flags.BoolVar(&cfg.minimal, "minimal", defaultMinimal, "Should the minimal template be used or not")

	if err := flags.Parse(args[1:]); err != nil {
		return cfg, err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return cfg, fmt.Errorf("no name given as the first argument")
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
		return cfg, err
	}

	if err := validateID("author", cfg.authorID); err != nil {
		return cfg, err
	}

	id := fmt.Sprintf("%s.%s", cfg.authorID, cfg.appID)

	zon, err := initZON(id)
	if err != nil {
		return cfg, err
	}

	cfg.zon = zon

	return cfg, nil
}

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stderr io.Writer) error {
	cfg, err := parse(args, stderr)
	if err != nil {
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

	var (
		writeFile = contentWriteFile
		srcFS     = content
		srcBase   = "content"
	)

	if cfg.minimal {
		writeFile = minimalWriteFile
		srcFS = minimal
		srcBase = "minimal"
	}

	entries, err := srcFS.ReadDir(srcBase)
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
				srcEntries, err := srcFS.ReadDir(srcBase + "/src")
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

type writeFileFunc func(cfg config, name string, dataFuncs ...dataFunc) error

func contentWriteFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := content.ReadFile("content/" + name)
	if err != nil {
		return fmt.Errorf("contentWriteFile: %w", err)
	}

	for i := range dataFuncs {
		data = dataFuncs[i](cfg, name, data)
	}

	return os.WriteFile(name, data, 0o644)
}

func minimalWriteFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := minimal.ReadFile("minimal/" + name)
	if err != nil {
		return fmt.Errorf("minimalWriteFile: %w", err)
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
	case "README.md", "spy.sh":
		data = replaceOne(data, "ff-app-id", cfg.appID)
		data = replaceOne(data, "ff-author-id", cfg.authorID)

		return data
	case "build.zig.zon":
		data = replaceOne(data, ".ff_app_name", cfg.zon.name)
		data = replaceOne(data, "0x9cdb93c8c3a4327e", cfg.zon.fingerprint)

		return data
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	return bytes.Replace(data, []byte(old), []byte(new), 1)
}

func initZON(dir string) (ZON, error) {
	tmp, err := os.MkdirTemp("", "ff-init-")
	if err != nil {
		return ZON{}, err
	}
	defer os.RemoveAll(tmp)

	cwd, err := os.Getwd()
	if err != nil {
		return ZON{}, err
	}
	defer os.Chdir(cwd)

	tmpDir := filepath.Join(tmp, dir)

	if err := os.Mkdir(tmpDir, 0o755); err != nil {
		return ZON{}, err
	}

	if err := os.Chdir(tmpDir); err != nil {
		return ZON{}, err
	}

	cmd := exec.Command("zig", "init")

	if err := cmd.Run(); err != nil {
		return ZON{}, err
	}

	zonPath := filepath.Join(tmpDir, "build.zig.zon")

	return extractZON(zonPath)
}

func extractZON(zonPath string) (ZON, error) {
	var zon ZON

	f, err := os.Open(zonPath)
	if err != nil {
		return zon, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if prefix := ".name = "; strings.Contains(text, prefix) {
			zon.name = strings.TrimSuffix(strings.TrimPrefix(text, prefix), ",")
		}

		if prefix := ".fingerprint = "; strings.Contains(text, prefix) {
			fingerprint, _, _ := strings.Cut(strings.TrimPrefix(text, prefix), ",")
			zon.fingerprint = fingerprint
		}
	}

	return zon, nil
}

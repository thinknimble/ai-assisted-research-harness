package main

import (
	"fmt"
	"io"
	"sort"
)

func runRepos(stdout io.Writer) error {
	cfg, err := loadGlobalConfigIfExists()
	if err != nil {
		return err
	}

	if cfg == nil || len(cfg.Repos) == 0 {
		fmt.Fprintln(stdout, "No repos registered. Run `research-assistant init` to create one.")
		return nil
	}

	// Sort names for stable output
	names := make([]string, 0, len(cfg.Repos))
	for name := range cfg.Repos {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		path := cfg.Repos[name]
		if name == cfg.Default {
			fmt.Fprintf(stdout, "* %s  %s\n", name, path)
		} else {
			fmt.Fprintf(stdout, "  %s  %s\n", name, path)
		}
	}

	return nil
}

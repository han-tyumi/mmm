package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/han-tyumi/mcf"
	"github.com/spf13/cobra"
)

var useSlug, useSearch bool

var getCmd = &cobra.Command{
	Use:   "get [-s] ...{id | slug}",
	Short: "Downloads the specified mods by ID",
	Run: func(cmd *cobra.Command, args []string) {
		version = cfg.GetString("version")

		var mods []mcf.Mod
		var err error

		if useSearch {
			mods, err = modsBySearch(args)
		} else if useSlug {
			mods, err = modsBySlug(args)
		} else {
			mods, err = modsByID(args)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(mods) == 0 {
			fmt.Fprintln(os.Stderr, "no mods found")
			os.Exit(1)
		}

		for i := range mods {
			mod := mods[i]

			modFile, err := findFile(&mod)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			name := path.Base(modFile.URL)

			fmt.Printf("downloading %s ...\n", name)
			res, err := http.Get(modFile.URL)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				fmt.Fprintln(os.Stderr, res.Status)
				os.Exit(1)
			}

			file, err := os.Create(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer file.Close()

			if _, err := io.Copy(file, res.Body); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&useSlug, "slug", "s", false, "Add mods based on its slug")
	getCmd.Flags().BoolVarP(&useSearch, "search", "S", false, "Add mods based on search terms")
	getCmd.Flags().StringVarP(&version, "version", "v", "", "Download the latest for a Minecraft version")

	cfg.BindPFlag("version", getCmd.Flags().Lookup("version"))
}

func modsByID(args []string) ([]mcf.Mod, error) {
	var ids []uint

	for i := range args {
		arg := args[i]

		id, err := strconv.ParseUint(arg, 10, 0)
		if err != nil {
			return nil, err
		}

		ids = append(ids, uint(id))
	}

	if len(ids) == 0 {
		return nil, errors.New("no ids specified")
	}

	mods, err := mcf.Many(ids)
	if err != nil {
		return nil, err
	}

	return mods, nil
}

func modsBySearch(args []string) ([]mcf.Mod, error) {
	var mods []mcf.Mod

	for i := range args {
		arg := args[i]

		r, err := mcf.Search(&mcf.SearchParams{
			PageSize: 1,
			Search:   arg,
			Version:  version,
		})
		if err != nil {
			return nil, err
		} else if len(r) == 0 {
			return nil, fmt.Errorf("%s not found", arg)
		}

		mods = append(mods, r[0])
	}

	return mods, nil
}

var allMods []mcf.Mod

func getAllMods() ([]mcf.Mod, error) {
	if allMods != nil {
		return allMods, nil
	}

	mods, err := mcf.Search(&mcf.SearchParams{
		Version: version,
	})
	if err != nil {
		return nil, err
	}

	allMods = mods
	return allMods, nil
}

func modsBySlug(args []string) ([]mcf.Mod, error) {
	var mods []mcf.Mod

	for i := range args {
		arg := args[i]

		mod, err := findBySlug(arg)
		if err != nil {
			return nil, err
		}

		mods = append(mods, *mod)
	}

	return mods, nil
}

func findBySlug(slug string) (*mcf.Mod, error) {
	mods, err := getAllMods()
	if err != nil {
		return nil, err
	}

	for i := range mods {
		mod := mods[i]
		if mod.Slug == slug {
			return &mod, nil
		}
	}

	return nil, fmt.Errorf("could not find mod with %s slug", slug)
}

func findFile(mod *mcf.Mod) (*mcf.ModFile, error) {
	if len(mod.LatestFiles) == 0 {
		return nil, fmt.Errorf("no files for %s", mod.Name)
	}

	if version == "" {
		return &mod.LatestFiles[0].ModFile, nil
	}

	files, err := mcf.Files(mod.ID)
	if err != nil {
		return nil, err
	}

	var latest *mcf.ModFile

	for i := range files {
		file := files[i]
		for j := range file.Versions {
			if file.Versions[j] != version {
				continue
			}

			if latest == nil || file.Uploaded.After(latest.Uploaded) {
				latest = &file
			}

			break
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("%s does not support %s", mod.Name, version)
	}

	return latest, nil
}

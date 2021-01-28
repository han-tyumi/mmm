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

var useSlug bool

var getCmd = &cobra.Command{
	Use:   "get [-s] ...{id | slug}",
	Short: "Downloads the specified mods by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var mods []mcf.Mod
		var err error

		if useSlug {
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

	getCmd.Flags().BoolVarP(&useSlug, "slug", "s", false, "Add a mod based on its URL slug")
	getCmd.Flags().StringVarP(&version, "version", "v", "", "Download the latest for a Minecraft version")
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

func modsBySlug(args []string) ([]mcf.Mod, error) {
	var mods []mcf.Mod

	for i := range args {
		arg := args[i]

		m, err := mcf.Search(&mcf.SearchParams{
			PageSize: 1,
			Search:   arg,
		})
		if err != nil {
			return nil, err
		} else if len(m) == 0 {
			return nil, fmt.Errorf("%s not found", arg)
		}

		mods = append(mods, m[0])
	}

	return mods, nil
}

func findFile(mod *mcf.Mod) (*mcf.ModFile, error) {
	if len(mod.LatestFiles) == 0 {
		return nil, fmt.Errorf("no files for %s", mod.Name)
	}

	if version == "" {
		return &mod.LatestFiles[0], nil
	}

	files := mod.LatestFiles
	for i := range files {
		file := files[i]
		for j := range file.Versions {
			if file.Versions[j] == version {
				return &file, nil
			}
		}
	}

	return nil, fmt.Errorf("%s not found for %s", version, mod.Name)
}

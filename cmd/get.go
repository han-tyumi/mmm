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
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useSlug, useSearch bool

var getCmd = &cobra.Command{
	Use:   "get [-s] ...{id | slug}",
	Short: "Downloads the specified mods by ID",
	Run: func(cmd *cobra.Command, args []string) {
		version = viper.GetString("version")

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
			utils.Error(err)
		}

		if len(mods) == 0 {
			utils.Error("no mods found")
		}

		for i := range mods {
			mod := mods[i]

			modFile, err := findLatestByMod(&mod)
			if err != nil {
				utils.Error(err)
			}

			name := path.Base(modFile.URL)

			fmt.Printf("downloading %s ...\n", name)
			res, err := http.Get(modFile.URL)
			if err != nil {
				utils.Error(err)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				utils.Error(res.Status)
			}

			file, err := os.Create(name)
			if err != nil {
				utils.Error(err)
			}
			defer file.Close()

			if _, err := io.Copy(file, res.Body); err != nil {
				utils.Error(err)
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

	viper.BindPFlag("version", getCmd.Flags().Lookup("version"))
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

func findLatestByMod(mod *mcf.Mod) (*mcf.ModFile, error) {
	if version == "" {
		if len(mod.LatestFiles) == 0 {
			return nil, fmt.Errorf("no files for %s", mod.Name)
		}

		return &mod.LatestFiles[0].ModFile, nil
	}

	return findLatestByID(mod.ID, mod.Name)
}

func findLatestByID(id uint, name string) (*mcf.ModFile, error) {
	files, err := mcf.Files(id)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files for %s", name)
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
		return nil, fmt.Errorf("%s does not support %s", name, version)
	}

	return latest, nil
}

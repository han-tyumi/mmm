package get

import (
	"fmt"
	"strconv"

	"github.com/han-tyumi/mcf"
)

var allMods = make(map[string][]mcf.Mod)

// AllMods returns all the mods for a given Minecraft version.
func AllMods(version string) ([]mcf.Mod, error) {
	if mods, ok := allMods[version]; ok {
		return mods, nil
	}

	mods, err := mcf.Search(&mcf.SearchParams{
		Version: version,
	})
	if err != nil {
		return nil, err
	}

	allMods[version] = mods
	return mods, nil
}

// ModsByArgs returns all mods for some given arguments and a Minecraft version.
func ModsByArgs(args []string, version string) ([]mcf.Mod, error) {
	ids := make([]uint, 0)
	slugs := make([]string, 0)

	for i := range args {
		arg := args[i]

		if id, err := strconv.ParseUint(arg, 10, 0); err == nil {
			ids = append(ids, uint(id))
		} else {
			slugs = append(slugs, arg)
		}
	}

	if len(ids) == 0 {
		return ModsBySlug(slugs, version)
	} else if len(slugs) == 0 {
		return mcf.Many(ids)
	}

	slugMods, err := ModsBySlug(slugs, version)
	if err != nil {
		return nil, err
	}

	idMods, err := mcf.Many(ids)
	if err != nil {
		return nil, err
	}

	return append(slugMods, idMods...), nil
}

// ModsBySlug returns the mods corresponding to each URL slug.
func ModsBySlug(slugs []string, version string) ([]mcf.Mod, error) {
	mods := make([]mcf.Mod, len(slugs))

	for i := range slugs {
		mod, err := ModBySlug(slugs[i], version)
		if err != nil {
			return nil, err
		}

		mods[i] = *mod
	}

	return mods, nil
}

// ModBySlug returns a mod by its URL slug.
func ModBySlug(slug, version string) (*mcf.Mod, error) {
	mods, err := AllMods(version)
	if err != nil {
		return nil, err
	}

	for i := range mods {
		mod := mods[i]
		if mod.Slug == slug {
			return &mod, nil
		}
	}

	return nil, fmt.Errorf("could not find mod with slug, %s", slug)
}

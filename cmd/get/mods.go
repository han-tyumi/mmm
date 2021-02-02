package get

import (
	"errors"
	"fmt"

	"github.com/han-tyumi/mcf"
)

var allMods map[string][]mcf.Mod

// AllMods returns all the mods for a given Minecraft version.
func AllMods(version string) (mods []mcf.Mod, err error) {
	if m, ok := allMods[version]; ok {
		return m, nil
	}

	mods, err = mcf.Search(&mcf.SearchParams{
		Version: version,
	})

	if err == nil {
		allMods[version] = mods
	}

	return
}

// ModsBySlug returns the mods corresponding to each URL slug.
func ModsBySlug(slugs []string, version string) (mods []mcf.Mod, err error) {
	for i := range slugs {
		mod, err := ModBySlug(slugs[i], version)
		if err != nil {
			return nil, err
		}

		mods = append(mods, *mod)
	}

	return
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

// ModsByID returns the mods corresponding to each ID.
func ModsByID(ids []uint) (mods []mcf.Mod, err error) {
	if len(ids) == 0 {
		return nil, errors.New("no ids specified")
	}

	return mcf.Many(ids)
}

// ModsBySearch returns the first mod search result for each search using a Minecraft version.
func ModsBySearch(searches []string, version string) (mods []mcf.Mod, err error) {
	if len(searches) == 0 {
		return nil, errors.New("no searches specified")
	}

	for i := range searches {
		search := searches[i]

		results, err := mcf.Search(&mcf.SearchParams{
			PageSize: 1,
			Search:   search,
			Version:  version,
		})
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			return nil, fmt.Errorf("%s not found", search)
		}

		mods = append(mods, results[0])
	}

	return
}

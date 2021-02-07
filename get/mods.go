/*
Package get provides functionality for getting mods and their files based on user input.

Copyright © 2021 Matthew Champagne <mmchamp95@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package get

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/utils"
)

var versionMods = make(map[string][]mcf.Mod)
var versionModsMu sync.Mutex

var versionSlugMod = make(map[string]map[string]*mcf.Mod)
var versionSlugModMu sync.Mutex

// AllMods returns all the mods for a given Minecraft version.
func AllMods(version string) ([]mcf.Mod, error) {
	versionModsMu.Lock()
	mods, ok := versionMods[version]
	versionModsMu.Unlock()

	if ok {
		return mods, nil
	}

	mods, err := mcf.Search(&mcf.SearchParams{
		Version: version,
	})
	if err != nil {
		return nil, err
	}

	versionModsMu.Lock()
	versionMods[version] = mods
	versionModsMu.Unlock()

	versionSlugModMu.Lock()
	versionSlugMod[version] = make(map[string]*mcf.Mod)
	versionSlugModMu.Unlock()

	return mods, nil
}

// ModsByArgs returns all mods for some given arguments and a Minecraft version.
func ModsByArgs(args []string, version string) ([]mcf.Mod, error) {
	ids := make([]uint, 0)
	slugs := make([]string, 0)

	var idsMu, slugsMu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(args))

	for i := range args {
		i := i

		go func() {
			defer wg.Done()

			arg := args[i]

			if id, err := strconv.ParseUint(arg, 10, 0); err == nil {
				idsMu.Lock()
				ids = append(ids, uint(id))
				idsMu.Unlock()
			} else {
				slugsMu.Lock()
				slugs = append(slugs, arg)
				slugsMu.Unlock()
			}
		}()
	}

	wg.Wait()

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
	ch := utils.NewErrCh(len(slugs))

	for i := range slugs {
		i := i

		go ch.Do(func() error {
			mod, err := ModBySlug(slugs[i], version)
			if err != nil {
				return err
			}

			mods[i] = *mod
			return nil
		})
	}

	if err := ch.Wait(func(err error) error {
		return err
	}); err != nil {
		return nil, err
	}

	return mods, nil
}

// ModBySlug returns a mod by its URL slug.
func ModBySlug(slug, version string) (*mcf.Mod, error) {
	mods, err := AllMods(version)
	if err != nil {
		return nil, err
	}

	versionSlugModMu.Lock()
	mod, ok := versionSlugMod[version][slug]
	versionSlugModMu.Unlock()

	if ok {
		return mod, nil
	}

	ch := make(chan *mcf.Mod)
	for i := range mods {
		i := i

		go func() {
			mod := mods[i]

			versionSlugModMu.Lock()
			_, ok := versionSlugMod[version][mod.Slug]
			versionSlugModMu.Unlock()

			if ok {
				ch <- nil
				return
			}

			versionSlugModMu.Lock()
			versionSlugMod[version][mod.Slug] = &mod
			versionSlugModMu.Unlock()

			if mod.Slug == slug {
				ch <- &mod
			} else {
				ch <- nil
			}
		}()
	}

	for range mods {
		if mod := <-ch; mod != nil {
			return mod, nil
		}
	}

	return nil, fmt.Errorf("could not find mod with slug, %s", slug)
}

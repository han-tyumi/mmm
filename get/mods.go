package get

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/utils"
)

var allMods = make(map[string][]mcf.Mod)
var allModsMu sync.Mutex

// AllMods returns all the mods for a given Minecraft version.
func AllMods(version string) ([]mcf.Mod, error) {
	allModsMu.Lock()
	mods, ok := allMods[version]
	allModsMu.Unlock()

	if ok {
		return mods, nil
	}

	mods, err := mcf.Search(&mcf.SearchParams{
		Version: version,
	})
	if err != nil {
		return nil, err
	}

	allModsMu.Lock()
	allMods[version] = mods
	allModsMu.Unlock()

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

	ch := make(chan *mcf.Mod)
	for i := range mods {
		i := i

		go func() {
			mod := mods[i]
			if mod.Slug == slug {
				ch <- &mod
			}
			ch <- nil
		}()
	}

	for range mods {
		if mod := <-ch; mod != nil {
			return mod, nil
		}
	}

	return nil, fmt.Errorf("could not find mod with slug, %s", slug)
}

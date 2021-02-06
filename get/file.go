package get

import (
	"fmt"
	"sync"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/utils"
)

// LatestFileByMod returns the latest mod file for a mod and an optional Minecraft version.
func LatestFileByMod(version string, mod *mcf.Mod) (*mcf.ModFile, error) {
	if version == "" {
		if len(mod.LatestFiles) == 0 {
			return nil, fmt.Errorf("no files for %s", mod.Name)
		}

		return &mod.LatestFiles[0].ModFile, nil
	}

	return LatestFileByID(version, mod.ID, mod.Name)
}

// LatestFileByID returns the latest mod file for a mod's ID and a Minecraft version.
func LatestFileByID(version string, id uint, name string) (*mcf.ModFile, error) {
	files, err := mcf.Files(id)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files for %s", name)
	}

	var latest *mcf.ModFile
	var latestMu sync.Mutex
	latestCh := make(chan *mcf.ModFile)

	for i := range files {
		i := i

		go func() {
			file := files[i]
			versionCh := make(chan bool)

			for j := range file.Versions {
				j := j

				go func() {
					versionCh <- file.Versions[j] == version
				}()
			}

			for range file.Versions {
				latestMu.Lock()
				isBefore := latest != nil && file.Uploaded.Before(latest.Uploaded)
				latestMu.Unlock()

				if isBefore {
					break
				} else if <-versionCh {
					latestCh <- &file
					return
				}
			}

			latestCh <- nil
		}()
	}

	for range files {
		if file := <-latestCh; file != nil && (latest == nil || file.Uploaded.After(latest.Uploaded)) {
			latestMu.Lock()
			latest = file
			latestMu.Unlock()
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("%s does not support %s", name, version)
	}

	return latest, nil
}

// LatestFileCallback is called concurrently with a mod and its latest file.
type LatestFileCallback func(mod *mcf.Mod, latest *mcf.ModFile) error

// LatestFileForEachMod concurrently calls cb with the latest file for each mod and the given Minecraft version.
func LatestFileForEachMod(mods []mcf.Mod, version string, cb LatestFileCallback) error {
	ch := utils.NewErrCh(len(mods))
	for i := range mods {
		i := i

		go ch.Do(func() error {
			mod := mods[i]

			latest, err := LatestFileByMod(version, &mod)
			if err != nil {
				return err
			}

			return cb(&mod, latest)
		})
	}

	return ch.Wait(func(err error) error {
		return err
	})
}

// LatestFileForEachArg concurrently calls cb with the latest file for each id or slug argument and the given Minecraft version.
func LatestFileForEachArg(args []string, version string, cb LatestFileCallback) error {
	mods, err := ModsByArgs(args, version)
	if err != nil {
		return err
	}

	return LatestFileForEachMod(mods, version, cb)
}

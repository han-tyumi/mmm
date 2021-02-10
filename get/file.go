/*
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
	"errors"
	"sync"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/utils"
)

// ErrNoFiles is returned when there are no files for a mod.
var ErrNoFiles = errors.New("no files")

// ErrVersionUnsupported is returned when a mod doesn't have any files supporting the specified version.
var ErrVersionUnsupported = errors.New("version unsupported")

// LatestFileByMod returns the latest mod file for a mod and an optional Minecraft version.
func LatestFileByMod(version string, mod *mcf.Mod) (*mcf.ModFile, error) {
	if version == "" {
		if len(mod.LatestFiles) == 0 {
			return nil, ErrNoFiles
		}

		return &mod.LatestFiles[0].ModFile, nil
	}

	return LatestFileByID(version, mod.ID)
}

// LatestFileByID returns the latest mod file for a mod's ID and a Minecraft version.
func LatestFileByID(version string, id uint) (*mcf.ModFile, error) {
	files, err := mcf.Files(id)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoFiles
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
		return nil, ErrVersionUnsupported
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

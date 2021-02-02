package get

import (
	"fmt"

	"github.com/han-tyumi/mcf"
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

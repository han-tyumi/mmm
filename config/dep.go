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

package config

import (
	"os"
	"time"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/get"
)

// Dependency is a mod managed in the user's dependency configuration file.
type Dependency struct {
	ID       uint      `mapstructure:"id"`
	Name     string    `mapstructure:"name"`
	URL      string    `mapstructure:"url"`
	File     string    `mapstructure:"file"`
	Uploaded time.Time `mapstructure:"uploaded"`
	Size     uint      `mapstructure:"size"`
}

// Download downloads the dependency to the current working directory.
func (d *Dependency) Download() error {
	return download.FromURL(d.File, d.URL)
}

// Downloaded returns whether the dependency has already been downloaded.
func (d *Dependency) Downloaded() (bool, error) {
	info, err := os.Stat(d.File)
	return err == nil && info.Size() == int64(d.Size), err
}

// SameFile returns whether the dependency is using the same mod file.
func (d *Dependency) SameFile(file *mcf.ModFile) bool {
	return file.Name == d.File && file.Uploaded == d.Uploaded && file.Size == d.Size
}

// SameDepFile returns whether a dependency has the same file information.
func (d *Dependency) SameDepFile(dep *Dependency) bool {
	return dep.File == d.File && dep.Uploaded == d.Uploaded && dep.Size == d.Size
}

// UpdateFile updates the dependency's file information.
func (d *Dependency) UpdateFile(file *mcf.ModFile) {
	d.URL = file.URL
	d.File = file.Name
	d.Uploaded = file.Uploaded
	d.Size = file.Size
}

// RemoveFile removes the associated mod file.
func (d *Dependency) RemoveFile() error {
	return os.Remove(d.File)
}

// LatestFile returns the latest mod file for this dependency.
func (d *Dependency) LatestFile(version string) (*mcf.ModFile, error) {
	return get.LatestFileByID(version, d.ID)
}

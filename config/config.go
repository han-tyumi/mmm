/*
Package config provides functionality for managing a user's mods.

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
	"errors"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/viper"
)

// ModsKey is the key used to store managed mods.
const ModsKey = "mods"

// ErrNotSet is returned when a mod slug Dependency has not been set.
var ErrNotSet = errors.New("not set")

// ErrNoMods is returned when there are no mods being managed.
var ErrNoMods = errors.New("no mods being managed")

// DependencyMap allows for safe concurrent usage of the map of a user's mod dependencies.
type DependencyMap struct {
	deps map[string]*Dependency
	mu   sync.Mutex
}

var viperMu sync.Mutex

// DepMapSync safely returns a map of mod slugs to Dependencies for the user's configuration file.
func DepMapSync() (map[string]*Dependency, error) {
	raw := map[string]*Dependency{}

	viperMu.Lock()
	err := viper.UnmarshalKey(ModsKey, &raw,
		viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
	viperMu.Unlock()

	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, ErrNoMods
	}

	return raw, nil
}

// DepMap safely returns a concurrency safe map of mod slugs to Dependencies for the user's configuration file.
func DepMap() (*DependencyMap, error) {
	deps, err := DepMapSync()
	if err != nil {
		return nil, err
	}

	depMap := &DependencyMap{
		deps: deps,
	}

	return depMap, nil
}

// Clone returns a copy of this dependency map.
func (d *DependencyMap) Clone() *DependencyMap {
	deps := make(map[string]*Dependency, len(d.deps))

	d.mu.Lock()
	for slug, dep := range d.deps {
		deps[slug] = dep.Clone()
	}
	d.mu.Unlock()

	return &DependencyMap{
		deps: deps,
	}
}

// Each calls the provided function for each mapped slug and dependency.
func (d *DependencyMap) Each(cb func(slug string, dep *Dependency)) {
	for slug, dep := range d.deps {
		cb(slug, dep)
	}
}

// Len returns the length of the DependencyMap.
func (d *DependencyMap) Len() int {
	return len(d.deps)
}

// Get safely returns a Dependency for a given mod's slug if it's present in the map.
func (d *DependencyMap) Get(slug string) (*Dependency, bool) {
	d.mu.Lock()
	dep, ok := d.deps[slug]
	d.mu.Unlock()

	return dep, ok
}

// Set safely sets the Dependency for a given mod's slug.
func (d *DependencyMap) Set(slug string, dep *Dependency) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.deps[slug] = dep
}

// Delete safely removes a Dependency for a given mod's slug if it's present in the map.
func (d *DependencyMap) Delete(slug string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.deps, slug)
}

// Write safely writes all Dependency map information.
func (d *DependencyMap) Write() error {
	viperMu.Lock()
	defer viperMu.Unlock()

	viper.Set(ModsKey, d.deps)
	return viper.WriteConfig()
}

// Dep safely returns a Dependency for a given mod's slug.
func Dep(slug string) (*Dependency, error) {
	key := ModsKey + "." + slug

	viperMu.Lock()
	isSet := viper.IsSet(key)
	viperMu.Unlock()

	if !isSet {
		return nil, ErrNotSet
	}

	dep := &Dependency{}

	viperMu.Lock()
	err := viper.UnmarshalKey(key, dep,
		viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
	viperMu.Unlock()

	if err != nil {
		return nil, err
	}
	return dep, nil
}

// SetDep safely sets Dependency information for a given mod's slug.
func SetDep(slug string, dep *Dependency) error {
	viperMu.Lock()
	defer viperMu.Unlock()

	viper.Set(ModsKey+"."+slug, dep)
	return viper.WriteConfig()
}

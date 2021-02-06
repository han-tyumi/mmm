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

// Dependency is a mod managed in the user's dependency configuration file.
type Dependency struct {
	ID       uint      `mapstructure:"id"`
	Name     string    `mapstructure:"name"`
	File     string    `mapstructure:"file"`
	Uploaded time.Time `mapstructure:"uploaded"`
	Size     uint      `mapstructure:"size"`
}

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
	raw, err := DepMapSync()
	if err != nil {
		return nil, err
	}

	depMap := &DependencyMap{
		deps: raw,
	}

	return depMap, nil
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
	viperMu.Lock()
	isSet := viper.IsSet(slug)
	viperMu.Unlock()

	if !isSet {
		return nil, ErrNotSet
	}

	dep := &Dependency{}

	viperMu.Lock()
	err := viper.UnmarshalKey(ModsKey+"."+slug, dep,
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

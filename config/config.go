package config

import (
	"errors"
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

// Deps returns a map of mod slugs to Dependencies for the user's dependency configuration file.
func Deps() (map[string]*Dependency, error) {
	deps := map[string]*Dependency{}
	err := viper.UnmarshalKey(ModsKey, &deps,
		viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
	if err != nil {
		return nil, err
	} else if len(deps) == 0 {
		return nil, ErrNoMods
	}
	return deps, nil
}

// Dep returns a Dependency for a given mod's slug.
func Dep(slug string) (*Dependency, error) {
	if !viper.IsSet(slug) {
		return nil, ErrNotSet
	}

	dep := &Dependency{}
	err := viper.UnmarshalKey(ModsKey+"."+slug, dep,
		viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
	if err != nil {
		return nil, err
	}

	return dep, nil
}

// SetDeps sets all Dependency information by mod slugs.
func SetDeps(deps map[string]*Dependency) error {
	viper.Set(ModsKey, deps)
	return viper.WriteConfig()
}

// SetDep sets Dependency information for a given mod's slug.
func SetDep(slug string, dep *Dependency) error {
	viper.Set(ModsKey+"."+slug, dep)
	return viper.WriteConfig()
}

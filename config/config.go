package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ShulkerboxUrl   string   `yaml:"shulker_box_url"`
	PluginCopyPaths []string `yaml:"plugin_paths"`
}

func getDefaultConfig() Config {
	return Config{
		ShulkerboxUrl: "https://minecraft-hgl-drew.s3.amazonaws.com/shulkerbox.zip",
	}
}

func ReadConfigFromFile(path string) (Config, error) {
	var config Config
	io, err := os.Open(filepath.FromSlash(path))
	if errors.Is(err, os.ErrNotExist) {
		return getDefaultConfig(), nil
	}

	decoder := yaml.NewDecoder(io)
	err = decoder.Decode(&config)

	if err != nil {
		return Config{}, err
	}

	finalConfig := mergeConfigs(getDefaultConfig(), config)

	return finalConfig, nil
}

func mergeConfigs(defaultConfig, userConfig Config) Config {
	// Create a copy of the defaultConfig to avoid modifying it directly
	mergedConfig := defaultConfig

	// Iterate over the fields of the Config struct using reflection
	configType := reflect.TypeOf(mergedConfig)
	configValue := reflect.ValueOf(&mergedConfig).Elem()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		userValue := reflect.ValueOf(userConfig).Field(i).Interface()

		// If the user-configured field is not the zero value for its type,
		// overwrite the corresponding field in the mergedConfig
		if !reflect.DeepEqual(userValue, reflect.Zero(field.Type).Interface()) {
			configValue.Field(i).Set(reflect.ValueOf(userValue))
		}
	}
	return mergedConfig
}

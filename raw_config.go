package main

import (
	"github.com/BurntSushi/toml"
)

type RawSettings struct {
	BaseUrl *string   `toml:"base_url"`
	Output  *string   `toml:"output"`
	LoadEnv *string   `toml:"load_env"`
	Include *[]string `toml:"include"`
}

type RawRequestConfig struct {
	Args     []string `toml:"args"`
	Method   *string  `toml:"method"`
	Url      *string  `toml:"url"`
	Path     *string  `toml:"path"`
	Select   []string `toml:"select"`
	Headers  []string `toml:"headers"`
	Body     *string  `toml:"body"`
	BodyType *string  `toml:"body_type"`
	Extend   *string  `toml:"extend"`
}

type RawConfig struct {
	Settings *RawSettings                `toml:"settings"`
	Requests map[string]RawRequestConfig `toml:"req"`
}

func ParseRawConfig(path string) (RawConfig, error) {
	var config RawConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}
	return config, nil
}

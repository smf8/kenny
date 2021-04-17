package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/sirupsen/logrus"
)

type (
	//Config represents a kenny config instance.
	Config struct {
		Debug    bool     `koanf:"debug"`
		Logger   Logger   `koanf:"logger"`
		Recorder Recorder `koanf:"recorder"`
	}

	//Logger represents logger(logrus) config information.
	Logger struct {
		Level   logrus.Level `koanf:"level"`
		Enabled bool         `koanf:"enabled"`
	}

	// Recorder represents audio recorder settings.
	Recorder struct {
		NumberOfChannels int
		SampleRate       float64
	}
)

//New creates a new config instance with this order : default -> config.yml.
func New() Config {
	var instance Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(def, "konaf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading file: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	logrus.Infof("following configuration is loaded:\n%+v", instance)

	return instance
}

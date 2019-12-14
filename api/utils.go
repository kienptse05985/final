package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func BindJSON(r *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}

func LoadEnv(config interface{}, prefix string, source string) error {
	if err := LoadEnvFromDir(config, prefix, source); err != nil {
		return LoadEnvFromFile(config, prefix, source)
	}

	return nil
}

func LoadEnvFromFile(config interface{}, configPrefix, envPath string) (err error) {
	godotenv.Load(envPath)
	err = envconfig.Process(configPrefix, config)
	return
}

func LoadEnvFromDir(config interface{}, configPrefix, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	filePaths := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dir, f.Name()))
	}

	if err := godotenv.Load(filePaths...); err != nil {
		return err
	}
	return envconfig.Process(configPrefix, config)
}

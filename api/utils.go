package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func GetQuery(r *http.Request, key string) (string, bool) {
	if values, ok := GetQueryArray(r, key); ok {
		return values[0], ok
	}
	return "", false
}

func GetQueryArray(r *http.Request, key string) ([]string, bool) {
	if values, ok := r.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func BindJSON(r *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}

func StringInArray(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if strings.TrimSpace(str) == strings.TrimSpace(val) {
			return true
		}
	}
	return false
}

func Int64InArray(i int64, arr []int64) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if val == i {
			return true
		}
	}
	return false
}

func IntInArray(i int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if val == i {
			return true
		}
	}
	return false
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

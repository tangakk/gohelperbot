package tgbot

import (
	"encoding/json"
	"os"
)

type Config struct {
	Key          string   `json:"key"`
	StartButtons []string `json:"start"`
}

func (c *Config) ReadFromJSON(path string) error {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		return err
	}
	return nil
}

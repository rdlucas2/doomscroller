package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SaveData struct {
	Character   *Character
	Floor       int
	CompletedStories []string
}

func SavePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".doomscroller", "save.json")
}

func Save(data SaveData) error {
	path := SavePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func Load() (*SaveData, error) {
	b, err := os.ReadFile(SavePath())
	if err != nil {
		return nil, err
	}
	var d SaveData
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func DeleteSave() {
	os.Remove(SavePath())
}

func HasSave() bool {
	_, err := os.Stat(SavePath())
	return err == nil
}

package game

import "math/rand"

const MaxFloors = 5

type Session struct {
	Character *Character
	Map       *WorldMap
	Floor     int
	Rng       *rand.Rand
	Settings  Settings
	StepCount int
}

func NewSession(c *Character, seed int64, settings Settings) *Session {
	rng := rand.New(rand.NewSource(seed))
	s := &Session{
		Character: c,
		Floor:     1,
		Rng:       rng,
		Settings:  settings,
	}
	s.Map = GenerateMap(1, rng)
	return s
}

func (s *Session) NextFloor() {
	s.Floor++
	s.Map = GenerateMap(s.Floor, s.Rng)
}

func (s *Session) IsComplete() bool {
	return s.Floor > MaxFloors
}

func (s *Session) ToSaveData() SaveData {
	var completed []string
	// Collect completed story IDs from registry state
	return SaveData{
		Character:        s.Character,
		Floor:            s.Floor,
		CompletedStories: completed,
	}
}

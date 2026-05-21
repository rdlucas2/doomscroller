package game

type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyNormal
	DifficultyHard
)

func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Easy"
	case DifficultyNormal:
		return "Normal"
	case DifficultyHard:
		return "Hard"
	}
	return "Unknown"
}

type Settings struct {
	MusicVolume  int
	SoundVolume  int
	Difficulty   Difficulty
	ShowFullMap  bool
}

func DefaultSettings() Settings {
	return Settings{
		MusicVolume: 80,
		SoundVolume: 80,
		Difficulty:  DifficultyNormal,
		ShowFullMap: false,
	}
}

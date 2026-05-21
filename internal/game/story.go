package game

import "sync"

type TriggerType int

const (
	TriggerEnterRoom TriggerType = iota
	TriggerBossDefeated
	TriggerItemFound
)

type Dialogue struct {
	Speaker string
	Text    string
}

type Reward struct {
	Gold     int
	XP       int
	ItemName string
}

type StoryBeat struct {
	ID        string
	Name      string
	FloorMin  int
	FloorMax  int
	Trigger   TriggerType
	Dialogues []Dialogue
	BossName  string
	Reward    Reward
	Completed bool
}

var (
	registryMu sync.RWMutex
	registry   = map[string]*StoryBeat{}
)

func RegisterStory(s *StoryBeat) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[s.ID] = s
}

func GetStoriesForFloor(floor int) []*StoryBeat {
	registryMu.RLock()
	defer registryMu.RUnlock()
	var out []*StoryBeat
	for _, s := range registry {
		if !s.Completed && floor >= s.FloorMin && floor <= s.FloorMax {
			out = append(out, s)
		}
	}
	return out
}

func MarkCompleted(id string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if s, ok := registry[id]; ok {
		s.Completed = true
	}
}

func ResetStories() {
	registryMu.Lock()
	defer registryMu.Unlock()
	for _, s := range registry {
		s.Completed = false
	}
}

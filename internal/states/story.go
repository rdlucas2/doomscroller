package states

import "github.com/rdlucas2/doomscroller/internal/game"

type StoryState struct {
	Beat *game.StoryBeat
	Line int
	Done bool
}

func NewStoryState(beat *game.StoryBeat) *StoryState { return &StoryState{Beat: beat} }

func (ss *StoryState) Advance() {
	if ss.Line < len(ss.Beat.Dialogues)-1 {
		ss.Line++
	} else {
		ss.Done = true
		game.MarkCompleted(ss.Beat.ID)
	}
}

package states

type IntroState struct {
	elapsed float64
	done    bool
}

func (s *IntroState) Update(dt float64) {
	if s.done {
		return
	}
	s.elapsed += dt
	if s.elapsed >= 3.2 {
		s.done = true
	}
}

func (s *IntroState) Done() bool { return s.done }

// Progress returns 0..1 controlling which intro phase to render.
func (s *IntroState) Progress() float64 {
	if s.elapsed < 0.4 {
		return 0
	}
	if s.elapsed < 1.0 {
		return 0.3 + (s.elapsed-0.4)/0.6*0.2
	}
	if s.elapsed < 2.8 {
		return 0.5 + (s.elapsed-1.0)/1.8*0.5
	}
	return 1.0
}

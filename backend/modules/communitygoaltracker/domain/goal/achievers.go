package goal

// Achievers represensts a map of achievers within a goal where the
// key is an achiever userid and the value contains the achiever.
type Achievers map[string]*Achiever

// UserIDs represent the achiever userids
func (a *Achievers) UserIDs() []string {
	keys := make([]string, len(*a))

	i := 0
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}

// NewAchievers with default state,progress, and messages for given initial achiever.
func NewAchievers(initialAchieverUserID string) (res Achievers) {
	res = make(Achievers)
	state, _ := NewState(nil)
	progress, _ := NewProgress(nil)
	messages := NewMessages(nil)
	res[initialAchieverUserID] = &Achiever{
		State:    &state,
		Progress: &progress,
		Messages: &messages,
	}
	return res

}

package achiever

import "encoding/json"

// Goals represents a map with the goal names as the key.
type Goals map[string]bool

// Names represent the goal names
func (g *Goals) Names() []string {
	keys := make([]string, len(*g))

	i := 0
	for k := range *g {
		keys[i] = k
		i++
	}
	return keys
}

// Remove a goal given the goal name from an Achiever's Goals.
func (g *Goals) Remove(goalName string) (e error) {
	if g == nil {
		return nil
	}
	if _, ok := (*g)[goalName]; !ok {
		return ErrAchieverGoalNotFound
	}
	delete(*g, goalName)
	return nil
}

// ToPrimativeType from Goals type
func (g *Goals) ToPrimativeType() map[string]bool {
	primativeType := make(map[string]bool)

	for k, v := range *g {
		primativeType[k] = v
	}
	return primativeType
}

// ToGoals given a primitive type
func ToGoals(g *map[string]bool) Goals {
	if g == nil {
		return nil
	}
	return Goals(*g)
}

// NewGoals creates new goals. Empty goals by if newGoals is nil.
func NewGoals(newGoals *[]byte) Goals {
	if newGoals != nil {
		g := new(Goals)
		json.Unmarshal(*newGoals, g)
		return *g
	}
	return make(Goals)

}

package goal

// Progress value object
type Progress int

// ToInt from Progress
func (n Progress) ToInt() int {
	return int(n)
}

// NewProgress creates a progress. 0 by defualt. Returns error if progress if
// invalid.
func NewProgress(newProgress *int) (res Progress, e error) {
	if newProgress == nil {
		return Progress(0), nil
	}
	if *newProgress < 0 || *newProgress > 100 {
		return 0, ErrGoalInvalidProgress
	}

	return Progress(*newProgress), nil
}

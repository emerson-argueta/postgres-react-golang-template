package mockgoal

//GoalService is a mock of church.Service
type GoalService struct {
	CreateFn       func(c *church.Church) error
	ReadFn         func(c *church.Church, byEmail bool) (*church.Church, error)
	ReadMultipleFn func(churchids []int64) ([]*church.Church, error)
	UpdateFn       func(c *church.Church, byEmail bool) error
	DeleteFn       func(c *church.Church, byEmail bool) error
}

// Create mocks church.Service.Create
func (s *GoalService) Create(c *church.Church) error {
	return s.CreateFn(c)
}

// Read mocks church.Service.Read
func (s *GoalService) Read(c *church.Church, byEmail bool) (*church.Church, error) {
	return s.ReadFn(c, byEmail)
}

// ReadMultiple mocks church.Service.ReadMultiple
func (s *GoalService) ReadMultiple(churchids []int64) ([]*church.Church, error) {
	return s.ReadMultipleFn(churchids)
}

// Update mocks church.Service.Update
func (s *GoalService) Update(c *church.Church, byEmail bool) error {
	return s.UpdateFn(c, byEmail)
}

// Delete mocks church.Service.Delete
func (s *GoalService) Delete(c *church.Church, byEmail bool) error {
	return s.DeleteFn(c, byEmail)
}

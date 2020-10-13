package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/goal"

	"github.com/lib/pq"
)

var _ goal.Processes = &GoalService{}

// GoalService represents a service for managing a Goal.
type GoalService struct {
	client *Client
}

// CreateGoal if successful. If the goal exists, returns ErrGoalExists.
func (s *GoalService) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	query, e := NewQuery(g)
	if e != nil {
		return nil, e
	}

	goalInsertQuery := query.Create(CommunitygoaltrackerSchema, GoalTable)
	goalInsertQuery = s.client.db.Rebind(goalInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	res = &goal.Goal{}

	e = s.client.db.Get(res, goalInsertQuery, queryParams...)

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return nil, e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return nil, communitygoaltracker.ErrGoalExists
	} else if pqError != nil {
		return nil, pqError
	}

	return res, e
}

// RetrieveGoal by email. If the goal does not exists, returns ErrGoalNotFound.
func (s *GoalService) RetrieveGoal(id int64) (res *goal.Goal, e error) {
	filter := "ID=?"
	queryParam := id

	query, e := NewQuery(&goal.Goal{})
	if e != nil {
		return nil, e
	}

	goalSelectQuery := query.Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.client.db.Rebind(goalSelectQuery)

	res = &goal.Goal{}
	e = s.client.db.Get(res, goalSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, communitygoaltracker.ErrGoalNotFound
	}

	return res, e

}

// UpdateGoal searching by id. If the goal does not exists, returns
// ErrGoalNotFound.
func (s *GoalService) UpdateGoal(g *goal.Goal) (e error) {
	filter := "ID=?"
	queryParam := g.ID

	query, e := NewQuery(g)
	if e != nil {
		return e
	}

	goalUpdateQuery := query.Update(CommunitygoaltrackerSchema, GoalTable, filter)
	goalUpdateQuery = s.client.db.Rebind(goalUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	e = s.client.db.Get(g, goalUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return communitygoaltracker.ErrGoalNotFound
	}

	return e
}

// DeleteGoal searching by id. If the goal does not exists, returns
// ErrGoalNotFound.
func (s *GoalService) DeleteGoal(id int64) (e error) {
	filter := "ID=?"
	queryParam := id

	query, e := NewQuery(&goal.Goal{})
	if e != nil {
		return e
	}
	goalDeleteQuery := query.Delete(CommunitygoaltrackerSchema, GoalTable, filter)
	goalDeleteQuery = s.client.db.Rebind(goalDeleteQuery)

	e = s.client.db.Get(&goal.Goal{}, goalDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return communitygoaltracker.ErrGoalNotFound
	}

	return e
}

// RetrieveGoals searching by ids. Returns ErrGoalNotFound if none of the goals
// are found.
func (s *GoalService) RetrieveGoals(ids []int64) (res []*goal.Goal, e error) {

	query, err := NewQuery(&goal.Goal{})
	if err != nil {
		return nil, err
	}
	var queryParams []interface{}
	for _, elem := range ids {
		queryParams = append(queryParams, elem)
	}

	filter := query.CreateMultipleValueFilter("ID", len(queryParams))

	goalSelectQuery := query.Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.client.db.Rebind(goalSelectQuery)

	if err := s.client.db.Select(&res, goalSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateGoals searching by ids. Returns ErrGoalNotFound if none of the goals
// not found. Return ErrGoalExists if any of the update names conflicts with
// another goal.
func (s *GoalService) UpdateGoals(gg []*goal.Goal) (e error) {
	// TODO
	searchKey := "ID"

	query, e := NewQuery(&goal.Goal{})
	if e != nil {
		return e
	}

	goalUpdateQuery := query.UpdateMultiple(CommunitygoaltrackerSchema, GoalTable, searchKey, len(gg))
	goalUpdateQuery = s.client.db.Rebind(goalUpdateQuery)

	models := make([]interface{}, len(gg))
	for i, g := range gg {
		models[i] = g
	}
	queryParams := MultipleModelValues(models)

	e = s.client.db.Get(gg, goalUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return communitygoaltracker.ErrGoalNotFound
	}

	return e
}

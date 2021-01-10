package persistence

import (
	"context"
	"database/sql"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/database"
	"strconv"
)

var _ repository.GoalRepo = &Goal{}

// Goal holds a a reference to the database client.
type Goal struct {
	Client database.Client
}

// CreateGoal if successful. If the goal exists, returns ErrGoalExists.
func (s *Goal) CreateGoal(g goal.Goal) (e error) {
	goalPersistenceModel := GoalDomainToPersistence(g)

	goalInsertQuery := s.Client.Query().Create(goalPersistenceModel, CommunitygoaltrackerSchema, GoalTable)
	goalInsertQuery = s.Client.DB().Rebind(goalInsertQuery)

	includeNil := true
	queryParams := s.Client.Query().ModelValues(goalPersistenceModel, includeNil)

	ctx := context.WithValue(context.Background(), "aggregateid", strconv.FormatInt(g.GetID(), 10))
	_, e = s.Client.DB().ExecContext(ctx, goalInsertQuery, queryParams...)

	return e
}

// RetrieveGoalByID If the goal does not exists, returns ErrGoalNotFound.
func (s *Goal) RetrieveGoalByID(id int64) (res goal.Goal, e error) {
	filter := "ID=?"
	queryParam := id

	goalSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.Client.DB().Rebind(goalSelectQuery)

	dbRes := &GoalDTO{}
	e = s.Client.DB().Get(dbRes, goalSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return res, goal.ErrGoalNotFound
	}

	return GoalPersistenceToDomain(*dbRes)

}

// RetrieveGoalByName by name. If the goal does not exists, returns ErrGoalNotFound.
func (s *Goal) RetrieveGoalByName(name string) (res goal.Goal, e error) {
	filter := "NAME=?"
	queryParam := name

	goalSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.Client.DB().Rebind(goalSelectQuery)

	dbRes := &GoalDTO{}
	e = s.Client.DB().Get(dbRes, goalSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return res, goal.ErrGoalNotFound
	}

	return GoalPersistenceToDomain(*dbRes)

}

// RetrieveGoalsByIDs searching by ids. Returns ErrGoalNotFound if none of the goals
// are found.
func (s *Goal) RetrieveGoalsByIDs(ids []int64) (res []goal.Goal, e error) {
	var queryParams []interface{}
	for _, elem := range ids {
		queryParams = append(queryParams, elem)
	}

	filter := s.Client.Query().CreateMultipleValueFilter("ID", len(queryParams))

	goalSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.Client.DB().Rebind(goalSelectQuery)

	var dbRes []GoalDTO
	if err := s.Client.DB().Select(&dbRes, goalSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res = make([]goal.Goal, len(dbRes))
	for i, g := range dbRes {
		if res[i], e = GoalPersistenceToDomain(g); e != nil {
			return nil, e
		}
	}
	if len(res) == 0 {
		return nil, goal.ErrGoalNotFound
	}

	return res, nil
}

// RetrieveGoalsByNames searching by ids. Returns ErrGoalNotFound if none of the goals
// are found.
func (s *Goal) RetrieveGoalsByNames(names []string) (res []goal.Goal, e error) {
	var queryParams []interface{}
	for _, elem := range names {
		queryParams = append(queryParams, elem)
	}

	filter := s.Client.Query().CreateMultipleValueFilter("NAME", len(queryParams))

	goalSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, GoalTable, filter)
	goalSelectQuery = s.Client.DB().Rebind(goalSelectQuery)

	var dbRes []GoalDTO
	if err := s.Client.DB().Select(&dbRes, goalSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res = make([]goal.Goal, len(dbRes))
	for i, g := range dbRes {
		if res[i], e = GoalPersistenceToDomain(g); e != nil {
			return nil, e
		}
	}
	if len(res) == 0 {
		return nil, goal.ErrGoalNotFound
	}

	return res, nil
}

// UpdateGoal searching by id. If the goal does not exists, returns
// ErrGoalNotFound.
func (s *Goal) UpdateGoal(g goal.Goal) (e error) {
	filter := "ID=?"

	goalPersistenceModel := GoalDomainToPersistence(g)
	goalUpdateQuery := s.Client.Query().Update(goalPersistenceModel, CommunitygoaltrackerSchema, GoalTable, filter)
	goalUpdateQuery = s.Client.DB().Rebind(goalUpdateQuery)

	includeNil := true
	queryParams := append(s.Client.Query().ModelValues(goalPersistenceModel, includeNil), g.GetID())

	e = s.Client.DB().Get(goalPersistenceModel, goalUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return goal.ErrGoalNotFound
	}

	return e
}

// UpdateGoals searching by ids. Returns ErrGoalNotFound if none of the goals
// not found. Return ErrGoalExists if any of the update names conflicts with
// another goal.
func (s *Goal) UpdateGoals(gg []goal.Goal) (e error) {
	// TODO
	searchKey := "ID"

	goalPersistenceModel := &GoalDTO{}
	goalUpdateQuery := s.Client.Query().UpdateMultiple(goalPersistenceModel, CommunitygoaltrackerSchema, GoalTable, searchKey, len(gg))
	goalUpdateQuery = s.Client.DB().Rebind(goalUpdateQuery)

	models := make([]interface{}, len(gg))
	for i, g := range gg {
		models[i] = GoalDomainToPersistence(g)
	}
	queryParams := s.Client.Query().MultipleModelValues(models)

	goalPersistenceModels := make([]*GoalDTO, len(gg))
	for i, g := range gg {
		goalPersistenceModels[i] = GoalDomainToPersistence(g)
	}
	e = s.Client.DB().Get(goalPersistenceModels, goalUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return goal.ErrGoalNotFound
	}

	return e
}

// DeleteGoal searching by id. If the goal does not exists, returns
// ErrGoalNotFound.
func (s *Goal) DeleteGoal(id int64) (e error) {
	filter := "ID=?"

	goalDeleteQuery := s.Client.Query().Delete(CommunitygoaltrackerSchema, GoalTable, filter)
	goalDeleteQuery = s.Client.DB().Rebind(goalDeleteQuery)

	goalPersistenceModel := &GoalDTO{}
	queryParam := id
	e = s.Client.DB().Get(goalPersistenceModel, goalDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return goal.ErrGoalNotFound
	}

	return e
}

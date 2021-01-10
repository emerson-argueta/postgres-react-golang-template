package persistence

import (
	"context"
	"database/sql"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/database"
)

var _ repository.AchieverRepo = &Achiever{}

// Achiever holds reference to database client.
type Achiever struct {
	Client database.Client
}

// CreateAchiever if successful. If the achiever
// exists, returns ErrAchieverExists.
func (s *Achiever) CreateAchiever(a achiever.Achiever) (e error) {
	achieverPersistenceModel := AchieverDomainToPersistence(a)

	achieverInsertQuery := s.Client.Query().Create(achieverPersistenceModel, CommunitygoaltrackerSchema, AchieverTable)
	achieverInsertQuery = s.Client.DB().Rebind(achieverInsertQuery)

	includeNil := true
	queryParams := s.Client.Query().ModelValues(achieverPersistenceModel, includeNil)

	ctx := context.WithValue(context.Background(), "aggregateid", a.GetID())
	_, e = s.Client.DB().ExecContext(ctx, achieverInsertQuery, queryParams...)

	return e
}

// RetrieveAchiever by id. If the achiever does not exists,
// returns ErrAchieverNotFound.
func (s *Achiever) RetrieveAchiever(id string) (res achiever.Achiever, e error) {
	filter := "ID=?"
	queryParam := id

	achieverSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverSelectQuery = s.Client.DB().Rebind(achieverSelectQuery)

	dbRes := &AchieverDTO{}
	e = s.Client.DB().Get(dbRes, achieverSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, achiever.ErrAchieverNotFound
	}

	return AchieverPersistenceToDomain(*dbRes)

}

// RetrieveAchieverByUserID and If the achiever does not exists,
// returns ErrAchieverNotFound.
func (s *Achiever) RetrieveAchieverByUserID(id string) (res achiever.Achiever, e error) {
	filter := "USERID=?"
	queryParam := id

	achieverSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverSelectQuery = s.Client.DB().Rebind(achieverSelectQuery)

	dbRes := &AchieverDTO{}
	e = s.Client.DB().Get(dbRes, achieverSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, achiever.ErrAchieverNotFound
	}

	return AchieverPersistenceToDomain(*dbRes)

}

// RetrieveAchieversByUserIDs searching by user ids. Returns ErrAchiverNotFound if none of the achievers
// are found.
func (s *Achiever) RetrieveAchieversByUserIDs(userIDs []string) (res []achiever.Achiever, e error) {
	var queryParams []interface{}
	for _, elem := range userIDs {
		queryParams = append(queryParams, elem)
	}

	filter := s.Client.Query().CreateMultipleValueFilter("USERID", len(queryParams))

	achieverSelectQuery := s.Client.Query().Read(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverSelectQuery = s.Client.DB().Rebind(achieverSelectQuery)

	dbRes := []AchieverDTO{}
	if err := s.Client.DB().Select(&dbRes, achieverSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res = make([]achiever.Achiever, len(dbRes))
	for i, a := range dbRes {
		if res[i], e = AchieverPersistenceToDomain(a); e != nil {
			return nil, e
		}
	}
	return res, nil
}

// UpdateAchiever searching by id. If the achiever does not exists, returns
// ErrAchieverNotFound.
func (s *Achiever) UpdateAchiever(a achiever.Achiever) (e error) {
	filter := "ID=?"

	achieverPersistenceModel := AchieverDomainToPersistence(a)
	achieverUpdateQuery := s.Client.Query().Update(achieverPersistenceModel, CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverUpdateQuery = s.Client.DB().Rebind(achieverUpdateQuery)

	includeNil := true
	queryParams := append(s.Client.Query().ModelValues(achieverPersistenceModel, includeNil), a.GetID())

	e = s.Client.DB().Get(achieverPersistenceModel, achieverUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return achiever.ErrAchieverNotFound
	}

	return e
}

// DeleteAchiever searching by id. If the achiever does not exists, returns
// ErrAchieverNotFound.
func (s *Achiever) DeleteAchiever(id string) (e error) {
	filter := "ID=?"

	achieverDeleteQuery := s.Client.Query().Delete(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverDeleteQuery = s.Client.DB().Rebind(achieverDeleteQuery)

	achieverPersistenceModel := &AchieverDTO{}
	queryParam := id
	e = s.Client.DB().Get(achieverPersistenceModel, achieverDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return achiever.ErrAchieverNotFound
	}

	return e
}

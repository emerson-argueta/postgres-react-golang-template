package postgres

import (
	"database/sql"

	"emersonargueta/m/v1/domain"
	"emersonargueta/m/v1/domain/transaction"
)

const transactionTable = "transaction"

var _ transaction.Service = &Transaction{}

// Transaction represents a service for managing transactions.
type Transaction struct {
	client *Client
}

// CreateManagementSession opens a session to start a series of actions taken to
// manage a transaction. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (s *Transaction) CreateManagementSession() error {
	return s.client.createManagementSession()
}

// EndManagementSession ends the session created to execute a series of actions
// taken to manage a transaction. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (s *Transaction) EndManagementSession() error {
	return s.client.endManagementSession()
}

// Create a new transaction.
func (s *Transaction) Create(txn *transaction.Transaction) (e error) {
	schema := s.client.config.Database.Schema

	query, e := NewQuery(txn)
	if e != nil {
		return e
	}

	transactionInsertQuery := query.Create(schema, transactionTable)
	transactionInsertQuery = s.client.db.Rebind(transactionInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(txn, transactionInsertQuery, queryParams...)
	} else {
		e = s.client.db.Get(txn, transactionInsertQuery, queryParams...)
	}
	if e != nil {
		return e
	}

	return nil
}

// Read transaction for an authenticated administrator.
// The number of read transactions can be limited if limit paramater is not nil.
func (s *Transaction) Read(donatorID int64, churchID int64, timeRange *transaction.TimeRange, limit *int64) (res []*transaction.Transaction, e error) {
	schema := s.client.config.Database.Schema
	filter := "DONATORID=? AND CHURCHID=? LIMIT ?"
	queryParams := []interface{}{
		donatorID,
		churchID,
		limit,
	}

	if timeRange != nil {
		filter = "DONATORID=? AND CHURCHID=? AND CREATEDAT BETWEEN ? AND ? LIMIT ?"
		queryParams = []interface{}{
			donatorID,
			churchID,
			timeRange.Lower,
			timeRange.Upper,
			limit,
		}

	}

	t := &transaction.Transaction{}
	query, err := NewQuery(t)
	if err != nil {
		return nil, err
	}

	transactionSelectQuery := query.Read(schema, transactionTable, filter)
	transactionSelectQuery = s.client.db.Rebind(transactionSelectQuery)

	if err := s.client.db.Select(
		&res,
		transactionSelectQuery,
		queryParams...,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// ReadWithFilter transactions, searching by non empty fields of filterTransaction model.
// The number of read transactions can be limited if limit paramater is not nil.
func (s *Transaction) ReadWithFilter(donatorID int64, churchID int64, timeRange *transaction.TimeRange, limit *int64, transactionFilter *domain.Filter) (res []*transaction.Transaction, e error) {
	schema := s.client.config.Database.Schema

	filter := "DONATORID=? AND CHURCHID=? LIMIT ?"
	queryParams := []interface{}{
		donatorID,
		churchID,
		limit,
	}

	if timeRange != nil {
		filter = "DONATORID=? AND CHURCHID=? AND CREATEDAT BETWEEN ? AND ? LIMIT ?"
		queryParams = []interface{}{
			donatorID,
			churchID,
			timeRange.Lower,
			timeRange.Upper,
			limit,
		}

	}

	t := &transaction.Transaction{}
	query, err := NewQuery(t)
	if err != nil {
		return nil, err
	}

	if transactionFilter != nil && len(*transactionFilter) != 0 {
		f, qp := query.CreateFilter(transactionFilter)
		filter = f + " AND " + filter

		queryParams = append(qp, queryParams...)
	}

	transactionSelectQuery := query.Read(schema, transactionTable, filter)
	transactionSelectQuery = s.client.db.Rebind(transactionSelectQuery)

	if err := s.client.db.Select(
		&res,
		transactionSelectQuery,
		queryParams...,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

//ReadMultiple donator's donations
// The number of read transactions can be limited if limit paramater is not nil.
func (s *Transaction) ReadMultiple(donatorids []int64, churchID int64, timeRange *transaction.TimeRange, limit *int64) (res []*transaction.Transaction, e error) {
	schema := s.client.config.Database.Schema

	t := &transaction.Transaction{}
	query, err := NewQuery(t)
	if err != nil {
		return nil, err
	}

	var queryParams []interface{}
	for _, elem := range donatorids {
		queryParams = append(queryParams, elem)
	}
	filter := query.CreateMultipleValueFilter("DONATORID", len(queryParams))

	if timeRange == nil {
		filter += " AND CHURCHID=? LIMIT ?"

		queryParams = append(
			queryParams,
			churchID,
			limit,
		)

	} else {
		filter += " AND CHURCHID=? AND CREATEDAT BETWEEN ? AND ? LIMIT ?"

		queryParams = append(
			queryParams,
			churchID,
			timeRange.Lower,
			timeRange.Upper,
			limit,
		)

	}

	transactionSelectQuery := query.Read(schema, transactionTable, filter)
	transactionSelectQuery = s.client.db.Rebind(transactionSelectQuery)

	if err := s.client.db.Select(
		&res,
		transactionSelectQuery,
		queryParams...,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// ReadMultipleWithFilter transactions, searching by non empty fields of filterTransaction model.
// The number of read transactions can be limited if limit paramater is not nil.
func (s *Transaction) ReadMultipleWithFilter(donatorids []int64, churchID int64, timeRange *transaction.TimeRange, limit *int64, transactionFilter *domain.Filter) (res []*transaction.Transaction, e error) {
	schema := s.client.config.Database.Schema

	t := &transaction.Transaction{}
	query, err := NewQuery(t)
	if err != nil {
		return nil, err
	}

	var queryParams []interface{}
	for _, elem := range donatorids {
		queryParams = append(queryParams, elem)
	}
	filter := query.CreateMultipleValueFilter("DONATORID", len(queryParams))

	if timeRange == nil {
		filter += " AND CHURCHID=? LIMIT ?"

		queryParams = append(
			queryParams,
			churchID,
			limit,
		)

	} else {
		filter += " AND CHURCHID=? AND CREATEDAT BETWEEN ? AND ? LIMIT ?"

		queryParams = append(
			queryParams,
			churchID,
			timeRange.Lower,
			timeRange.Upper,
			limit,
		)

	}
	if transactionFilter != nil && len(*transactionFilter) != 0 {
		f, qp := query.CreateFilter(transactionFilter)

		if len(f) != 0 {
			filter = f + " AND " + filter
			queryParams = append(qp, queryParams...)
		}

	}

	transactionSelectQuery := query.Read(schema, transactionTable, filter)
	transactionSelectQuery = s.client.db.Rebind(transactionSelectQuery)

	if err := s.client.db.Select(
		&res,
		transactionSelectQuery,
		queryParams...,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete all transactions for an administrator. This should only be done if an
// administrator is deleted.
func (s *Transaction) Delete(donatorID int64, churchID int64) (e error) {
	schema := s.client.config.Database.Schema
	filter := "DONATORID=? AND CHURCHID=?"

	t := &transaction.Transaction{}
	query, e := NewQuery(t)
	if e != nil {
		return e
	}

	transactionDeleteQuery := query.Delete(schema, transactionTable, filter)
	transactionDeleteQuery = s.client.db.Rebind(transactionDeleteQuery)

	tt := []*transaction.Transaction{}

	if s.client.transaction != nil {
		e = s.client.transaction.Select(&tt, transactionDeleteQuery, donatorID, churchID)
	} else {
		e = s.client.db.Select(&tt, transactionDeleteQuery, donatorID, churchID)
	}
	if e != nil && e == sql.ErrNoRows {
		return transaction.ErrTransactionNotFound
	} else if e != nil {
		return e
	}

	return nil
}

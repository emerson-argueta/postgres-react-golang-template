package mocktransaction

import (
	"trustdonations.org/m/v2/domain/transaction"
)

//TransactionService is a mock of transaction.Service
type TransactionService struct {
	CreateManagementSessionFn func() error
	EndManagementSessionFn    func() error
	CreateFn                  func(transaction *transaction.Transaction) error
	ReadFn                    func(donatorID int64, churchID int64, timeRange *transaction.TimeRange, limit *int64) ([]*transaction.Transaction, error)
	DeleteFn                  func(donatorID int64, churchID int64) error
}

// Create mocks transaction.Service.Create.
func (s *TransactionService) Create(transaction *transaction.Transaction) error {
	return s.CreateFn(transaction)
}

// Read mocks transaction.Service.Read.
func (s *TransactionService) Read(donatorID int64, churchID int64, timeRange *transaction.TimeRange, limit *int64) ([]*transaction.Transaction, error) {
	return s.ReadFn(donatorID, churchID, timeRange, limit)
}

// Delete mocks transaction.Service.Delete.
func (s *TransactionService) Delete(donatorID int64, churchID int64) error {
	return s.DeleteFn(donatorID, churchID)
}

// CreateManagementSession mocks transaction.Service.CreateManagementSession
func (s *TransactionService) CreateManagementSession() error {
	return s.CreateManagementSessionFn()
}

// EndManagementSession mocks transaction.Service.EndManagementSession
func (s *TransactionService) EndManagementSession() error {
	return s.EndManagementSessionFn()
}

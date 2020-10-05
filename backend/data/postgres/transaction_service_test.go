package postgres_test

import (
	"reflect"
	"testing"

	"emersonargueta/m/v1/domain/transaction"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.TransactionService()

	churchid_t := int64(1)
	donatorid_t := int64(1)
	amount_t := float64(11.11)

	createdAt := c.Now()
	transactionType := transaction.Credit

	donationType := transaction.Cash
	currency := "usd"
	account := transaction.Donator
	category := "general"
	txn := transaction.Transaction{
		ChurchID:  &churchid_t,
		DonatorID: &donatorid_t,
		Amount:    &amount_t,
		CreatedAt: &createdAt,
		Type:      &transactionType,
		Donation:  &transaction.Donation{Type: &donationType, Currency: &currency, Account: &account, Category: &category},
	}

	// Create two transaction.
	if err := s.Create(&txn); err != nil {
		t.Fatal(err)
	} else if err = s.Create(&txn); err != nil {
		t.Fatal(err)
	}

	// Retrieve transactions by ADMINID and compare.
	if other, err := s.Read(*txn.DonatorID, *txn.ChurchID, &transaction.TimeRange{Lower: &createdAt, Upper: &createdAt}, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&txn, other[0]) {
		t.Fatalf("unexpected transaction: %#v", other[0])
	} else if !reflect.DeepEqual(&txn, other[1]) {
		t.Fatalf("unexpected transaction: %#v", other[1])
	}

	// Clean up database.
	if err := s.Delete(*txn.DonatorID, *txn.ChurchID); err != nil {
		t.Fatal(err)
	}
}

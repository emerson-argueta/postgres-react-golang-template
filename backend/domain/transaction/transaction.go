package transaction

import (
	"time"

	"trustdonations.org/m/v2/domain"
)

// TimePeriod can be day, week, month, or year
type TimePeriod string

func (t TimePeriod) String() string {
	return string(t)
}

// ParseTimePeriod from string
func ParseTimePeriod(timePeriod string) (res TimePeriod, e error) {

	switch timePeriod {
	case Day.String():
		return Day, nil
	case Week.String():
		return Week, nil
	case Month.String():
		return Month, nil
	case Year.String():
		return Year, nil
	default:
		return res, ErrTransactionInternal
	}
}

// Day is a basic unit of TimePeriod
const Day = TimePeriod("day")

// Week is a basic unit of TimePeriod
const Week = TimePeriod("week")

// Month is a basic unit of TimePeriod
const Month = TimePeriod("month")

// Year is a basic unit of TimePeriod
const Year = TimePeriod("year")

// Type can be either credit or debit.
type Type string

func (t Type) String() string {
	return string(t)
}

// Credit transactions are tied to the donator
// account when a donation is made by the donator.
const Credit = Type("credit")

// Debit transactions are tied to the church account
// when a donation is made by the donator.
const Debit = Type("debit")

// DonationType can be online,cash,check,cryptocurrency,etc
type DonationType string

func (dt DonationType) String() string {
	return string(dt)
}

// Online donation type
const Online = DonationType("online")

// Cash donation type
const Cash = DonationType("cash")

// Check donation type
const Check = DonationType("check")

// Account can be a church or donator account.
type Account string

// Church account will hold the debit transactions created by donations.
const Church = Account("church")

// Donator account will hold the credit transactions created by donations.
const Donator = Account("donator")

// DonatorDonationsMap is a map whose keys are a string value that can
// be used to group common donations and the values are an array of donations.
type DonatorDonationsMap map[string][]*Donation

// Values of DonatorDonationsMap
func (d DonatorDonationsMap) Values() (res []interface{}) {
	for _, value := range d {
		res = append(res, value)
	}
	return res
}

// Keys to group common donations
func (d DonatorDonationsMap) Keys() []string {
	keys := make([]string, len(d))

	i := 0
	for k := range d {
		keys[i] = k
		i++
	}
	return keys
}

// DonationsMap is a map whose keys are donatorids and values are maps of the
// donator's donations.
type DonationsMap map[int64]DonatorDonationsMap

// Values of DonatorDonationsMap
func (d DonationsMap) Values() (res []interface{}) {
	for _, value := range d {
		res = append(res, value)
	}
	return res
}

// Keys to group common donations
func (d DonationsMap) Keys() []int64 {
	keys := make([]int64, len(d))

	i := 0
	for k := range d {
		keys[i] = k
		i++
	}
	return keys
}

// DonationsSum is a map whose keys are a string value that can be used to group
// common donations and the values the sum of the common donations.
type DonationsSum map[string]float64

// DonationsSumMap is a map whose keys are donation categories and values are
// maps of category donations sums.
type DonationsSumMap map[string]DonationsSum

// DonationReport represents a information used to create a church donation
// report or a single donator's donation report.
type DonationReport struct {
	DonatorIDs         []int64         `json:"donatorid,omitempty"`
	ChurchID           *int64          `json:"churchid,omitempty"`
	SumFilter          *SumFilter      `json:"sumfilter,omitempty"`
	TimeRange          *TimeRange      `json:"timerange,omitempty"`
	DonationCategories []string        `json:"donationcategories,omitempty"`
	Donations          DonationsMap    `json:"donations,omitempty"`
	DonationsSum       DonationsSumMap `json:"donationssum,omitempty"`
}

// Donation represents information about a transaction which is created from a donation.
type Donation struct {
	DonatorID *int64        `json:"donatorid,omitempty"`
	ChurchID  *int64        `json:"churchid,omitempty"`
	Amount    *float64      `json:"amount,omitempty"`
	Type      *DonationType `json:"type,omitempty"`
	Currency  *string       `json:"currency,omitempty"`
	Account   *Account      `json:"account,omitempty"`
	Category  *string       `json:"category,omitempty"`
	Details   *string       `json:"details,omitempty"`
	Date      *string       `json:"date,omitempty"`
}

// SumFilter represents a unit of time period and a multiplier for that unit of time period
type SumFilter struct {
	TimePeriod *TimePeriod `json:"timeperiod,omitempty"`
	Multiplier *int64      `json:"multiplier,omitempty"`
}

// TimeRange represents an upper and lower range of time for retrieving transactions.
type TimeRange struct {
	Upper *time.Time `json:"upper,omitempty"`
	Lower *time.Time `json:"lower,omitempty"`
}

// Transaction represents a transaction for a donator or church account, which
// is created when a donation is made by a donator.
type Transaction struct {
	DonatorID *int64     `db:"donatorid"`
	ChurchID  *int64     `db:"churchid"`
	Amount    *float64   `db:"amount" json:"amount"`
	Type      *Type      `db:"type" json:"type"`
	Donation  *Donation  `db:"donation" json:"donation"`
	CreatedAt *time.Time `db:"createdat" json:"createdat"`
	Updatedat *time.Time `db:"updatedat" json:"updatedat"`
}

// Client creates a connection to a service. In this case, a transaction service.
type Client interface {
	Service() Service
}

// Service provides functions managing transactions.
type Service interface {
	CreateManagementSession() error
	EndManagementSession() error
	Create(*Transaction) error
	// Read transaction for an authenticated administrator.
	// The number of read transactions can be limited if limit paramater is not nil.
	Read(donatorID int64, churchID int64, timeRange *TimeRange, limit *int64) ([]*Transaction, error)
	// ReadWithFilter transactions, searching by key and values of filterTransaction map.
	// The number of read transactions can be limited if limit paramater is not nil.
	ReadWithFilter(donatorID int64, churchID int64, timeRange *TimeRange, limit *int64, transactionFilter *domain.Filter) ([]*Transaction, error)
	// ReadMultiple transactions for an authenticated administrator.
	// The number of read transactions can be limited if limit paramater is not nil.
	ReadMultiple(donatorids []int64, churchID int64, timeRange *TimeRange, limit *int64) ([]*Transaction, error)
	// ReadMultipleWithFilter transactions, where the transactionFilter is a map whose keys are the fields to be filtered and the values
	// The number of read transactions can be limited if limit paramater is not nil.
	ReadMultipleWithFilter(donatorids []int64, churchID int64, timeRange *TimeRange, limit *int64, transactionFilter *domain.Filter) ([]*Transaction, error)
	// Delete all transactions for an administrator.
	// This should only be done if an administrator is deleted.
	Delete(donatorID int64, churchID int64) error
}

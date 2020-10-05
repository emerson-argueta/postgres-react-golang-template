package postgres

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/domain"
	"emersonargueta/m/v1/domain/administrator"
	"emersonargueta/m/v1/domain/church"
	"emersonargueta/m/v1/domain/donator"
	"emersonargueta/m/v1/domain/transaction"
	"emersonargueta/m/v1/user"

	"github.com/jmoiron/sqlx"
	// using postgres implementation of sqlx
	_ "github.com/lib/pq"
)

// Client represents a client to the underlying PostgreSQL database.
type Client struct {

	// Returns the current time.
	Now func() time.Time

	config *config.Config

	Services Services

	db          *sqlx.DB
	transaction *sqlx.Tx
}

// Services represents the services that the postgres service provides
type Services struct {
	Administrator Administrator
	Transaction   Transaction
	User          User
	Church        Church
	Donator       Donator
}

// NewClient function
func NewClient() *Client {
	c := &Client{Now: time.Now, transaction: nil}
	c.Services.Administrator.client = c
	c.Services.Transaction.client = c
	c.Services.User.client = c
	c.Services.Church.client = c
	c.Services.Donator.client = c

	// get configuration stucts via .env file
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	c.config = config

	return c
}

// Open and initializes the PostgreSQL database.
func (c *Client) Open() error {

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s search_path=%s sslmode=disable",
		c.config.Database.Host, c.config.Database.Port, c.config.Database.User,
		c.config.Database.Password, c.config.Database.DB, c.config.Database.Schema,
	)

	db, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}
	c.db = db

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return nil
}

// Close closes then underlying postgres database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// AdministratorService returns the admin service associated with the client.
func (c *Client) AdministratorService() administrator.Service { return &c.Services.Administrator }

// TransactionService returns the transaction service associated with the client.
func (c *Client) TransactionService() transaction.Service { return &c.Services.Transaction }

// UserService returns the user service associated with the client.
func (c *Client) UserService() user.Service { return &c.Services.User }

// ChurchService returns the church service associated with the client.
func (c *Client) ChurchService() church.Service { return &c.Services.Church }

// DonatorService returns the donator service associated with the client.
func (c *Client) DonatorService() donator.Service { return &c.Services.Donator }

// Query represents an object to create generic CRUD queries based on a model.
type Query struct {
	model interface{}
}

// NewQuery returns a Query object with a model.
func NewQuery(model interface{}) (*Query, error) {
	if model == nil {
		return nil, fmt.Errorf("model is required for new query")
	}
	q := &Query{}
	q.model = model
	return q, nil
}

// Create returns a generic create query for a model.
func (q *Query) Create(schema string, table string) string {
	fields := structFieldsToStringSlice(q.model)

	queryFields := "("
	queryCreateParameters := "VALUES ("

	for _, field := range fields {
		queryFields += field + ","
		queryCreateParameters += "?,"
	}

	queryFields = strings.TrimSuffix(queryFields, ",") + ")"
	queryCreateParameters = strings.TrimSuffix(queryCreateParameters, ",") + ")"

	return fmt.Sprintf("INSERT INTO %s.%s %s %s RETURNING *", schema, table, queryFields, queryCreateParameters)
}

// Read returns a generic read query for a model.
func (q *Query) Read(schema string, table string, filter string) string {
	return fmt.Sprintf("SELECT * FROM %s.%s WHERE %s", schema, table, filter)
}

// Update returns a generic update query for a model.
func (q *Query) Update(schema string, table string, filter string) string {
	fields := structFieldsToStringSlice(q.model)

	queryFields := "("
	queryUpdateParameters := "("

	for _, field := range fields {
		queryFields += field + ","
		queryUpdateParameters += "COALESCE(?," + field + "),"
	}

	queryFields = strings.TrimSuffix(queryFields, ",") + ")"
	queryUpdateParameters = strings.TrimSuffix(queryUpdateParameters, ",") + ")"

	return fmt.Sprintf("UPDATE %s.%s SET %s = %s WHERE %s RETURNING *", schema, table, queryFields, queryUpdateParameters, filter)
}

// Delete returns a generic delete query for a model.
func (q *Query) Delete(schema string, table string, filter string) string {
	return fmt.Sprintf("DELETE FROM %s.%s WHERE %s RETURNING *", schema, table, filter)
}

// ModelValues returns the all values for the Query's model.
func (q *Query) ModelValues(includeNil bool) []interface{} {
	var vv []interface{}
	modelType := reflect.TypeOf(q.model).Elem()
	modelValue := reflect.ValueOf(q.model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		if _, ignore := modelType.Field(i).Tag.Lookup("dbignoreinsert"); ignore {
			continue
		}
		if _, dbTag := modelType.Field(i).Tag.Lookup("db"); dbTag {
			if includeNil {
				fieldValue := modelValue.Field(i).Interface()
				vv = append(vv, fieldValue)
			} else {
				fieldValue := modelValue.Field(i).Interface()
				if !reflect.ValueOf(fieldValue).IsNil() {
					vv = append(vv, fieldValue)
				}
			}
		}

	}

	return vv

}

// CreateFilter for a query from a domain.Filter.
func (q *Query) CreateFilter(filter *domain.Filter) (resQuery string, resQueryParams []interface{}) {

	for field, values := range *filter {
		// if field contains dot then split string on dot and form filter value substr1-->'substr2' for json field types
		if strings.Contains(field, ".") {
			// TODO: incorporate nested json fields
			field = strings.Split(field, ".")[0] + "->>'" + strings.Split(field, ".")[1] + "'"
		}

		if len(values) <= 0 {
			continue
		}
		prevComparatorOperator := ""
		for _, value := range values {
			logicalOperator := "AND "
			if len(prevComparatorOperator) == 0 || prevComparatorOperator == value.ComparatorOperator {
				logicalOperator = "OR "
			}

			rhs := field + " " + value.ComparatorOperator
			if reflect.TypeOf(value.Value).Kind() == reflect.Array || reflect.TypeOf(value.Value).Kind() == reflect.Slice {
				values := value.Value.([]interface{})
				rhs = q.CreateMultipleValueFilter(field, len(values))
				resQueryParams = append(resQueryParams, values...)
				resQuery += rhs
			} else {
				resQueryParams = append(resQueryParams, value.Value)
				resQuery += rhs + " ? "
			}

			resQuery += logicalOperator

			prevComparatorOperator = value.ComparatorOperator
		}
	}
	resQuery = strings.TrimSuffix(resQuery, "OR ")
	resQuery = strings.TrimSuffix(resQuery, "AND ")

	return resQuery, resQueryParams
}

// CreateFilterFromModel creates a filter from fields of model that are not nil
func (q *Query) CreateFilterFromModel() string {
	filter := ""
	modelType := reflect.TypeOf(q.model).Elem()
	modelValue := reflect.ValueOf(q.model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		if _, dbTag := modelType.Field(i).Tag.Lookup("db"); dbTag {
			fieldValue := modelValue.Field(i).Interface()

			if !reflect.ValueOf(fieldValue).IsNil() {
				fieldName := modelType.Field(i).Tag.Get("db")
				filter += fieldName + "=? AND "
			}

		}
	}
	filter = strings.TrimSuffix(filter, "AND ")

	return filter
}

// CreateMultipleValueFilter creates a sql filter for multiple values from then length af an array of elements
func (q *Query) CreateMultipleValueFilter(field string, valuesLen int) string {
	if valuesLen <= 0 {
		return ""
	}

	filter := field + " IN ("
	for i := 0; i < valuesLen; i++ {
		filter += "?, "
	}

	filter = strings.TrimSuffix(filter, ", ") + ")"

	return filter
}

func structFieldsToStringSlice(model interface{}) []string {
	var ss []string
	modelType := reflect.TypeOf(model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		if _, ignore := modelType.Field(i).Tag.Lookup("dbignoreinsert"); ignore {
			continue
		} else if _, dbTag := modelType.Field(i).Tag.Lookup("db"); dbTag {
			fieldName := modelType.Field(i).Tag.Get("db")
			ss = append(ss, fieldName)
		}
	}

	return ss
}

// createManagementSession opens a session to start a series of actions taken to
// manage an church. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (c *Client) createManagementSession() error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	c.transaction = tx

	return nil
}

// endManagementSession ends the session created to execute a series of actions
// taken. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (c *Client) endManagementSession() error {
	if c.transaction == nil {
		return nil
	}
	if err := c.transaction.Commit(); err != nil {
		if err := c.transaction.Rollback(); err != nil {
			return err
		}

		return err
	}
	c.transaction = nil

	return nil
}

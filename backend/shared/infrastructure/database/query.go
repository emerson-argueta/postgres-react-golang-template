package database

// Filter is a map used to filter domain models by their field values. The keys represent the field names.
// The values represent an array of FilterValue structs.
type Filter map[string][]FilterValue

// FilterValue holds a ComparatorOperator that can be used on each Value for filtering.
type FilterValue struct {
	ComparatorOperator string
	Value              interface{}
}

// Values of Filter map
func (f Filter) Values() (res []interface{}) {
	for _, value := range f {
		res = append(res, value)
	}
	return res
}

// Query provides functions fot database to implement crud
type Query interface {
	// Create query string is a sql statement to create a new entry.
	Create(model interface{}, schema string, table string) string
	// Read query string is a sql statement to read an existing record.
	Read(schema string, table string, filter string) string
	// Update query string is a sql statement to update an existing record.
	Update(model interface{}, schema string, table string, filter string) string
	// UpdateUpdateMultiple query string is a sql statement to update multiple existing records.
	UpdateMultiple(model interface{}, schema string, table string, searchKey string, numValues int) string
	// Delete query string is a sql statement to delete an existing record.
	Delete(schema string, table string, filter string) string
	// CreateFilter for a query
	CreateFilter(model interface{}, filter *Filter) (query string, queryParams []interface{})
	CreateMultipleValueFilter(field string, numVals int) string
	// ModelValues returned, all values including nil returned if includeNil is true
	ModelValues(model interface{}, includeNil bool) []interface{}
	// MultipleModelValues returns the values of each field for all all models.
	MultipleModelValues(models []interface{}) []interface{}
	// CreateFilterFromModel
	CreateFilterFromModel(model interface{}) string
}

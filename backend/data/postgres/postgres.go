package postgres

import (
	"fmt"
	"reflect"
	"strings"
)

// Query represents an object to create generic CRUD queries based on a model.
type Query struct {
	model interface{}
}

// Filter is a map used to filter domain models by their field values. The
// fields are specified in the key of the map. The value or values are specified
// by the map's FilterValue struct. The operator used to filter by field and
// value is specified by the map's FilterValue struct.
type Filter map[string][]FilterValue

// FilterValue hold the value or values and comparator operator used on each value or values to filter transactions.
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

// CreateFilter for a query from a Filter.
func (q *Query) CreateFilter(filter *Filter) (resQuery string, resQueryParams []interface{}) {

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

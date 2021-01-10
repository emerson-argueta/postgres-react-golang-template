package postgres

import (
	"emersonargueta/m/v1/shared/infrastructure/database"
	"fmt"
	"reflect"
	"strings"
)

var _ database.Query = &postgresQuery{}

type postgresQuery struct {
}

// NewQuery with query methods
func NewQuery() database.Query {
	return &postgresQuery{}
}

// Create returns a generic create query for a model.
func (q *postgresQuery) Create(model interface{}, schema string, table string) string {

	fields := structFieldsToStringSlice(model)

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
func (q *postgresQuery) Read(schema string, table string, filter string) string {

	return fmt.Sprintf("SELECT * FROM %s.%s WHERE %s", schema, table, filter)
}

// Update returns a generic update query for a model.
func (q *postgresQuery) Update(model interface{}, schema string, table string, filter string) string {

	fields := structFieldsToStringSlice(model)

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

// UpdateMultiple returns a generic multiple update query for a model.
func (q *postgresQuery) UpdateMultiple(model interface{}, schema string, table string, searchKey string, numValues int) string {
	fields := updateFields(model)

	queryFields := "("
	queryUpdateParameters := "("
	for _, field := range fields {
		queryFields += fmt.Sprintf("%s = COALESCE(VALUELIST.%s,T.%s)", field, field, field)
		queryUpdateParameters += "?,"
	}
	queryFields = strings.TrimSuffix(queryFields, ",") + ")"
	queryUpdateParameters = strings.TrimSuffix(queryUpdateParameters, ",") + ")"

	values := ""
	for i := 0; i < numValues; i++ {
		values += queryUpdateParameters + ","
	}
	values = strings.TrimSuffix(queryUpdateParameters, ",")

	filter := fmt.Sprintf("VALUELIST.%s=T.%s", searchKey, searchKey)

	return fmt.Sprintf("UPDATE %s.%s as T SET %s FROM ( VALUES %s ) WHERE %s RETURNING *", schema, table, queryFields, values, filter)
}

// Delete returns a generic delete query for a model.
func (q *postgresQuery) Delete(schema string, table string, filter string) string {
	return fmt.Sprintf("DELETE FROM %s.%s WHERE %s RETURNING *", schema, table, filter)
}

// ModelValues returns the all values for the Query's model.
func (q *postgresQuery) ModelValues(model interface{}, includeNil bool) []interface{} {

	var vv []interface{}
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

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

// MultipleModelValues returns the values of each field for all all models.
func (q *postgresQuery) MultipleModelValues(models []interface{}) []interface{} {
	var vv []interface{}

	for i := 0; i < len(models); i++ {
		modelType := reflect.TypeOf(models[i]).Elem()
		modelValue := reflect.ValueOf(models[i]).Elem()

		for i := 0; i < modelType.NumField(); i++ {
			if _, dbTag := modelType.Field(i).Tag.Lookup("db"); dbTag {
				fieldValue := modelValue.Field(i).Interface()
				vv = append(vv, fieldValue)
			}
		}
	}

	return vv

}

// CreateFilter for a query from a Filter.
func (q *postgresQuery) CreateFilter(model interface{}, filter *database.Filter) (resQuery string, resQueryParams []interface{}) {
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

// CreateMultipleValueFilter creates a sql filter for multiple values from then length af an array of elements
func (q *postgresQuery) CreateMultipleValueFilter(field string, valuesLen int) string {

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

// CreateFilterFromModel creates a filter from fields of model that are not nil
func (q *postgresQuery) CreateFilterFromModel(model interface{}) string {

	filter := ""
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

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
func updateFields(model interface{}) []string {
	var ss []string
	modelType := reflect.TypeOf(model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		if _, dbTag := modelType.Field(i).Tag.Lookup("db"); dbTag {
			fieldName := modelType.Field(i).Tag.Get("db")
			ss = append(ss, fieldName)
		}
	}

	return ss
}

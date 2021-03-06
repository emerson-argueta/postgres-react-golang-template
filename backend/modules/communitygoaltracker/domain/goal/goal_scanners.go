package goal

// Value marshalls Achievers.
// func (a Achievers) Value() (driver.Value, error) {
// 	if a == nil {
// 		return nil, nil
// 	}

// 	return json.Marshal(a)
// }

// Scan converts raw JSON ([]byte) to Achievers.
// func (a *Achievers) Scan(value interface{}) (e error) {
// 	if value == nil {
// 		return nil
// 	}

// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	m := make(map[string]interface{})
// 	json.Unmarshal(b, &m)
// 	for key, value := range m {
// 		uuid := key
// 		achieverGoal := value.(map[string]interface{})

// 		if state, e := ToState(achieverGoal["state"].(string)); e != nil {
// 			return e
// 		} else {
// 			m[uuid].(map[string]interface{})["state"] = state
// 		}

// 	}

// 	if b, e = json.Marshal(m); e != nil {
// 		return e
// 	}

// 	return json.Unmarshal(b, &a)
// }

// // Value marshalls Messages.
// func (m Messages) Value() (driver.Value, error) {
// 	if m == nil {
// 		return nil, nil
// 	}

// 	return json.Marshal(m)
// }

// // Scan converts raw JSON ([]byte) to Messages.
// func (m *Messages) Scan(value interface{}) error {
// 	if value == nil {
// 		return nil
// 	}

// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	return json.Unmarshal(b, &m)
// }

// // Value to store State in database.
// func (s State) Value() (driver.Value, error) {
// 	return s.String()
// }

// // Scan string stored in database to State type.
// func (s *State) Scan(value interface{}) (e error) {
// 	if value == nil {
// 		return nil
// 	}

// 	i, ok := value.(string)
// 	if !ok {
// 		return errors.New("type assertion to string failed")
// 	}

// 	*s, e = ToState(i)

// 	return e
// }

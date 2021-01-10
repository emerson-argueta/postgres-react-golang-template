package goal

import "encoding/json"

// Messages represents a map of messages for an achiever within a goal
// where the key is a timestamp and the value is the message
type Messages map[Timestamp]Message

// Timestamps represent the timestamps
func (m *Messages) Timestamps() []Timestamp {
	keys := make([]Timestamp, len(*m))

	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	return keys
}

// ToPrimativeType from Messages type
func (m *Messages) ToPrimativeType() map[string]string {
	primativeType := make(map[string]string)

	for k, v := range *m {
		primativeType[k.ToString()] = v.ToString()
	}
	return primativeType
}

// NewMessagesFromByteArray creates new messages. Empty messages by default
func NewMessagesFromByteArray(newMessages *[]byte) Messages {
	if newMessages != nil {
		m := new(Messages)
		json.Unmarshal(*newMessages, m)
		return *m
	}

	return make(Messages)

}

// NewMessages creates new message. Empty messages by default
func NewMessages(newMessages *map[string]string) (res Messages) {
	res = make(Messages)
	if newMessages == nil {
		return res
	}

	for k, v := range *newMessages {
		timestamp, _ := NewTimestamp(&k)
		message := NewMessage(&v)
		res[timestamp] = message
	}
	return res

}

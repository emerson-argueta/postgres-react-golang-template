package goal

// Message holds a message for an achiever's goal
type Message string

// NewMessage create a new message. Empty string if message is nil.
func NewMessage(message *string) Message {
	if message == nil {
		return Message("")
	}
	return Message(*message)
}

// ToString from Message type
func (m *Message) ToString() string {
	return string(*m)
}

package consts

// ChatType
type ChatType string

func (c ChatType) IsGroupChatType() bool {
	return c == "group"
}

func (c ChatType) IsUserChatType() bool {
	return c == "p2p"
}

const (
	ChatTypeGroup ChatType = "group"
	ChatTypeUser  ChatType = "personal"
)

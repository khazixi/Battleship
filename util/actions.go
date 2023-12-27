package util

type ActionType int

type Action interface {
	getEvent() ActionType
}

const (
	List = iota
	Create
	Join
)

type JoinAction struct {
	event  ActionType
	RoomID int
	ID     int
}

type CreateAction struct {
	event ActionType
	ID    int
}

type ListAction struct {
	event ActionType
	ID    int
}

func (a JoinAction) getEvent() ActionType {
	return a.event
}

func CreateJoinAction(id int, roomID int) JoinAction {
	return JoinAction{event: Join, ID: id, RoomID: roomID}
}

func (c CreateAction) getEvent() ActionType {
	return c.event
}

func CreateCreateAction(id int) CreateAction {
	return CreateAction{event: Create, ID: id}
}

func (l ListAction) getEvent() ActionType {
	return l.event
}

func CreateListAction(id int) ListAction {
	return ListAction{event: List, ID: id}
}

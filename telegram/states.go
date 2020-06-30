package telegram

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	HOME = iota
	VERIFICAR
)

type State struct {
	ID   int
	Data map[int]primitive.ObjectID
}

func (s *State) ClearData() {
	s.Data = nil
}

func (s *State) SetData(id int, data primitive.ObjectID) {
	if s.Data == nil {
		s.Data = make(map[int]primitive.ObjectID)
	}
	s.Data[id] = data
}

func (s *State) GetData(id int) (primitive.ObjectID, bool) {
	data, exists := s.Data[id]
	return data, exists
}

package dto

import (
	"encoding/json"
	"errors"
)

type Task struct {
	Title       string
	Description string
}

func (t Task) ValidateToCreate() error {
	switch {
	case t.Title == "":
		return errors.New("title is empty")
	case t.Description == "":
		return errors.New("description is empty")
	default:
		return nil
	}
}

func ToJSON(v any) []byte {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}

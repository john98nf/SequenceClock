package sequence

import "fmt"

type Sequence struct {
	Name      string   `json:"name,omitempty"`
	Framework string   `json:"framework,omitempty"`
	Functions []string `json:"functions,omitempty"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework string,
	functions ...string) (*Sequence, error) {
	if name == "" || framework == "" || len(functions) == 0 {
		return nil, fmt.Errorf("can't create a sequence with empty fields")
	}
	return &Sequence{
		Name:      name,
		Framework: framework,
		Functions: functions,
	}, nil
}

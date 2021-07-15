package sequence

import "fmt"

type Sequence struct {
	Name          string   `json:"name,omitempty"`
	Framework     string   `json:"framework,omitempty"`
	AlgorithmType string   `json:"algorithmType,omitempty"`
	Functions     []string `json:"functions,omitempty"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework,
	algorithmType string,
	functions ...string) (*Sequence, error) {
	if name == "" || algorithmType == "" || framework == "" || len(functions) == 0 {
		return nil, fmt.Errorf("can't create a sequence with empty fields")
	}
	return &Sequence{
		Name:          name,
		Framework:     framework,
		AlgorithmType: algorithmType,
		Functions:     functions,
	}, nil
}

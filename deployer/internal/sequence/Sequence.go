package sequence

import "fmt"

type Sequence struct {
	Name          string   `form:"name" binding:"required"`
	Framework     string   `form:"framework" binding:"required"`
	AlgorithmType string   `form:"algorithm" binding:"required"`
	Functions     []string `form:"functions" binding:"required"`
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

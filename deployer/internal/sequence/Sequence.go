package sequence

import "fmt"

type Sequence struct {
	Name      string   `json:"name,omitempty"`
	Functions []string `json:"functions,omitempty"`
	Framework string   `json:"framework,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework,
	namespace string,
	functions ...string) (*Sequence, error) {
	if name == "" || framework == "" || namespace == "" || len(functions) == 0 {
		return nil, fmt.Errorf("can't create a sequence with empty fields")
	}
	return &Sequence{
		Name:      name,
		Functions: functions,
		Framework: framework,
		Namespace: namespace,
	}, nil
}

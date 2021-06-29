package sequence

type Sequence struct {
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
	Framework string   `json:"framework"`
	Namespace string   `json:"namespace"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework,
	namespace string,
	functions ...string) *Sequence {
	res := &Sequence{
		Name:      name,
		Functions: functions,
		Framework: framework,
		Namespace: namespace,
	}
	return res
}

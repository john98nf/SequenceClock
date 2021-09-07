package sequence

import (
	"fmt"
)

type Sequence struct {
	Name                   string   `form:"name" binding:"required"`
	Framework              string   `form:"framework" binding:"required"`
	AlgorithmType          string   `form:"algorithm" binding:"required"`
	Functions              []string `form:"functions" binding:"required"`
	ProfiledExecutionTimes []int64  `form:"profiledExecutionTimes" binding:"required"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework,
	algorithmType string,
	functions []string,
	profiledExecutionTimes []int64) (*Sequence, error) {
	return &Sequence{
		Name:                   name,
		Framework:              framework,
		AlgorithmType:          algorithmType,
		Functions:              functions,
		ProfiledExecutionTimes: profiledExecutionTimes,
	}, nil
}

/*
	Validates Sequence struct
*/
func (s *Sequence) Validate() error {
	if len(s.ProfiledExecutionTimes) != len(s.Functions) {
		return fmt.Errorf(("inconsistent sequence"))
	}
	return nil
}

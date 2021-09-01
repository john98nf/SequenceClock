package sequence

import (
	"fmt"
	"time"
)

type Sequence struct {
	Name                   string          `form:"name" binding:"required" schema:"name,required"`
	Framework              string          `form:"framework" binding:"required" schema:"framework,required"`
	AlgorithmType          string          `form:"algorithm" binding:"required" schema:"algorithm,required"`
	Functions              []string        `form:"functions" binding:"required" schema:"functions,required"`
	ProfiledExecutionTimes []time.Duration `form:"elapsedTimes" binding:"required" schema:"-"`
}

/*
	Creates a new Sequence
*/
func NewSequence(
	name,
	framework,
	algorithmType string,
	functions []string,
	profiledExecutionTimes []time.Duration) (*Sequence, error) {
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

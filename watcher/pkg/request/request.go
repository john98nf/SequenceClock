package request

type Request struct {
	Function string   `form:"function" binding:"required"`
	Metrics  *Metrics `form:"metrics" binding:"required"`
}

type ResetRequest struct {
	Id       uint64 `form:"id" binding:"required"`
	Function string `form:"function" binding:"required"`
}

type Metrics struct {
	Slack                 int64 `form:"slack" binding:"required"`         // Used by P controller
	SumOfSlack            int64 `form:"sumOfSlack" binding:"required"`    // Used by I controller
	PreviousSlack         int64 `form:"previousSlack" binding:"required"` // Used by D controller
	ProfiledExecutionTime int64 `form:"profiledExecutionTime" binding:"required"`
}

func NewResetRequest(id uint64, function string) *ResetRequest {
	return &ResetRequest{
		Id:       id,
		Function: function,
	}
}

func NewRequest(f string, m *Metrics) *Request {
	return &Request{
		Function: f,
		Metrics:  m,
	}
}

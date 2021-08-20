package request

type Request struct {
	Function string   `form:"function" binding:"required" schema:"function"`
	Metrics  *Metrics `form:"metrics" binding:"required" schema:"metrics"`
}

type ResetRequest struct {
	Id       uint64 `form:"id" binding:"required" schema:"id"`
	Function string `form:"function" binding:"required" schema:"function"`
}

type Metrics struct {
	Slack                 int64 `form:"slack" binding:"required" schema:"slack"`                 // Used by P controller
	SumOfSlack            int64 `form:"sumOfSlack" binding:"required" schema:"sumfSlack"`        // Used by I controller
	PreviousSlack         int64 `form:"previousSlack" binding:"required" schema:"previousSlack"` // Used by D controller
	ProfiledExecutionTime int64 `form:"profiledExecutionTime" binding:"required" schema:"profiledExecutionTime"`
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

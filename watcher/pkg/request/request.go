package request

const (
	SpeedUpRequest RequestType = iota + 1
	SlowDownRequest
	ResetRequest
)

type Request struct {
	Function string      `form:"Function" binding:"required"`
	Type     RequestType `form:"Type" binding:"required"`
	// TO DO: Balance resources using the metrics mentioned below.
	//slack time.Duration `form:"Function" binding:"required"`
	//elapsedTime time.Duration `form:"Function" binding:"required"`
}

func newSpeedUpRequest(fName string) *Request {
	return &Request{
		Function: fName,
		Type:     SpeedUpRequest,
	}
}

func newSlowDownRequest(fName string) *Request {
	return &Request{
		Function: fName,
		Type:     SlowDownRequest,
	}
}

func newResetRequest(fName string) *Request {
	return &Request{
		Function: fName,
		Type:     ResetRequest,
	}
}

type RequestType int

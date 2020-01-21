package handler

type (
	SampleRequest struct{
		PayloadString *string
	}

	SampleResponse struct{
		ResultString *string
	}

	SampleHandler struct{

	}
)

func NewSampleHandler() interface{} {
	return &SampleHandler{}
}

func (h *SampleHandler) Handle(request *SampleRequest) (*SampleResponse, error){
	return &SampleResponse{ResultString: request.PayloadString}, nil
}





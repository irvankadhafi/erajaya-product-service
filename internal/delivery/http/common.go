package http

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func setSuccessResponse(data any) successResponse {
	return successResponse{
		Success: true,
		Data:    data,
	}
}

type errorResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func setErrorMessage(msg any) errorResponse {
	return errorResponse{
		Success: false,
		Message: msg,
	}
}

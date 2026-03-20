package handlerx

import "net/http"

// Response represents the result of a handler execution.
//
// It encapsulates all possible outcomes of a request, including:
//   - HTTP status code
//   - response data (JSON, string, etc.)
//   - file responses
//   - error handling
//   - middleware flow control (Next)
//
// This struct is designed to be interpreted by a Converter,
// which translates it into a specific framework response (Fiber, Echo, etc.).
type Response struct {
	// Err represents an error that occurred during request handling.
	//
	// If set, it usually takes priority over Data and will be
	// serialized as an error response by the converter.
	Err error

	// StatusCode defines the HTTP status code to be returned.
	//
	// If not explicitly set or invalid, it defaults to 200 (OK).
	StatusCode int

	// next indicates whether the request should be passed
	// to the next middleware/handler in the chain.
	next bool

	// FilePath, if set, indicates that a file should be sent
	// as the response instead of JSON or raw data.
	FilePath string

	// Data holds the response payload.
	//
	// It can be:
	//   - struct/map → serialized as JSON
	//   - string     → sent as plain text (depending on converter)
	//   - any other type supported by the converter
	Data any
}

// Next marks the response to pass execution to the next handler.
//
// This is typically used in middleware scenarios.
//
// Example:
//
//	return Response{}.Next()
func (r Response) Next() Response {
	r.next = true
	return r
}

// JSON creates a Response with a JSON payload.
//
// The provided data will be assigned to the Data field,
// and should be serializable by the converter (usually to JSON).
//
// Example:
//
//	return Response{}.JSON(map[string]string{
//	    "message": "ok",
//	})
func (r Response) JSON(data any) Response {
	return Response{
		StatusCode: r.Status(),
		Data:       data,
	}
}

// GoNext returns whether the handler chain should continue.
//
// This method should be used instead of accessing internal fields directly,
// ensuring proper encapsulation.
//
// Example:
//
//	if res.GoNext() {
//	    // call next middleware
//	}
func (r Response) GoNext() bool {
	return r.next
}

// Status returns a valid HTTP status code.
//
// If the defined StatusCode is less than or equal to 200,
// it defaults to http.StatusOK (200).
//
// This ensures that invalid or unset status codes
// do not break the response flow.
func (r Response) Status() int {
	statusCode := r.StatusCode
	initialStatus := http.StatusOK

	if statusCode <= initialStatus {
		statusCode = initialStatus
	}

	return statusCode
}

// SendStatus creates a Response with only a status code.
//
// Useful for simple responses without a body.
//
// Example:
//
//	return SendStatus(404)
func SendStatus(statusCode int) Response {
	return Response{StatusCode: statusCode}
}

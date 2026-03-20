package handlerx

// Handler represents a generic request handler.
//
// It receives a Context abstraction and returns a Response,
// allowing full control over request parsing and response generation.
//
// This design enables framework-independent business logic.
//
// Example:
//
//	func HelloHandler(c Context) Response {
//	    name := c.Query("name", "guest")
//
//	    return Response{
//	        Data: map[string]string{
//	            "message": "Hello " + name,
//	        },
//	    }
//	}
//
// Example with error:
//
//	func ErrorHandler(c Context) Response {
//	    return Response{
//	        Err: errors.New("something went wrong"),
//	        StatusCode: 500,
//	    }
//	}
type Handler func(c Context) Response
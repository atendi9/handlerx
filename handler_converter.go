package handlerx

// Converter defines a generic adapter that transforms a Handler
// into a framework-specific handler type.
//
// This abstraction allows the same business logic (Handler)
// to be reused across different HTTP frameworks.
//
// T represents the target handler type (e.g., fiber.Handler, echo.HandlerFunc).
//
// Example (Fiber):
//
//	type FiberConverter struct{}
//
//	func (f FiberConverter) Convert(h Handler) fiber.Handler {
//	    return func(c *fiber.Ctx) error {
//	        // Wrap framework context into our abstraction
//	        ctx := config.FiberContext{Ctx: c}
//
//	        // Execute handler
//	        res := h(Atendi9Context{ctx})
//
//	        // Middleware flow control
//	        if res.GoNext() {
//	            return c.Next()
//	        }
//
//	        // File response
//	        if len(res.FilePath) > 0 {
//	            return c.SendFile(res.FilePath)
//	        }
//
//	        // Error handling (priority over Data)
//	        if err := res.Err; err != nil {
//	            return c.Status(res.Status()).JSON(fiber.Map{
//	                "err": err.Error(),
//	            })
//	        }
//
//	        // String response optimization
//	        if v, ok := res.Data.(string); ok {
//	            return c.Status(res.Status()).SendString(v)
//	        }
//
//	        // Default: JSON response
//	        return c.Status(res.Status()).JSON(res.Data)
//	    }
//	}
type Converter[T any] interface {
	// Convert transforms a generic Handler into a framework-specific handler.
	Convert(h Handler) T
}
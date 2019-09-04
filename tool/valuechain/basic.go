package valuechain

// Basic handler
type Basic struct {
	next Handler
}

func (handler *Basic) resolve(tag, value string) string {
	return value
}

func (handler *Basic) setNext(next Handler) Handler {
	handler.next = next
	return next
}

package valuechain

// Handler executer interface
type Handler interface {
	resolve(tag, value string) string
	setNext(next Handler) Handler
}

func makeChain(handlders ...Handler) (top Handler) {
	var chain Handler
	for _, handler := range handlders {
		if top == nil {
			top = handler // record first handler
			chain = handler
			continue
		}
		chain = chain.setNext(handler)
	}
	return
}

var chain = makeChain(
	&URLEncoder{},
	&URLDecoder{},
	&Basic{},
)

// Execute chain
func Execute(tag, value string) string {
	return chain.resolve(tag, value)
}

package payments

// PaymentProcessorRegistry provides access to payment processors by provider name.
type PaymentProcessorRegistry interface {
	GetProcessor(provider string) (PaymentProcessor, bool)
}

// MapProcessorRegistry is a map-based implementation of PaymentProcessorRegistry.
type MapProcessorRegistry struct {
	processors map[string]PaymentProcessor
}

// NewMapProcessorRegistry creates a registry from a map of provider name to processor.
func NewMapProcessorRegistry(processors map[string]PaymentProcessor) *MapProcessorRegistry {
	if processors == nil {
		processors = make(map[string]PaymentProcessor)
	}
	return &MapProcessorRegistry{processors: processors}
}

// GetProcessor returns the processor for the given provider, or (nil, false) if not found.
func (r *MapProcessorRegistry) GetProcessor(provider string) (PaymentProcessor, bool) {
	p, ok := r.processors[provider]
	return p, ok
}

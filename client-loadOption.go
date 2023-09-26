package sherlock

const LoadDocumentTypeUnknown = 0

const LoadDocumentTypeActor = 1

const LoadDocumentTypeCollection = 2

const LoadDocumentTypeDocument = 3

type LoadConfig struct {
	DocumentType     int
	MaximumRedirects int
	DefaultValue     map[string]any
}

type LoadOption func(*LoadConfig)

func NewLoadConfig(options ...any) LoadConfig {
	result := LoadConfig{
		MaximumRedirects: 6,
		DocumentType:     LoadDocumentTypeUnknown,
		DefaultValue:     make(map[string]any),
	}

	for _, option := range options {
		if typed, ok := option.(LoadOption); ok {
			typed(&result)
		}
	}
	return result
}

func AsActor() LoadOption {
	return asDocumentType(LoadDocumentTypeActor)
}

func AsDocument() LoadOption {
	return asDocumentType(LoadDocumentTypeDocument)
}

func AsCollection() LoadOption {
	return asDocumentType(LoadDocumentTypeCollection)
}

func asDocumentType(documentType int) LoadOption {
	return func(config *LoadConfig) {
		config.DocumentType = documentType
	}
}

func WithMaximumRedirects(maximumRedirects int) LoadOption {
	return func(config *LoadConfig) {
		config.MaximumRedirects = maximumRedirects
	}
}

func WithDefaultValue(defaultValue map[string]any) LoadOption {
	return func(config *LoadConfig) {
		config.DefaultValue = defaultValue
	}
}

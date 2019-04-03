package converters

type Parser interface {
	Parse(data []byte) (*ParsedData, error)
}

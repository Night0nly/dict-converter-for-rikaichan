package converters

type Converter interface {
	ConvertToEdict(dirname string) error
}

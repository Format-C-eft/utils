package password

type Config struct {
	Length         uint32 `json:"length" yaml:"length"`
	IncludeNumbers bool   `json:"include_numbers" yaml:"include_numbers"`
	IncludeSymbols bool   `json:"include_symbols" yaml:"include_symbols"`
}

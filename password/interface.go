package password

type Generator interface {
	Generate() (string, error)
}

package tool

type Tool interface {
	Run(payload string) (string, error)
}

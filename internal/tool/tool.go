package tool

import "encoding/json"

type Tool[T any] struct {
	args T
	Fn   func(T) string
}

func (t *Tool[T]) parse(payload string) error {
	return json.Unmarshal([]byte(payload), &t.args)
}

func (t *Tool[T]) Run(payload string) (string, error) {
	err := t.parse(payload)
	if err != nil {
		return "", err
	}
	return t.Fn(t.args), nil
}

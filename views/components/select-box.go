package components

type Option struct {
	Key, Value string
}

type SelectBox struct {
	Label   string
	Name    string
	Options []Option
}

package components

type Component struct {
	ID         string
	Name       string
	Type       int
	LoxoneType string `mapstructure:"loxone_type"`
}

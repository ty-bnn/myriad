package values

type MapValue struct {
	Kind ValueKind
	Name string
	Key  Value
}

func (m MapValue) GetKind() ValueKind {
	return m.Kind
}

func (m MapValue) GetName() string {
	return m.Name
}

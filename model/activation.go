package model

type Activation struct {
	Token string `sql:",unique,notnull"`
	Email string `sql:",unique,notnull"`
}

func (a Activation) Valid() bool {
	return true
}

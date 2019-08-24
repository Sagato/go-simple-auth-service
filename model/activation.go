package model

type Activation struct {
	Token string `sql:",pk,unique,notnull"`
	Email string `sql:",unique,notnull"`
}

func (a Activation) Valid() bool {
	return true
}

func (a Activation) IsEmpty() bool {
	return a.Token == "" && a.Email == ""
}

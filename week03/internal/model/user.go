package model

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "user"
}

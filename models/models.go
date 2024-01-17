package models

type Product struct {
	Id   int64 `gorm:"primaryKey"`
	Name string
	Code string
}

type Massage struct {
	ClientId string
	UserId   string
	Msg      []byte
}

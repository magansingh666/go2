package models

type Product struct {
	Id   int64 `gorm:"primaryKey"`
	Name string
	Code string
}

type Request struct {
	ClientId    string
	UserId      string
	MetaData    string
	Data        string
	ContextInfo string
	Command     string
}

type Response struct {
	ClientId    string
	UserId      string
	MetaData    string
	Data        string
	ContextInfo string
	Command     string
}

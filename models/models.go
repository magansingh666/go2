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

// incoming message structure ; C -- command; M -- meta data; D -- actual data
type IMsg struct {
	ClientId string
	UserId   string
	C        []string
	M        string
	D        interface{}
}

// outgoing message structure ; E -- error may be true or false. if it is true and D will have Error info other wise D will be data
type OMsg struct {
	ClientId string
	UserId   string
	E        string
	D        string
}

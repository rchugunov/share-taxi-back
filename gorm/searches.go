package gorm

type SearchesDao interface {
	Connect()
	Disconnect() error
}

type SearchesDaoImpl struct {
	Connection
}

package database

type IpAddress struct {
	AddressType string `db:"address_type"`
	IpStart     []byte `db:"ip_start"`
	IpEnd       []byte `db:"ip_end"`
	Country     string `db:"country"`
	City        string `db:"city"`
}

package jualin

type Order struct {
	Name     string `json:"name" bson:"name"`
	Quantity int    `json:"quantity" bson:"quantity"`
	Price    int    `json:"price" bson:"price"`
}

type User struct {
	Name     string `json:"name" bson:"name"`
	Whatsapp string `json:"whatsapp" bson:"whatsapp"`
	Address  string `json:"address" bson:"address"`
}

type PaymentRequest struct {
	Orders        []Order `json:"orders" bson:"orders"`
	Total         int     `json:"total" bson:"total"`
	User          User    `json:"user" bson:"user"`
	Payment       string  `json:"payment" bson:"payment"`
	PaymentMethod string  `json:"paymentMethod" bson:"paymentMethod"`
}

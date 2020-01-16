package main

//CUSTOMER document structure
type CUSTOMER struct {
	CustomerID     string       `json:"CustomerID"`
	Active         bool         `json:"Active"`
	FirstName      string       `json:"FirstName"`
	LastName       string       `json:"LastName"`
	Email          string       `json:"Email"`
	PhoneNumbers   []string     `json:"PhoneNumbers"`
	Password       string       `json:"Password"`
	AddressBook    []ADDRESS    `json:"AddressBook"`
	PaymentOptions []CREDITCARD `json:"PaymentOptions"`
	WishList       []string     `json:"WishList"`
	SaveForLater   []string     `json:"SaveForLater"`
	Updated        int64        `json:"Updated"`
}
package main

type User struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Msg struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Root struct {
	Project   string `json:"project"`
	Status    int    `json:"status"`
	GoLang    string `json:"golang"`
	DataBase  string `json:"database"`
	Version   string `json:"version"`
	Timezone  string `json:"timezone"`
	Timestamp int64  `json:"timestamp"`
}

type ResponseCustomerOK struct {
	Code    int      `json:"code"`
	Status  string   `json:"status"`
	Message Customer `json:"message"`
}

type ResponseCustomerERROR struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Customer struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	LastName       string `json:"last_name"`
	DocumentNumber int    `json:"document_number"`
	CustomerNumber int    `json:"customer_number"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Enabled        bool   `json:"enabled"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastEntryIVR   string `json:"last_entry_ivr"`
}

type CreateCredit struct {
	ID             *int64  `json:"id,omitempty"`
	UniqueID       string  `json:"uniqueid"`
	Card           int64   `json:"card"`
	ExpirationDate int     `json:"expiration_date"`
	SecurityCode   int     `json:"security_code"`
	Amount         float64 `json:"amount"`
	IDCustomer     int     `json:"id_customer"`
	CreationAt     *string `json:"created_at,omitempty"`
}

type Survey struct {
	UniqueID   string `json:"uniqueid"`
	Agent      int    `json:"agent"`
	Queue      int    `json:"queue"`
	Phone      int64  `json:"phone"`
	IDCustomer int    `json:"id_customer"`
}

type SurveyUpdate struct {
	Q1 int `json:"q1,omitempty"`
	Q2 int `json:"q2,omitempty"`
	Q3 int `json:"q3,omitempty"`
}

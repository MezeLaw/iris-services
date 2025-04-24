package models

type PatientRequest struct {
	ClientID       string `json:"client_id" required:"true"`
	FirstName      string `json:"first_name" required:"true"`
	LastName       string `json:"last_name" required:"true"`
	DocType        string `json:"doc_type" required:"true"`
	DocNumber      string `json:"doc_number" required:"true"`
	BirthDate      string `json:"birth_date" required:"true"`
	Gender         string `json:"gender" required:"true"`
	CountryCode    string `json:"country_code" required:"true"`
	PhoneNumber    string `json:"phone_number" required:"true"`
	Email          string `json:"email" required:"true"`
	AddressStreet  string `json:"address_street" required:"true"`
	AddressNumber  string `json:"address_number" required:"true"`
	AddressCity    string `json:"address_city" required:"true"`
	AddressCountry string `json:"address_country" required:"true"`
	ZipCode        string `json:"zip_code" required:"true"`
}

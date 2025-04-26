package models

type PatientRequest struct {
	ClientID       string                 `json:"client_id" required:"true"`
	FirstName      string                 `json:"first_name" required:"true"`
	LastName       string                 `json:"last_name" required:"true"`
	DocType        string                 `json:"doc_type" required:"true"`
	DocNumber      string                 `json:"doc_number" required:"true"`
	BirthDate      string                 `json:"birth_date" required:"true"`
	Gender         string                 `json:"gender" required:"true"`
	CountryCode    string                 `json:"country_code" required:"true"`
	PhoneNumber    string                 `json:"phone_number" required:"true"`
	Email          string                 `json:"email" required:"true"`
	AddressStreet  string                 `json:"address_street" required:"true"`
	AddressNumber  string                 `json:"address_number" required:"true"`
	AddressCity    string                 `json:"address_city" required:"true"`
	AddressCountry string                 `json:"address_country" required:"true"`
	ZipCode        string                 `json:"zip_code" required:"true"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      string                 `json:"created_at,omitempty"`
	UpdatedAt      string                 `json:"updated_at,omitempty"`
}

type GetPatientRequest struct {
	ClientID  string `json:"client_id"`
	ID        string `json:"id"`
	DocType   string `json:"doc_type"`
	DocNumber string `json:"doc_number"`
}

type Patient struct {
	ID             string                 `dynamodbav:"id"`
	ClientID       string                 `dynamodbav:"client_id"`
	FirstName      string                 `dynamodbav:"first_name"`
	LastName       string                 `dynamodbav:"last_name"`
	DocType        string                 `dynamodbav:"doc_type"`
	DocNumber      string                 `dynamodbav:"doc_number"`
	DocKey         string                 `dynamodbav:"doc_key"`
	BirthDate      string                 `dynamodbav:"birth_date"`
	Gender         string                 `dynamodbav:"gender"`
	CountryCode    string                 `dynamodbav:"country_code"`
	PhoneNumber    string                 `dynamodbav:"phone_number"`
	Email          string                 `dynamodbav:"email"`
	AddressStreet  string                 `dynamodbav:"address_street"`
	AddressNumber  string                 `dynamodbav:"address_number"`
	AddressCity    string                 `dynamodbav:"address_city"`
	AddressCountry string                 `dynamodbav:"address_country"`
	ZipCode        string                 `dynamodbav:"zip_code"`
	CreatedAt      string                 `dynamodbav:"created_at"`
	UpdatedAt      string                 `dynamodbav:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata" dynamodbav:"metadata"`
}

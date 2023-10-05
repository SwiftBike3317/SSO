package types

type Account struct {
	Id         string `db:"id"`
	Patronymic string `db:"patronymic"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	Email      string `db:"email"`
}

type JWT struct {
	AllowedServices string
	User            Account
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

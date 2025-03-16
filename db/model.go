package db

type SignUpForm struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address1    string `json:"address1"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
	DateOfBirth string `json:"dob"`
	AadharNo    string `json:"aadharNo"`
	UserId      string `json:"userId"`
}

type SignInForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
type BankUser struct {
	// gorm.Model
	// UserID      string `gorm:"primaryKey"`
	Email             string `gorm:"primaryKey"`
	Password          string `gorm:"not null"`
	FirstName         string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	Address1          string `gorm:"not null"`
	City              string `gorm:"not null"`
	State             string `gorm:"not null"`
	PostalCode        string `gorm:"not null"`
	DateOfBirth       string `gorm:"not null"`
	DwollaCustomerUrl string `gorm:"not null"`
	DwollaCustomerId  string `gorm:"not null"`
	AadharNo          string `gorm:"unique;not null"`
}
*/

func (BankUser) TableName() string {
	return "bank_users"
}
func (s *SignUpForm) ConvertToUser() *BankUser {
	return &BankUser{
		Email:       s.Email,
		Password:    s.Password,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Address1:    s.Address1,
		City:        s.City,
		State:       s.State,
		PostalCode:  s.PostalCode,
		DateOfBirth: s.DateOfBirth,
		AadharNo:    s.AadharNo,
		UserID:      s.UserId,
	}
}

type LoggedInUser struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address1    string `json:"address1"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
	DateOfBirth string `json:"dob"`
	AadharNo    string `json:"aadharNo"`
	Email       string `json:"email"`
	UserId      string `json:"userId"`
}

type BankUser struct {
	// gorm.Model
	Email             string `gorm:"primaryKey" json:"email"`
	Password          string `gorm:"not null" json:"password"`
	FirstName         string `gorm:"not null" json:"firstName"`
	LastName          string `gorm:"not null" json:"lastName"`
	Address1          string `gorm:"not null" json:"address1"`
	City              string `gorm:"not null" json:"city"`
	State             string `gorm:"not null" json:"state"`
	PostalCode        string `gorm:"not null" json:"postalCode"`
	DateOfBirth       string `gorm:"not null" json:"dateOfBirth"`
	DwollaCustomerUrl string `gorm:"not null" json:"dwollaCustomerUrl"`
	DwollaCustomerId  string `gorm:"not null" json:"dwollaCustomerId"`
	AadharNo          string `gorm:"unique;not null" json:"aadharNo"`
	UserID            string `gorm:"unique;not null" json:"userId"`
}

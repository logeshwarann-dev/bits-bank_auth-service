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
}

type SignInForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BankUser struct {
	// gorm.Model
	// UserID      string `gorm:"primaryKey"`
	Email       string `gorm:"primaryKey"`
	Password    string `gorm:"not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Address1    string `gorm:"not null"`
	City        string `gorm:"not null"`
	State       string `gorm:"not null"`
	PostalCode  string `gorm:"not null"`
	DateOfBirth string `gorm:"not null"`
	AadharNo    string `gorm:"unique;not null"`
}

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
	}
}

package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthServiceJWT is struct jwt token implemented on auth service
type AuthServiceJWT struct {
	jwt.StandardClaims
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	Sub            int       `json:"sub"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	IsAdmin        int       `json:"isAdmin"`
	Status         string    `json:"status"`
	IsNbiot        bool      `json:"is_nbiot"`
	SubdistrictID  int32     `json:"subdistrict_id"`
	ProvCode       int64     `json:"prov_code"`
	CustomerNumber string    `json:"customer_number"`
	AppUserID      int       `json:"app_user_id"`
	UserName       string    `json:"username"`
	Address        string    `json:"address"`
	Picture        string    `json:"picture"`
	CompanyName    string    `json:"company_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	AppToken       string    `json:"app_token"`
	CustID         string    `json:"cust_id"`
	Location       string    `json:"location"`
	Interest       string    `json:"interest"`
	Version        int       `json:"version"`
}

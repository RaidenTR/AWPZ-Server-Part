package authorization

import (
	"AWPZ/internal/authorizationdata"
	"AWPZ/internal/database"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	Admin  = "admin"
	Lector = "lector"
	Device = "device"
)

var const secretKey = []byte("very secret key")

type JWTClaims struct {
	jwt.StandardClaims
	ID int
}

func ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, getKeyString)
	if err == nil && token.Valid {
		return true
	}
	return false
}

func GetIDFromToken(tokenString string) int {
	token, _ := jwt.Parse(tokenString, getKeyString)
	claims := token.Claims.(jwt.MapClaims)
	ID := int(claims["ID"].(float64))
	return ID
}
func GetDeviceRoom(tokenString string) string {
	ID := GetIDFromToken(tokenString)
	room := database.GetDeviceRoom(ID)
	return room
}

func IsAdmin(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, getKeyString)
	claims := token.Claims.(jwt.MapClaims)
	return (&claims).VerifyAudience(Admin, true)
}

func getKeyString(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Sign method is invalid")
	}
	return secretKey, nil
}


func getLectorToken(loginData authorizationdata.Set) (string, error) {
	if !database.IsAuthenticatedLector(loginData) {
		return "", errors.New("Wrong login or password")
	}
	ID := database.GetLectorID(loginData)
	result := createToken(ID, loginData.AccessLvl)
	return result, nil
}


func getAdminToken(loginData authorizationdata.Set) (string, error) {
	if !database.IsAuthenticatedAdmin(loginData) {
		return "", errors.New("Wrong login or password")
	}
	ID := database.GetAdminID(loginData)
	result := createToken(ID, loginData.AccessLvl)
	return result, nil
}

func createToken(id int, user string) string {
	claims := createClaims(user, id)
	tokenString := returnSignedString(claims)
	return tokenString
}

func createClaims(user string, id int) JWTClaims {
	expires := time.Now().Add(time.Hour * 4).Unix()
	resultClaims := JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: expires,
			Audience:  user,
		}, id,
	}
	return resultClaims
}

func returnSignedString(claims JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secretKey)
	return tokenString
}

func getDeviceToken(loginData authorizationdata.Set) (string, error) {
	if !database.IsAuthenticatedDevice(loginData) {
		return "", errors.New("Wrong device ID")
	}
	ID := database.GetDeviceID(loginData)
	result := createToken(ID, loginData.AccessLvl)
	return result, nil
}

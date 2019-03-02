package authorization

import (
	"AWPZ/internal/authorizationdata"
)

func GetLectorToken(loginData authorizationdata.Set) (string, error) {
	result, err := getLectorToken(loginData)
	return result, err
}

func GetAdminToken(loginData authorizationdata.Set) (string, error) {
	result, err := getAdminToken(loginData)
	return result, err
}

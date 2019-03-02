package authorization

import (
	"AWPZ/internal/authorizationdata"
)

func GetDeviceToken(id string) (string, error) {
	loginData := authorizationdata.Set{
		Login:     id,
		AccessLvl: Device}
	result, err := getDeviceToken(loginData)
	return result, err
}

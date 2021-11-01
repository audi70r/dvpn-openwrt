package auth

import "encoding/json"

func ValidateAndUnmarshal(req []byte) (login LoginRequest, err error) {
	if err = json.Unmarshal(req, &login); err != nil {
		return login, nil
	}

	return login, nil
}

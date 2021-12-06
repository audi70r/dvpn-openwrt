package keys

import "encoding/json"

func ValidateAndUnmarshalRecovery(req []byte) (keys AddRecoverRequest, err error) {
	if err = json.Unmarshal(req, &keys); err != nil {
		return keys, nil
	}

	return keys, nil
}

func ValidateAndUnmarshalDeletion(req []byte) (keys DeleteRequest, err error) {
	if err = json.Unmarshal(req, &keys); err != nil {
		return keys, nil
	}

	return keys, nil
}

package sonar

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func NewToken(sonarURL string) (string, error) {
	name := uuid.NewString()

	url := fmt.Sprintf("%s/api/user_tokens/generate?name=%s&type=USER_TOKEN", sonarURL, name)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("fail to create a request to get a new token: %w", err)
	}

	req.SetBasicAuth("admin", "admin")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get a new token, status code from api: %d", res.StatusCode)
	}

	data := make(map[string]interface{})

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", fmt.Errorf("fail to decode the new token from api: %w", err)
	}

	return data["token"].(string), nil
}

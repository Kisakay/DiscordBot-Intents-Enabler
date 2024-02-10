package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

func checkIntent(token string) (map[string]interface{}, error) {
	url := "https://discord.com/api/v10/applications/@me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	flags := result["flags"].(float64)

	intentMap := map[string]float64{
		"APPLICATION_AUTO_MODERATION_RULE_CREATE_BADGE": 1 << 6,
		"Presence":                         1 << 13,
		"GuildMember":                      1 << 14,
		"VERIFICATION_PENDING_GUILD_LIMIT": 1 << 16,
		"EMBEDDED":                         1 << 17,
		"MessageContent":                   1 << 18,
		"APPLICATION_COMMAND_BADGE":        1 << 23,
	}

	flagsMissing := make([]string, 0)

	for k, v := range intentMap {
		if int(flags)&int(v) == 0 {
			flagsMissing = append(flagsMissing, k)
		}
	}

	if len(flagsMissing) > 0 {
		return map[string]interface{}{
			"state": false,
			"data":  flagsMissing,
		}, nil
	} else {
		return map[string]interface{}{
			"data": make([]interface{}, 0),
		}, nil
	}
}

func activeIntent(token string) error {
	url := "https://discord.com/api/v10/applications/@me"

	payload := map[string]interface{}{
		"flags": 565248,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func user(token string) (map[string]interface{}, error) {
	url := "https://discord.com/api/v10/users/@me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var userData UserResponse
	json.Unmarshal(body, &userData)

	return map[string]interface{}{
		"valid": resp.Status == "OK",
		"id":    userData.ID,
		"tag":   userData.Username + "#" + userData.Discriminator,
	}, nil
}

func main() {
	// utilisation des fonctions
	fmt.Print(activeIntent("MTE5Nzk1NjI3MDc0NTk5MzMzOA.GHtiEe.mRu2I0X63tKtZ4md2wrEzy6dD7VWjTOHIjq4t0"))
}

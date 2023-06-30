package main

import (
	"encoding/json"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/VichyGopher/hcmonitor/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type HcaptchaResponse struct {
	Features struct {
		A11YChallenge bool `json:"a11y_challenge"`
	} `json:"features"`
	C struct {
		Type string `json:"type"`
		Req  string `json:"req"`
	} `json:"c"`
	Pass bool `json:"pass"`
}

type JwtToken struct {
	F                         float64 `json:"f,omitempty"`
	S                         float64 `json:"s,omitempty"`
	ScriptType                string  `json:"s,omitempty"`
	EncryptionKey             string  `json:"d,omitempty"`
	RequestChallengeTimestamp float64 `json:"e,omitempty"`
	NType                     string  `json:"n,omitempty"`
	C                         float64 `json:"c,omitempty"`
	VersionBaseUrl            string  `json:"l,omitempty"`
}

func ScrapeJwt(v *HcVersion) (*JwtToken, error) {
	client := cycletls.Init()

	response, err := client.Do(fmt.Sprintf("https://api2.hcaptcha.com/checksiteconfig?v=%s&host=discord.com&sitekey=4c672d35-0701-42b2-88c3-78380b0db560&sc=1&swa=1&spst=0", v.Version), cycletls.Options{
		Body:      "",
		Ja3:       "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,27-51-23-18-45-13-10-16-43-5-17513-11-0-65281-35-21,29-23-24,0",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		Headers: map[string]string{
			`authority`:          `hcaptcha.com`,
			`accept`:             `application/json`,
			`accept-language`:    `fr-FR,fr;q=0.9`,
			`content-type`:       `text/plain`,
			`origin`:             `https://newassets.hcaptcha.com`,
			`referer`:            `https://newassets.hcaptcha.com/`,
			`sec-ch-ua`:          `"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`,
			`sec-ch-ua-mobile`:   `?0`,
			`sec-ch-ua-platform`: `"macOS"`,
			`sec-fetch-dest`:     `empty`,
			`sec-fetch-mode`:     `cors`,
			`sec-fetch-site`:     `same-site`,
			`user-agent`:         `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36`,
		},
	}, "GET")
	if utils.HandleError(err) {
		return nil, err
	}

	var data HcaptchaResponse
	if err := json.Unmarshal([]byte(response.Body), &data); err != nil {
		return nil, err
	}

	payload, _ := jwt.Parse(data.C.Req, nil)
	claims, _ := payload.Claims.(jwt.MapClaims)

	if claims["l"].(string) == "" {
		return nil, fmt.Errorf("cant parse jwt token")
	}

	return &JwtToken{
		F:                         claims["f"].(float64),
		S:                         claims["s"].(float64),
		ScriptType:                claims["t"].(string),
		EncryptionKey:             claims["d"].(string),
		RequestChallengeTimestamp: claims["e"].(float64),
		NType:                     claims["n"].(string),
		C:                         claims["c"].(float64),
		VersionBaseUrl:            claims["l"].(string),
	}, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"net/url"
)

var scanUrlInternal = "%s/scan"

type ScanUrlPayload struct {
	Url               string `json:"url"`
	RecaptchaResponse string `json:"recaptcha_token"`
	Id                string `json:"id"`
}

type InternalAPIResult struct {
	Screenshot string `json:"screen_shot"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Prediction string `json:"prediction"`
	Percentage string `json:"percentage"`
}

func ScanURL(c *gin.Context) {
	cmd := new(ScanUrlPayload)
	if err := BindJSON(c.Request, cmd); err != nil {
		c.JSON(402, ResponseBody{
			Message: "Invalid payload",
		})
		return
	}

	if !validateURL(cmd.Url) {
		c.JSON(402, ResponseBody{
			Message: "invalid url",
		})
		return
	}

	isHuman, err := VerifyCaptcha(cmd.RecaptchaResponse, config.GoogleReCaptchaSecret)
	if err != nil || !isHuman {
		c.JSON(401, ResponseBody{
			Message: "Prove that you are not a robot",
		})
		return
	}

	if !checkWebsiteAlive(cmd.Url) {
		c.JSON(402, ResponseBody{
			Message: "the site is unreachable",
		})
		return
	}
	cmd.Id = bson.NewObjectId().Hex()
	internalJsonPayload, _ := json.Marshal(cmd)
	req, err := http.NewRequest("POST", fmt.Sprintf(scanUrlInternal, config.InternalAPI), bytes.NewBuffer(internalJsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(402, ResponseBody{
			Message: err.Error(),
		})
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	result := new(InternalAPIResult)
	if err := json.Unmarshal(body, result); err != nil {
		c.JSON(402, ResponseBody{
			Message: "Invalid response" + err.Error(),
		})
		return
	}
	c.JSON(200, ResponseBody{
		Message: "Analyzed successfully",
		Data: map[string]interface{}{
			"screenshot": result.Screenshot,
			"prediction": result.Prediction,
			"percentage": result.Percentage,
			"code":       result.Code,
			"url":        cmd.Url,
		},
	})
}

func validateURL(raw string) bool {
	_, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return true
}

func checkWebsiteAlive(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode == 404 || resp.StatusCode == 500 {
		return false
	}

	return true
}

const GOOGLE_VERIFY_CAPCHA_URL = "https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s"

type GoogleCaptchaResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"error-codes"`
}

func VerifyCaptcha(captcha string, secret string) (bool, error) {
	response, e := http.Post(fmt.Sprintf(GOOGLE_VERIFY_CAPCHA_URL, secret, captcha), "application/json", nil)
	if e != nil {
		return false, e
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return false, fmt.Errorf("Validation Error")
	}

	buffer, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return false, e
	}

	data := new(GoogleCaptchaResponse)
	e = json.Unmarshal(buffer, data)
	if e != nil {
		return false, e
	}
	return data.Success, nil
}

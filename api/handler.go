package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
		c.JSON(422, ResponseBody{
			Message: "Invalid payload",
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

	if !govalidator.IsURL(cmd.Url) {
		c.JSON(422, ResponseBody{
			Message: "invalid url",
		})
		return
	}

	if !checkWebsiteAlive(cmd.Url) {
		c.JSON(422, ResponseBody{
			Message: "the site is unreachable",
		})
		return
	}
	result, err := InternalAPI(*cmd)
	if err != nil {
		c.JSON(422, ResponseBody{
			Message: "Error" + err.Error(),
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

func checkWebsiteAlive(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode == 404 || resp.StatusCode == 500 {
		return false
	}

	return true
}

type AddMonitorPayload struct {
	ID                bson.ObjectId `json:"_id" bson:"_id"`
	URL               string        `json:"url" bson:"url"`
	Email             string        `json:"email" bson:"email"`
	Interval          int           `json:"interval" bson:"interval"`
	RecaptchaResponse string        `json:"recaptcha_token" bson:"-"`
}

func AddMonitorSchedule(c *gin.Context) {
	cmd := new(AddMonitorPayload)
	if err := BindJSON(c.Request, cmd); err != nil {
		c.JSON(422, ResponseBody{
			Message: "Invalid payload",
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

	if !govalidator.IsURL(cmd.URL) {
		c.JSON(422, ResponseBody{
			Message: "invalid url",
		})
		return
	}

	if !govalidator.IsEmail(cmd.Email) {
		c.JSON(422, ResponseBody{
			Message: "invalid email",
		})
		return
	}

	if cmd.Interval <= 0{
		c.JSON(422, ResponseBody{
			Message: "interval must be > 0",
		})
		return
	}

	cmd.ID = bson.NewObjectId()
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	//_, err = container.MongoClient.Database(container.Config.MongoDatabase).Collection(container.Config.MongoCollection).InsertOne(ctx, cmd)
	//if err != nil {
	//	c.JSON(422, ResponseBody{
	//		Message: "Cannot insert new schedule",
	//	})
	//	fmt.Println(err)
	//	return
	//}

	cronInterval := fmt.Sprintf("@every %dm", cmd.Interval)
	err = container.CronDaemon.AddFunc(cronInterval, func() {
		err := MonitorJob(*cmd)
		if err != nil {
			fmt.Printf("Error monitoring URL: %s %s\n", cmd.URL, err)
		}
		return
	})
	defacedLog := fmt.Sprintf("[%s] New Monitor: %s - %s - %d minutes\n", time.Now().Format(time.RFC3339), cmd.URL, cmd.Email, cmd.Interval)
	SaveLog(defacedLog, "monitor.log")
	if err != nil {
		c.JSON(422, ResponseBody{
			Message: "Cannot create new cron job",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, ResponseBody{
		Message: "Your monitor has been created",
	})
}

func MonitorJob(result AddMonitorPayload) error {
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//var result AddMonitorPayload
	//err := container.MongoClient.Database(container.Config.MongoDatabase).Collection(container.Config.MongoCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	//if err != nil {
	//	return err
	//}

	monitorInternal := ScanUrlPayload{
		Url: result.URL,
	}
	prediction, err := InternalAPI(monitorInternal)
	if err != nil {
		return err
	}
	if prediction.Prediction == "0" {
		return nil
	}
	defacedLog := fmt.Sprintf("[%s] Detected new defacement: %s - %s\n", time.Now().Format(time.RFC3339), result.URL, result.Email)
	SaveLog(defacedLog, "deface.log")
	log.Println(defacedLog)
	content := fmt.Sprintf(`
		Hi, <br/>
		<br/>
		Your site at %s may have been defaced. Please have it checked.<br/><br/>
		Sincerely, <br>
		Defacetor
	`, monitorInternal.Url)

	body := map[string]interface{}{
		"from":        container.MailRepository.Username,
		"to":          []string{result.Email},
		"subject":     fmt.Sprintf("Defacetor Alert!"),
		"contentType": "text/html",
		"content":     content,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = container.MailRepository.SendByGmail(bodyJson)
	if err != nil {
		return err
	}

	return nil
}

func InternalAPI(cmd ScanUrlPayload) (InternalAPIResult, error) {
	cmd.Id = bson.NewObjectId().Hex()
	internalJsonPayload, _ := json.Marshal(cmd)
	req, err := http.NewRequest("POST", fmt.Sprintf(scanUrlInternal, config.InternalAPI), bytes.NewBuffer(internalJsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return InternalAPIResult{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	result := new(InternalAPIResult)
	if err := json.Unmarshal(body, result); err != nil {
		return InternalAPIResult{}, err
	}
	return *result, err
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

func SaveLog(text string, fileName string) {
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Println(err)
	}
}

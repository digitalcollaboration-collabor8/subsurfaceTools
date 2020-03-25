package cloud

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type AuthConfig struct {
	ClientID        string
	ClientSecret    string
	TokenURL        string
	ResourceID      string
	SubScriptionKey string
	Token           string
}

//StringRFC3339ToTime will take a RFC3339 formatted string, e.g. 2005-01-05T21:59:59.999Z
// and convert it to a time object
func StringRFC3339ToTime(timeStamp string) (time.Time, error) {
	return time.Parse(
		time.RFC3339,
		timeStamp)
}

func TimeToStr(timeObj time.Time, format string) string {

	return timeObj.Format(format)
}

func Authenticate() (string, error) {
	zap.S().Infof("Authenticating")
	var resp *http.Response
	var err error
	var clientId, clientSecret, tokenUrl, resourceId string
	//read all of the variables from the environment settings
	clientId = os.Getenv(AzureClientIdEnvName)
	clientSecret = os.Getenv(AzureClientSecretEnvName)
	tokenUrl = os.Getenv(AzureTokenUrlEnvName)
	resourceId = os.Getenv(AzureResourceIdEnvName)
	if clientId == "" || clientSecret == "" || tokenUrl == "" || resourceId == "" {
		errorMessage := "Failed in authentication, not all required environment variables seem to be set:AzureClientId,AzureClientSecret,AzureTokenUrl,AzureResourceId,AzureSubscriptionKey"
		zap.S().Errorf(errorMessage)
		return "", errors.New(errorMessage)
	}

	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"resource":      {resourceId},
	}
	if resp, err = http.PostForm(tokenUrl, formData); err != nil {
		zap.S().Errorf("Error in authentication post:%s", err.Error())
		return "", err
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode == 200 {

		zap.S().Debugf("Token type%s, Resource:%s, Token:%s", result["token_type"],
			result["resource"], result["access_token"])
		return fmt.Sprintf("%v", result["access_token"]), nil
	} else {
		//need to check the error
		zap.S().Errorf("Authentication service responded with code%s,error:%s", resp.StatusCode, result["error_description"])
		return "", errors.New(fmt.Sprintf("%v", result["error_description"]))
	}

}

//DaysBetween will check the numder of days between startdate and enddate
//where the dates are in the format YYYY-MM-DD e.g. 2019-03-28
func DaysBetween(startDate, endDate string) (float64, error) {
	//parse the times
	var start, end time.Time
	var err error
	if start, err = time.Parse("2006-01-02", startDate); err != nil {
		return 0, err
	}
	if end, err = time.Parse("2006-01-02", endDate); err != nil {
		return 0, err
	}
	return math.Ceil(end.Sub(start).Hours() / 24.0), nil
}

func CreateUUID() string {

	u1 := uuid.NewV4()
	return u1.String()
}

func CloudConfigArrayToStruct(data []byte) (CloudDownload, error) {
	var cConfig CloudDownload
	var err error
	//unmarshal it
	if err = xml.Unmarshal(data, &cConfig); err != nil {
		return cConfig, err
	}

	return cConfig, err
}

package cloud

import "encoding/xml"

//globals pointing to enviroinment variables to look for
var AzureClientIdEnvName = "AzureClientId"
var AzureClientSecretEnvName = "AzureClientSecret"
var AzureTokenUrlEnvName = "AzureTokenUrl"
var AzureResourceIdEnvName = "AzureResourceId"
var AzureFileDownloadUrlEnvName = "AzureFileDownloadUrl"
var AzureSubscriptionKeyEnvName = "AzureSubscriptionKey"
var AzureGraphUrlEnvName = "AzureGraphUrl"
var MaxRollDays = 30             //control the max number of days that you can roll back in the download setup
var MaxNumberOfDaysPeriod = 91.0 //controls the number of days that can be queried for using the period functionality
var utcTimeSTamp = "2005-01-05T21:59:59.999Z"

type CloudDownload struct {
	XMLName     xml.Name                `xml:"subsurface"`
	CloudConfig CloudConfig             `xml:"config"`
	DPRS        []CloudProductionConfig `xml:"dpr"`
	MPRGovs     []CloudProductionConfig `xml:"mprmlGov"`
	MPRPartners []CloudProductionConfig `xml:"mprmlPartner"`
	DDRMLS      []CloudDDRMLConfig      `xml:"ddrml"`
}

type CloudConfig struct {
	ClientId        string `xml:"clientId"`
	ClientSecret    string `xml:"clientSecret"`
	TokenURL        string `xml:"tokenUrl"`
	ResourceId      string `xml:"resourceId"`
	FileDownloadUrl string `xml:"fileDownloadUrl"`
	SubscriptionKey string `xml:"subscriptionKey"`
	GraphURL        string `xml:"graphUrl"`
}

type CloudProductionConfig struct {
	FieldName       string `xml:"fieldName"`
	DateFrom        string `xml:"dateFrom"`
	DateTo          string `xml:"dateTo"`
	RollDays        int    `xml:"rollDays"`
	RollMonths      int    `xml:"rollMonths"`
	UseUploadedFrom bool   `xml:"useUploadedFrom"`
	LogFile         string `xml:"logFile"`
	Common          CloudCommonConfig
}

type CloudDDRMLConfig struct {
	DateFrom        string `xml:"dateFrom"`
	DateTo          string `xml:"dateTo"`
	RollDays        int    `xml:"rollDays"`
	UseUploadedFrom bool   `xml:"useUploadedFrom"`
	LogFile         string `xml:"logFile"`
	Common          CloudCommonConfig
}

type CloudCommonConfig struct {
	XMLName          xml.Name `xml:"common"`
	Format           string   `xml:"format"`
	OutputFolder     string   `xml:"outputFolder"`
	FileOutputPrefix string   `xml:"fileOutputPrefix"`
}

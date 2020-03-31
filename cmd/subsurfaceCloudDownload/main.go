package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/cloud"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Version string
	Build   string
)

func readCloudConfig2Struct(configFile string) (cloud.CloudDownload, error) {
	var err error
	var config cloud.CloudDownload
	var data []byte
	var xmlFile *os.File
	if xmlFile, err = os.Open(configFile); err != nil {
		return config, err
	}
	defer xmlFile.Close()
	if data, err = ioutil.ReadAll(xmlFile); err != nil {
		return config, err
	}
	//unmarshal it
	if err = xml.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}

func runDownload(configFile, logConfig string) {
	var cfg zap.Config
	var cloudCnfg cloud.CloudDownload
	logFileName := "log_" + cloud.TimeToStr(time.Now(), "2006-01-02T15_04_05") + ".json"
	errorFileName := "log_errors_" + cloud.TimeToStr(time.Now(), "2006-01-02T15_04_05") + ".json"
	if logConfig == "" {
		cfg = zap.Config{

			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout", "./" + logFileName},
			ErrorOutputPaths: []string{"stderr", "./" + errorFileName},
		}
		//just encode the time in RFC3339 as the default is epoch from production config
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	} else {
		//parse in the log configuration
		if rawJSON, err := ioutil.ReadFile(logConfig); err != nil {
			panic(err)
		} else {
			if err := json.Unmarshal(rawJSON, &cfg); err != nil {
				panic(err)
			}
		}
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	//read the cloud config
	if cloudCnfg, err = readCloudConfig2Struct(configFile); err != nil {
		zap.S().Errorf("Failed in reading cloud configuration xml file:%s", err.Error())
		return
	}
	//now have everything start processing
	//first check if we have any of the cloud params specified in the incoming xml file
	//if so set them as environment variables
	setEnvironments(cloudCnfg)
	cloud.ProcessAndRunDownload(cloudCnfg)
}

//setEnvironments just checks if we have recieved any configuration
//parameters for the cloud through the xml file if so just override the environment variables
func setEnvironments(cloudCnfg cloud.CloudDownload) {
	if cloudCnfg.CloudConfig.ClientId != "" {
		os.Setenv(cloud.AzureClientIdEnvName,
			cloudCnfg.CloudConfig.ClientId)
	}
	if cloudCnfg.CloudConfig.ClientSecret != "" {
		os.Setenv(cloud.AzureClientSecretEnvName,
			cloudCnfg.CloudConfig.ClientSecret)
	}
	if cloudCnfg.CloudConfig.TokenURL != "" {
		os.Setenv(cloud.AzureTokenUrlEnvName,
			cloudCnfg.CloudConfig.TokenURL)
	}
	if cloudCnfg.CloudConfig.ResourceId != "" {
		os.Setenv(cloud.AzureResourceIdEnvName,
			cloudCnfg.CloudConfig.ResourceId)
	}
	if cloudCnfg.CloudConfig.FileDownloadUrl != "" {
		os.Setenv(cloud.AzureFileDownloadUrlEnvName,
			cloudCnfg.CloudConfig.FileDownloadUrl)
	}
	if cloudCnfg.CloudConfig.SubscriptionKey != "" {
		os.Setenv(cloud.AzureSubscriptionKeyEnvName,
			cloudCnfg.CloudConfig.SubscriptionKey)
	}
	if cloudCnfg.CloudConfig.GraphURL != "" {
		os.Setenv(cloud.AzureGraphUrlEnvName,
			cloudCnfg.CloudConfig.GraphURL)
	}
}

func main() {

	configFile := flag.String("configuration", "", "Path to the xml configuration file to use")
	logConfig := flag.String("logconfiguration", "", "Path to the log configuration file")
	showVersion := flag.Bool("version", false, "If specified will print out the version information and then exit")

	flag.Parse()
	if *showVersion {
		fmt.Println("Version:", Version)
		fmt.Println("Build time:", Build)
		return
	}
	if *configFile != "" {
		runDownload(*configFile, *logConfig)
	} else {
		fmt.Println("Missing configuration parameter that should point to a valid configuration xml file")
		flag.PrintDefaults()
		return
	}

}

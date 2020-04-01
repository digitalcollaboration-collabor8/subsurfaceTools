package cloud

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

//GenerateFileName builds a filename based on the incoming criteria
func GenerateFileName(fObj FileObject, filePrefix, format string) string {
	createdStr := ""
	//generate a filename in the form of
	//filePrefix_created_reportType_periodStart_periodEnd_fileReferenceId_storedFileName.extension

	if created, err := StringRFC3339ToTime(fObj.Created); err != nil {
		zap.S().Errorf("Failed in converting RFC3339 timestamp to timeobj:%s", err.Error())

	} else {
		createdStr = TimeToStr(created, "2006-01-02T15_04_05")
	}

	fileName := filePrefix + "_" + createdStr
	baseFileName, _ := common.GetFileNameAndExtension(fObj.FileName)
	baseFileName = SafeEncodeNameForWinFiles(baseFileName)
	reportType := MapReportType(fObj.ReportType)
	sourceString := ""
	if strings.ToLower(reportType) == "ddrml" {
		for i := 0; i < len(fObj.Sources); i++ {
			sourceString = sourceString + "_" + SafeEncodeNameForWinFiles(fObj.Sources[i].Name)
		}
	}
	fileName = fileName + "_" + reportType + "_" + sourceString + "_" + fObj.MetaData.PeriodStart + "_" +
		fObj.MetaData.PeriodEnd +
		"_" + fObj.FileReference + "_" + baseFileName + "." + strings.ToLower(format)
	return fileName
}

func SafeEncodeNameForWinFiles(name string) string {
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "*", "_")
	return name
}

func MapReportType(reportType int) string {
	switch reportType {
	case 1:
		return "DPR10"
	case 2:
		return "DPR20"
	case 3:
		return "DDRML"
	case 4:
		return "MPRMLGov"
	case 5:
		return "MPRMLPartner"
	default:
		return "NONE_REPORTTYPE"
	}
}

//DownloadFiles downloads a set of files from a array of fileobjects, it will call download file function for each entry
// and dowload the data to the outputfolder, will return an array of possible errors
//assetName will be used to prefix the files
//addTimeStampInName will add a unix timestamp to the outputfile name identifying when the file was downloaded
func DownloadFiles(files []FileObject, fileURL, token, subscriptionKey, format,
	outputFolder, filePrefix string, addTimeStampInName bool) []error {
	var errorsEncountered []error
	for i := 0; i < len(files); i++ {
		zap.S().Debugf("Processing file:" + files[i].FileName)

		if fileData, err := DownloadFile(files[i].FileReference,
			fileURL, token, subscriptionKey, format); err != nil {
			errorMsg := fmt.Sprintf("Failed in download of file with referenceId:%s,fileName:%s,format:%s,error:%s",
				files[i].FileReference, files[i].FileName, format, err.Error())
			zap.S().Error(errorMsg)
			errorsEncountered = append(errorsEncountered, errors.New(errorMsg))
		} else {
			//we have the file now write it to disk.
			outputFiles := BuildOutputPathForReportType(files[i], filePrefix, outputFolder, format)
			//need to handle several paths and create the folders if needed
			for x := 0; x < len(outputFiles); x++ {
				//just check that the folder exists
				if !FileOrFolderExists(filepath.Dir(outputFiles[x])) {
					//folder does not exist just create it
					if err = os.MkdirAll(filepath.Dir(outputFiles[x]), os.ModePerm); err != nil {
						errorMsg := fmt.Sprintf("Failed in creating folder for storage with path:%s,error:%s", outputFiles[x], err.Error())
						zap.S().Error(errorMsg)
						errorsEncountered = append(errorsEncountered, errors.New(errorMsg))
						//just abort operation
						return errorsEncountered
					}
				}
				if err := common.Write2File(outputFiles[x], fileData); err != nil {
					errorMsg := fmt.Sprintf("Failed in write of file with referenceId:%s,fileName:%s,outputLocation:%s,error:%s",
						files[i].FileReference, files[i].FileName, outputFiles[x], err.Error())
					zap.S().Error(errorMsg)
					errorsEncountered = append(errorsEncountered, errors.New(errorMsg))
				} else {
					zap.S().Infof("Wrote file to:%s", outputFiles[x])
				}
			}
		}

	}
	return errorsEncountered
}

func FileOrFolderExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

//VerifyCloudDownloadConfig will check the configuration data for errors
func VerifyCloudDownloadConfig(downloadConfig CloudDownload) error {
	//first check that no rolldays are greater than 10
	for i := 0; i < len(downloadConfig.DPRS); i++ {
		if downloadConfig.DPRS[i].RollDays > MaxRollDays {
			return fmt.Errorf("Invalid number of rolldays specified:%d, allowed:%d",
				downloadConfig.DPRS[i].RollDays,
				MaxRollDays)
		}
		//check the days
		if downloadConfig.DPRS[i].DateFrom != "" && downloadConfig.DPRS[i].DateTo != "" {
			if days, err := DaysBetween(downloadConfig.DPRS[i].DateFrom, downloadConfig.DPRS[i].DateTo); err != nil {
				return err
			} else {
				if days > MaxNumberOfDaysPeriod {
					return fmt.Errorf("Max number of days for a given period exceeded:%.2f, max set to:%.2f",
						days, MaxNumberOfDaysPeriod)
				}
			}
		}
	}
	for i := 0; i < len(downloadConfig.DDRMLS); i++ {
		if downloadConfig.DDRMLS[i].RollDays > MaxRollDays {
			return fmt.Errorf("Invalid number of rolldays specified:%d, allowed:%d",
				downloadConfig.DDRMLS[i].RollDays,
				MaxRollDays)
		}
		//check the days
		if downloadConfig.DDRMLS[i].DateFrom != "" && downloadConfig.DDRMLS[i].DateTo != "" {
			if days, err := DaysBetween(downloadConfig.DDRMLS[i].DateFrom, downloadConfig.DDRMLS[i].DateTo); err != nil {
				return err
			} else {
				if days > MaxNumberOfDaysPeriod {
					return fmt.Errorf("Max number of days for a given period exceeded:%.2f, max set to:%.2f",
						days, MaxNumberOfDaysPeriod)
				}
			}
		}
	}
	for i := 0; i < len(downloadConfig.MPRGovs); i++ {
		if downloadConfig.MPRGovs[i].RollDays > MaxRollDays {
			return fmt.Errorf("Invalid number of rolldays specified:%d, allowed:%d",
				downloadConfig.MPRGovs[i].RollDays,
				MaxRollDays)

		}
		//check the days
		if downloadConfig.MPRGovs[i].DateFrom != "" && downloadConfig.MPRGovs[i].DateTo != "" {
			if days, err := DaysBetween(downloadConfig.MPRGovs[i].DateFrom,
				downloadConfig.MPRGovs[i].DateTo); err != nil {
				return err
			} else {
				if days > MaxNumberOfDaysPeriod {
					return fmt.Errorf("Max number of days for a given period exceeded:%.2f, max set to:%.2f",
						days, MaxNumberOfDaysPeriod)
				}
			}
		}
		if strings.ToLower(downloadConfig.MPRGovs[i].Common.Format) == "pdf" {
			return fmt.Errorf("Unsupported format specified for MPRML Gov data, pdf not supported")
		}

	}
	for i := 0; i < len(downloadConfig.MPRPartners); i++ {
		if downloadConfig.MPRPartners[i].RollDays > MaxRollDays {
			return fmt.Errorf("Invalid number of rolldays specified:%d, allowed:%d",
				downloadConfig.MPRPartners[i].RollDays,
				MaxRollDays)
		}
		//check the days
		if downloadConfig.MPRPartners[i].DateFrom != "" && downloadConfig.MPRPartners[i].DateTo != "" {
			if days, err := DaysBetween(downloadConfig.MPRPartners[i].DateFrom,
				downloadConfig.MPRPartners[i].DateTo); err != nil {
				return err
			} else {
				if days > MaxNumberOfDaysPeriod {
					return fmt.Errorf("Max number of days for a given period exceeded:%.2f, max set to:%.2f",
						days, MaxNumberOfDaysPeriod)
				}
			}
		}
		if strings.ToLower(downloadConfig.MPRPartners[i].Common.Format) == "pdf" {
			return fmt.Errorf("Unsupported format specified for MPRML Partner data, pdf not supported")
		}
	}
	return nil
}

//ProcessAndRunDownload will process and query for a set of files as defined in the download config
//any errors will be tracked on a file basis and returned
func ProcessAndRunDownload(downloadConfig CloudDownload) []error {
	var err error
	var errList []error
	var token string
	var subscriptionKey, graphQLUrl, fileDownloadUrl string
	var fQueries []FileQuery

	if err = VerifyCloudDownloadConfig(downloadConfig); err != nil {
		return append(errList, err)
	}
	if token, err = Authenticate(); err != nil {
		zap.S().Errorf("Failed in getting token:%s", err.Error())
		return []error{err}
	}
	zap.S().Debugf("Got token:%s", token)
	//get the subscription key from the environment variables
	subscriptionKey = os.Getenv(AzureSubscriptionKeyEnvName)
	graphQLUrl = os.Getenv(AzureGraphUrlEnvName)
	fileDownloadUrl = os.Getenv(AzureFileDownloadUrlEnvName)
	if subscriptionKey == "" {
		errorMsg := fmt.Sprintf("Unable to find subscription key in environment variable:%s", AzureSubscriptionKeyEnvName)
		zap.S().Errorf(errorMsg)
		return []error{errors.New(errorMsg)}
	}
	if graphQLUrl == "" {
		errorMsg := fmt.Sprintf("Unable to locate environment variable for the graphqlurl:%s",
			AzureGraphUrlEnvName)
		zap.S().Errorf(errorMsg)
		return []error{errors.New(errorMsg)}
	}
	if fileDownloadUrl == "" {
		errorMsg := fmt.Sprintf("Unable to locate environment variable for the filedownloadurl:%s",
			AzureFileDownloadUrlEnvName)
		zap.S().Errorf(errorMsg)
		return []error{errors.New(errorMsg)}
	}
	//loop through all of the entities and process them one by one
	//first dprs
	zap.S().Debugf("Processing DPR's:%d", len(downloadConfig.DPRS))
	for i := 0; i < len(downloadConfig.DPRS); i++ {
		fQueries = append(fQueries,
			createDPR10Query(downloadConfig.DPRS[i], token))
		fQueries = append(fQueries,
			createDPR20Query(downloadConfig.DPRS[i], token))
	}
	zap.S().Debugf("Processing MPRML Govs:%d", len(downloadConfig.MPRGovs))
	for x := 0; x < len(downloadConfig.MPRGovs); x++ {
		//make sure that we do not have queries for pdfs here
		if strings.ToLower(downloadConfig.MPRGovs[x].Common.Format) == "pdf" {
			errList = append(errList, errors.New("PDFS are not supported for MPRML Government reports"))
		} else {
			fQueries = append(fQueries, createMPRMLGovQuery(downloadConfig.MPRGovs[x], token))
		}
	}
	zap.S().Debugf("Processing MPRML Partners:%d", len(downloadConfig.MPRPartners))
	for y := 0; y < len(downloadConfig.MPRPartners); y++ {
		if strings.ToLower(downloadConfig.MPRPartners[y].Common.Format) == "pdf" {
			errList = append(errList, errors.New("PDFS are not supported for MPRML Partner reports"))
		} else {
			fQueries = append(fQueries, createMPRMLPartnerQuery(downloadConfig.MPRPartners[y], token))
		}
	}
	zap.S().Debugf("Processing ddrmls:%s", len(downloadConfig.DDRMLS))
	for z := 0; z < len(downloadConfig.DDRMLS); z++ {
		fQueries = append(fQueries, createDDRMLQuery(downloadConfig.DDRMLS[z], token))
	}
	//now run all of the queries one by one
	for a := 0; a < len(fQueries); a++ {
		err = RunQueryAndDownloadFiles(fQueries[a], token, subscriptionKey, graphQLUrl,
			fileDownloadUrl)
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}

//RunQueryAndDownloadFiles will take an filequery object and run a graphql query for the specified files
//and using the file result it will try to download each file locally using the filedownloadurl and the file reference
func RunQueryAndDownloadFiles(fQuery FileQuery, token, subscriptionKey, graphQLUrl,
	fileDownloadUrl string) error {
	var dObj DataObject
	var err error
	var query []byte
	if fQuery.UseUploadedFrom {
		if query, err = BuildQueryForAssetUsingCreated(fQuery); err != nil {
			zap.S().Errorf("Failed in generation of asset query:%s", err.Error())
			return err
		}
	} else {

		if query, err = BuildQueryForAssetUsingPeriod(fQuery); err != nil {
			zap.S().Errorf("Failed in generation of asset query:%s", err.Error())
			return err
		}
	}
	zap.S().Infof("Running file query for, field:%s,fileType:%s, reportType:%s,timeRange:%s-%s, useUploadedFrom:%s",
		fQuery.Field,
		fQuery.FileType, fQuery.ReportType, fQuery.TimeFrom,
		fQuery.TimeTo, fQuery.UseUploadedFrom)
	zap.S().Debugf("Generated query:%s", query)
	if dObj, _, err = RunGraphQueryForFiles(token, graphQLUrl,
		subscriptionKey, query); err != nil {
		errorMsg := fmt.Sprintf("RunGraphQL query for files failed:%s", err.Error())
		zap.S().Errorf(errorMsg)
		return errors.New(errorMsg)
	}
	zap.S().Infof("Got number of files:%d", len(dObj.Files))
	if errorList := DownloadFiles(dObj.Files, fileDownloadUrl,
		token, subscriptionKey, strings.ToUpper(fQuery.FileType),
		fQuery.OutputLocation, fQuery.OutputPrefix, false); len(errorList) > 0 {
		for i := 0; i < len(errorList); i++ {
			zap.S().Errorf(errorList[i].Error())

		}
		return errors.New("Failed in download of cloud files, please check logs...")
	}
	return nil
}

//createDPR10Query creates a DPR 1.0 query object
func createDPR10Query(dprCnfg CloudProductionConfig, token string) FileQuery {

	fQuery := buildFileQuery(dprCnfg.RollDays, dprCnfg.DateFrom, dprCnfg.DateTo)
	fQuery.Field = dprCnfg.FieldName
	fQuery.ReportType = "DPR10"
	fQuery.FileType = dprCnfg.Common.Format
	fQuery.UseUploadedFrom = dprCnfg.UseUploadedFrom
	fQuery.OutputLocation = dprCnfg.Common.OutputFolder
	fQuery.OutputPrefix = dprCnfg.Common.FileOutputPrefix
	return fQuery

}

func createDPR20Query(dprCnfg CloudProductionConfig, token string) FileQuery {

	fQuery := buildFileQuery(dprCnfg.RollDays, dprCnfg.DateFrom, dprCnfg.DateTo)
	fQuery.Field = dprCnfg.FieldName
	fQuery.ReportType = "DPR20"
	fQuery.FileType = dprCnfg.Common.Format
	fQuery.UseUploadedFrom = dprCnfg.UseUploadedFrom
	fQuery.OutputLocation = dprCnfg.Common.OutputFolder
	fQuery.OutputPrefix = dprCnfg.Common.FileOutputPrefix
	return fQuery

}

func createMPRMLGovQuery(mpmrmlCnfg CloudProductionConfig, token string) FileQuery {

	fQuery := buildFileQuery(mpmrmlCnfg.RollDays, mpmrmlCnfg.DateFrom, mpmrmlCnfg.DateTo)
	fQuery.Field = mpmrmlCnfg.FieldName
	fQuery.ReportType = "MPRMLGov"
	fQuery.FileType = mpmrmlCnfg.Common.Format
	fQuery.UseUploadedFrom = mpmrmlCnfg.UseUploadedFrom
	fQuery.OutputLocation = mpmrmlCnfg.Common.OutputFolder
	fQuery.OutputPrefix = mpmrmlCnfg.Common.FileOutputPrefix
	return fQuery

}

func createMPRMLPartnerQuery(mpmrmlCnfg CloudProductionConfig, token string) FileQuery {

	fQuery := buildFileQuery(mpmrmlCnfg.RollDays, mpmrmlCnfg.DateFrom, mpmrmlCnfg.DateTo)
	fQuery.Field = mpmrmlCnfg.FieldName
	fQuery.ReportType = "MPRMLPartner"
	fQuery.FileType = mpmrmlCnfg.Common.Format
	fQuery.UseUploadedFrom = mpmrmlCnfg.UseUploadedFrom
	fQuery.OutputLocation = mpmrmlCnfg.Common.OutputFolder
	fQuery.OutputPrefix = mpmrmlCnfg.Common.FileOutputPrefix
	return fQuery

}

func createDDRMLQuery(ddrmlConfig CloudDDRMLConfig, token string) FileQuery {

	fQuery := buildFileQuery(ddrmlConfig.RollDays, ddrmlConfig.DateFrom, ddrmlConfig.DateTo)
	fQuery.ReportType = "DDRML"
	fQuery.FileType = ddrmlConfig.Common.Format
	fQuery.UseUploadedFrom = ddrmlConfig.UseUploadedFrom
	fQuery.OutputLocation = ddrmlConfig.Common.OutputFolder
	fQuery.OutputPrefix = ddrmlConfig.Common.FileOutputPrefix
	return fQuery

}

//BuildOutputPathForReportType will build the output folder path based on the report type
//DDRML data will be stored in outputFolder + name of wellbore and then the actual file
//the function will return an array of paths based on the source objects as especially ddrml files can contain several wellbores
func BuildOutputPathForReportType(fObj FileObject, filePrefix, outputFolder, format string) []string {
	var paths []string
	outputFileName := GenerateFileName(fObj, filePrefix, format)
	outputFolder = outputFolder + string(filepath.Separator)

	//if the report type is a ddrml report we need to generate a subfolder for the wellbore itself which is part of the file
	if strings.Contains(strings.ToLower(MapReportType(fObj.ReportType)), "ddrml") {
		//loop through all of the sources to separate out the data
		for i := 0; i < len(fObj.Sources); i++ {
			paths = append(paths, outputFolder+SafeEncodeNameForWinFiles(fObj.Sources[i].Name)+string(filepath.Separator)+outputFileName)
		}
	} else {
		paths = append(paths, outputFolder+outputFileName)
	}
	return paths
}

/*DownloadFilesToZip downloads a set of files using the given array of fileReference id
fileURL -> to download for from Azure:
token -> recieved from oauth2
subscription key -> to use API service
format -> the format to download (pdf or xml)
*/
func DownloadFilesToZip(fileReferences []FileObject, fileURL, token, subscriptionKey, format string) ([]byte, error) {
	var resp *resty.Response
	var err error
	var data []byte

	for i := 0; i < len(fileReferences); i++ {
		zap.S().Debugf("Downloading file with reference id:%s, format:%s, from url:%s",
			fileReferences[i].FileReference, format, fileURL)
	}

	//add the required headers
	if fileURL == "" {
		return data, errors.New("Missing fileurl...")
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Ocp-Apim-Subscription-Key"] = subscriptionKey
	//create the query params
	queryString := ""
	for i := 0; i < len(fileReferences); i++ {
		queryString = queryString + "&referenceIds=" + fileReferences[i].FileReference
	}
	queryString = "format=" + strings.ToLower(format) + queryString
	fmt.Printf("QUERYSTRING:%s\n", queryString)
	client := resty.New()
	client.SetTimeout(time.Duration(1 * time.Minute))
	if ce := zap.S().Desugar().Check(zap.DebugLevel, "debugging"); ce != nil {
		client.SetDebug(true)
	}
	if resp, err = client.R().
		SetQueryString(queryString).
		SetHeaders(headers).
		SetAuthToken(token).
		Get("https://apis.collabor8.no/test/api/Files/files"); err != nil {
		zap.S().Errorf("Failed in get of file, format:%s,error:%s",
			format, err.Error())
		return data, err
	}

	if resp.IsError() {
		//log the error
		errorMsg := fmt.Sprintf("Failed in get of file server responded with http error > 400, format:%s,httpStatusCode:%d,httpStatus:%s,body:%s",
			format, resp.StatusCode(), resp.Status(), string(resp.Body()))
		zap.S().Errorf(errorMsg)
		return data, errors.New(errorMsg)
	}
	//return the body, it is already closed by the resty library and copied to new array

	return resp.Body(), nil
}

//DownloadFile downloads a single file using the given fileReference id,
// fileURL for Azure, token from oauth2, subscription key for service and format to download (pdf or xml)
func DownloadFile(fileReference, fileURL, token, subscriptionKey, format string) ([]byte, error) {
	var resp *resty.Response
	var err error
	var data []byte
	zap.S().Debugf("Downloading file with reference id:%s, format:%s, from url:%s",
		fileReference, format, fileURL)
	//add the required headers
	if fileURL == "" {
		return data, errors.New("Missing fileurl...")
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Ocp-Apim-Subscription-Key"] = subscriptionKey
	client := resty.New()
	client.SetTimeout(time.Duration(1 * time.Minute))
	if ce := zap.S().Desugar().Check(zap.DebugLevel, "debugging"); ce != nil {
		client.SetDebug(true)
	}
	if resp, err = client.R().
		SetQueryParams(map[string]string{
			"format": strings.ToLower(format),
		}).
		SetHeaders(headers).
		SetAuthToken(token).
		Get(fileURL + "/" + fileReference); err != nil {
		zap.S().Errorf("Failed in get of file, referenceId:%s, format:%s,error:%s",
			fileReference, format, err.Error())
		return data, err
	}

	if resp.IsError() {
		//log the error
		errorMsg := fmt.Sprintf("Failed in get of file server responded with http error > 400, referenceId:%s, format:%s,httpStatusCode:%d,httpStatus:%s,body:%s",
			fileReference, format, resp.StatusCode(), resp.Status(), string(resp.Body()))
		zap.S().Errorf(errorMsg)
		return data, errors.New(errorMsg)
	}
	//return the body, it is already closed by the resty library and copied to new array

	return resp.Body(), nil
}

func buildFileQuery(rolldays int, dateFrom, dateTo string) FileQuery {
	fQuery := FileQuery{}
	if rolldays != common.Rolldays_default && rolldays != 0 {
		zap.S().Debugf("Rolldays is not set to the default:%d, but:%d will generate date range based on this", common.Rolldays_default, rolldays)
		start, end := common.RollDays(rolldays)

		fQuery.TimeFrom = common.FormatTime2QueryDayString(start)
		fQuery.TimeTo = common.FormatTime2QueryDayString(end)
		zap.S().Debugf("Generated date range from rolldays, start:%s ,end:%s", fQuery.TimeFrom, fQuery.TimeTo)
	} else {
		//not rolling days using fixed setup
		fQuery.TimeFrom = dateFrom
		fQuery.TimeTo = dateTo
		zap.S().Debugf("Using date range, start:%s, end:%s", dateFrom, dateTo)
	}
	return fQuery
}

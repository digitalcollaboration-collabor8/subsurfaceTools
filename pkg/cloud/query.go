package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type GraphQuery struct {
	Query         string `json:"query"`
	OperationName string `json:"operationName,omitempty"`
	Variables     string `json:"variables,omitempty"`
}

type FileQueryConfig struct {
	AuthCnfg        AuthConfig
	GraphURL        string
	FileDownloadUrl string
	Queries         []FileQuery
}

type FileQuery struct {
	TimeFrom        string
	TimeTo          string
	Field           string
	FileType        string
	ReportType      string
	UseUploadedFrom bool
	OutputLocation  string
	OutputPrefix    string
}

type FileGraphResult struct {
	Data DataObject `json:"data"`
}

type DataObject struct {
	Files  []FileObject `json:"files"`
	Errors []DataErrors `json:"errors"`
}

type DataErrors struct {
	Message string `json:"message"`
}

type FileObject struct {
	FileName      string       `json:"fileName"`
	FileReference string       `json:"fileReferenceId"`
	Created       string       `json:"created"`
	MetaData      FileMetaData `json:"metadata"`
	ReportType    int          `json:"reportType"`
	Sources       []DataSource `json:"sources"`
}

type DataSource struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	NamingSystem string `json:"namingSystem"`
}

type FileMetaData struct {
	FileType     string `json:"fileType"`
	PeriodEnd    string `json:"periodEnd"`
	PeriodStart  string `json:"periodStart"`
	ReportId     string `json:"reportId"`
	ReportStatus string `json:"reportStatus"`
}

var queryForXMLFilesForAssetUsingCreated = `query{
	files(created_after: "{{.TimeFrom}}",
		created_before: "{{.TimeTo}}",
		{{if ne .Field ""}}field:"{{.Field}}",{{end}}fileType:XML{{if ne .ReportType ""}},report_type:{{.ReportType}}{{end}}){
		  fileName
		  fileReferenceId
		  created
		  metadata{
			fileType
			periodEnd
			periodStart
			reportId
			reportStatus
			
		  }
		  reportType
		  sources{
			kind
			name
			namingSystem
		  }
		}
	  }`

var queryForXMLFilesForAssetUsingPeriod = `query{
		files(period_start: "{{.TimeFrom}}",
			period_end: "{{.TimeTo}}",
			{{if ne .Field ""}}field:"{{.Field}}",{{end}}fileType:XML{{if ne .ReportType ""}},report_type:{{.ReportType}}{{end}}){
			  fileName
			  fileReferenceId
			  created
			  metadata{
				fileType
				periodEnd
				periodStart
				reportId
				reportStatus
				
			  }
			  reportType
			  sources{
				kind
				name
				namingSystem
			  }
			}
		  }`

//BuildQueryForAssetUsingCreated builds a qraphql query for file using created to/from date
//if the report type is either ddrml or dpr20 and the format is pdf it will ask for xml files instead
//as the minio storage just stores xml files for these report types and the pdf report is generated on the fly.
func BuildQueryForAssetUsingCreated(fQuery FileQuery) ([]byte, error) {

	var tpl bytes.Buffer

	if strings.ToLower(fQuery.ReportType) == "dpr20" || strings.ToLower(fQuery.ReportType) == "ddrml" {
		//setting thje filetype to xml as it is a dpr20 or ddrml report
		fQuery.FileType = "XML"

	}
	zap.S().Debugf("Building asset query:%s,%s,%s,%s", fQuery.Field, fQuery.FileType,
		fQuery.TimeFrom, fQuery.TimeTo)
	tmpl, err := template.New("assetsQuery").Parse(queryForXMLFilesForAssetUsingCreated)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&tpl, fQuery)
	if err != nil {
		return nil, err
	}
	zap.S().Debugf("Generated query:%s", tpl.String())

	return tpl.Bytes(), nil
}

//BuildQueryForAssetUsingPeriod builds a qraphql query for file using period to/from date
//if the report type is either ddrml or dpr20 and the format is pdf it will ask for xml files instead
//as the minio storage just stores xml files for these report types and the pdf report is generated on the fly.
func BuildQueryForAssetUsingPeriod(fQuery FileQuery) ([]byte, error) {
	var tpl bytes.Buffer
	if strings.ToLower(fQuery.ReportType) == "dpr20" || strings.ToLower(fQuery.ReportType) == "ddrml" {
		//setting thje filetype to xml as it is a dpr20 or ddrml report
		fQuery.FileType = "XML"

	}
	zap.S().Debugf("Building asset query:%s,%s,%s,%s", fQuery.Field, fQuery.FileType,
		fQuery.TimeFrom, fQuery.TimeTo)
	tmpl, err := template.New("assetsQuery").Parse(queryForXMLFilesForAssetUsingPeriod)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&tpl, fQuery)
	if err != nil {
		return nil, err
	}
	zap.S().Debugf("Generated query:%s", tpl.String())

	return tpl.Bytes(), nil
}

func RunGraphQueryForFiles(token, url, subscriptionKey string, query []byte) (DataObject, interface{}, error) {
	var resp *resty.Response
	var dObject DataObject
	var fResult FileGraphResult
	headers := make(map[string]string)

	var err error
	//create the graphqlQueryObject
	queryObj := GraphQuery{
		Query: string(query),
	}
	payload, err := json.Marshal(queryObj)
	if err != nil {
		zap.S().Errorf("Failed in marshal of graphql object:%s", err.Error())
		return dObject, nil, err
	}
	//create the headers
	//fmt.Printf("GRAPH_TOKEN:%s\n", token)
	headers["Content-Type"] = "application/json"
	headers["Ocp-Apim-Subscription-Key"] = subscriptionKey

	//build the client
	client := resty.New()
	if resp, err = client.R().
		SetHeaders(headers).
		SetAuthToken(token).
		SetBody(payload).
		Post(url); err != nil {
		zap.S().Errorf("Failed in rest request:%s", err.Error())
		return dObject, nil, err

	}

	if resp.IsError() {

		return dObject, resp.Error(), fmt.Errorf("Code:%d,Status:%s,Body:%s",
			resp.StatusCode(), resp.Status(), string(resp.Body()))
	} else {
		zap.S().Debugf("Got response from server:%s", string(resp.Body()))
		//unmarshal it
		if err = json.Unmarshal(resp.Body(), &fResult); err != nil {
			zap.S().Errorf("Failed in unmarshalling of response error:%s,got:%s", err.Error(), (resp.Body()))
			return dObject, nil, err
		}

		return fResult.Data, resp.Error(), err
	}
}

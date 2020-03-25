package cloud

import (
	"os"
	"strings"
	"testing"
)

func TestFailDownloadWithInvalidToken(t *testing.T) {
	token := "XXXXX"
	if result, err := DownloadFile("XXXX",
		os.Getenv("AzureFileDownloadUrl"), token, os.Getenv(AzureSubscriptionKeyEnvName), "xml"); err == nil {
		t.Errorf("Download of file should fail without subscription key:got body,%s", string(result))
	} else {
		if !strings.Contains(err.Error(), "401") {
			//fail if we do not get an unauthorized back
			t.Errorf("Response back doens not contain an 401 unauthorized...")
		}
	}
}

func TestFailDownloadWithInvalidFileId(t *testing.T) {
	var token string
	var err error
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in getting token, got error back instead of token:%s", err.Error())
	}
	if result, err := DownloadFile("XXXX",
		os.Getenv("AzureFileDownloadUrl"), token, os.Getenv(AzureSubscriptionKeyEnvName), "xml"); err == nil {
		t.Errorf("Download of file should fail with 400 bad request:got body,%s", string(result))
	} else {
		t.Logf("error:%s", err.Error())
		if !strings.Contains(err.Error(), "400") {
			//fail if we do not get an unauthorized back
			t.Errorf("Response back does not contain an 400 unauthorized...")
		}
	}
}

func TestFailDownloadWithInvalidFormat(t *testing.T) {
	var token string
	var err error
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in getting token, got error back instead of token:%s", err.Error())
	}
	if result, err := DownloadFile(AzureFileIDToTest,
		os.Getenv("AzureFileDownloadUrl"), token, os.Getenv(AzureSubscriptionKeyEnvName), "XML"); err == nil {
		t.Errorf("Download of file should fail with 400 bad request:got body,%s", string(result))
	} else {
		t.Logf("error:%s", err.Error())
		if !strings.Contains(err.Error(), "400") {
			//fail if we do not get an unauthorized back
			t.Errorf("Response back does not contain an 400 unauthorized...")
		}
	}
}

func TestDownloadXMLFile(t *testing.T) {
	var token string
	var err error
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in getting token, got error back instead of token:%s", err.Error())
	}
	if result, err := DownloadFile(AzureFileIDToTest,
		os.Getenv("AzureFileDownloadUrl"), token, os.Getenv(AzureSubscriptionKeyEnvName), "xml"); err != nil {
		t.Errorf("Download of file should not fail:%s", err.Error())
	} else {
		//check that the body contains xml
		if !strings.Contains(string(result), "xml") {
			t.Errorf("Expected xml back but got something else:%s", string(result))
		}

	}
}

func TestDownloadPDFFile(t *testing.T) {
	var token string
	var err error
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in getting token, got error back instead of token:%s", err.Error())
	}
	if result, err := DownloadFile(AzureFileIDToTest,
		os.Getenv("AzureFileDownloadUrl"), token, os.Getenv(AzureSubscriptionKeyEnvName), "pdf"); err != nil {
		t.Errorf("Download of file should not fail:%s", err.Error())
	} else {
		//check that the body contains xml
		if !strings.Contains(string(result), "PDF") {
			t.Errorf("Expected pdf back but got something else:%s", "dd")
		}

	}
}

func TestQueryAndDownloadXMLFiles(t *testing.T) {
	var token string
	var query []byte
	var err error
	var dObj DataObject
	fQuery := FileQuery{
		TimeFrom: "2020-02-24T23:00:00.000Z",
		TimeTo:   "2020-02-29T19:00:00.000Z",
		Field:    "ÅSGARD",
		FileType: "XML",
	}
	if query, err = BuildQueryForAssetUsingCreated(fQuery); err != nil {
		t.Errorf("Failed in generation of asset query:%s", err.Error())
		return
	} /*else {
		t.Log(string(query))
	}*/
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in test of authentication in run query for files, got error back instead of token:%s", err.Error())
	}
	if dObj, _, err = RunGraphQueryForFiles(token, os.Getenv(AzureGraphUrlEnvName), os.Getenv(AzureSubscriptionKeyEnvName), query); err != nil {
		t.Errorf("Resty post for files should not fail, failed with:%s", err.Error())
		return
	}
	t.Logf("Got number of files:%d", len(dObj.Files))
	if errorList := DownloadFiles(dObj.Files, os.Getenv("AzureFileDownloadUrl"),
		token, os.Getenv(AzureSubscriptionKeyEnvName), "XML", LocalStorageLocation, "ASGARD", false); len(errorList) > 0 {
		//download should not fail check errors
		for i := 0; i < len(errorList); i++ {
			t.Error(err.Error())
		}
		return
	}

}

func TestQueryAndDownloadPDFFiles(t *testing.T) {
	var token string
	var query []byte
	var err error
	var dObj DataObject
	fQuery := FileQuery{
		TimeFrom: "2020-02-24T23:00:00.000Z",
		TimeTo:   "2020-02-29T19:00:00.000Z",
		Field:    "ÅSGARD",
		FileType: "PDF",
	}
	if query, err = BuildQueryForAssetUsingCreated(fQuery); err != nil {
		t.Errorf("Failed in generation of asset query:%s", err.Error())
		return
	} /*else {
		t.Log(string(query))
	}*/
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in test of authentication in run query for files, got error back instead of token:%s", err.Error())
	}
	if dObj, _, err = RunGraphQueryForFiles(token,
		os.Getenv(AzureGraphUrlEnvName),
		os.Getenv(AzureSubscriptionKeyEnvName),
		query); err != nil {
		t.Errorf("Resty post for files should not fail, failed with:%s", err.Error())
		return
	}
	t.Logf("Got number of files:%d", len(dObj.Files))
	if errorList := DownloadFiles(dObj.Files,
		os.Getenv("AzureFileDownloadUrl"),
		token, os.Getenv(AzureSubscriptionKeyEnvName),
		"PDF", LocalStorageLocation, "ASGARD", false); len(errorList) > 0 {
		//download should not fail check errors
		for i := 0; i < len(errorList); i++ {
			t.Error(errorList[i].Error())
		}
		return
	}

}

func TestReplaceIllegalChars(t *testing.T) {
	expected := "A________B"
	string2Replace := `A<>:/\|?*B`
	string2Replace = SafeEncodeNameForWinFiles(string2Replace)
	if string2Replace != expected {
		t.Errorf("Replaced string:%s does not match expected:%s", string2Replace, expected)
	}
}

func TestGenerateFileNameDDRML(t *testing.T) {
	fObj := FileObject{
		FileName:      "DDR for NO 7120_8-L-4 H 2020-03-02.xml",
		FileReference: "12345",

		ReportType: 3,
		MetaData: FileMetaData{
			FileType:    "XML",
			PeriodStart: "2019-01-19",
			PeriodEnd:   "2019-01-20",
		},
	}
	fileName := GenerateFileName(fObj, "DDRML_", "PDF")
	if fileName != "DDRML__DDRML__2019-01-19_2019-01-20_12345_DDR for NO 7120_8-L-4 H 2020-03-02.pdf" {
		t.Errorf("Invalid filename:%s", fileName)
	}
	t.Logf("Filename:%s", fileName)
}

func TestQueryAndDownloadDDRMLPDFFiles(t *testing.T) {
	var token string
	var query []byte
	var err error
	var dObj DataObject
	fQuery := FileQuery{
		TimeFrom:   "2020-02-11T23:00:00.000Z",
		TimeTo:     "2020-03-13T19:00:00.000Z",
		FileType:   "PDF",
		ReportType: "DDRML",
	}
	if query, err = BuildQueryForAssetUsingCreated(fQuery); err != nil {
		t.Errorf("Failed in generation of asset query:%s", err.Error())
		return
	} /*else {
		t.Log(string(query))
	}*/
	//need to get the token first
	if token, err = GetValidToken(); err != nil {
		t.Errorf("Failed in test of authentication in run query for files, got error back instead of token:%s", err.Error())
	}
	if dObj, _, err = RunGraphQueryForFiles(token, os.Getenv(AzureGraphUrlEnvName), os.Getenv(AzureSubscriptionKeyEnvName), query); err != nil {
		t.Errorf("Resty post for files should not fail, failed with:%s", err.Error())
		return
	}
	t.Logf("Got number of files:%d", len(dObj.Files))
	if errorList := DownloadFiles(dObj.Files, os.Getenv("AzureFileDownloadUrl"),
		token, os.Getenv(AzureSubscriptionKeyEnvName), "PDF", LocalStorageLocation, "ASGARD", false); len(errorList) > 0 {
		//download should not fail check errors
		for i := 0; i < len(errorList); i++ {
			t.Error(errorList[i].Error())
		}
		return
	}

}

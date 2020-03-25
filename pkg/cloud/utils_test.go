package cloud

import (
	"testing"
)

var AzureFileIDToTest = "f58264b8-c905-4c4e-8a50-8004ed9c0243"
var LocalStorageLocation = "../../test/results"

func TestCalculateDaysBetweenDates(t *testing.T) {
	if days, err := DaysBetween("2019-01-01", "2019-01-31"); err != nil {
		t.Errorf("Failed in calculating days between:%s", err.Error())
	} else {
		t.Logf("Days between:%.2f", days)
	}
	if days, err := DaysBetween("2019-01-01", "2019-03-31"); err != nil {
		t.Errorf("Failed in calculating days between:%s", err.Error())
	} else {
		t.Logf("Days between:%.2f", days)
	}
}

func TestConvertRFC3339StrToOtherFormat(t *testing.T) {
	timeStampStr := "2020-03-01T10:44:51.526Z"
	format := "2006-01-02T15_04_05"
	if timeObj, err := StringRFC3339ToTime(timeStampStr); err != nil {
		t.Errorf("Failed in converting rfc3339 string to timestamp and to time string with format:%s",
			err.Error())
	} else {
		t.Logf("convertedTime:%s", TimeToStr(timeObj, format))
	}
}

func TestSuccessAuthenticate(t *testing.T) {

	if token, err := Authenticate(); err != nil {
		t.Errorf("Failed in test of authentication, got error back instead of token:%s", err.Error())
	} else {
		t.Logf("Got token:%s", token)
	}

}

func TestFailureAuthenticate(t *testing.T) {

	if _, err := Authenticate(); err == nil {
		t.Errorf("Failed in test of authentication, should get an error back but got nothing")
	} else {
		t.Logf("Got error:%s", err.Error())
	}

}

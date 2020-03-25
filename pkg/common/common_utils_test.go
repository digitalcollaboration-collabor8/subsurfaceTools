/*
* @Author: magsv
* @Date:   2018-08-07 12:58:48
* @Last Modified by:   magsv
* @Last Modified time: 2018-08-07 14:54:52
 */
package common

import (
	"testing"
	"time"
)

const timeVariant_NoTimeZoneDST = "2018-06-08T00:00:00"
const timeVariant_NoTimeZoneWinter = "2018-01-08T00:00:00"

const timeVariant_UTC = "2018-06-08T00:00:00Z"
const timeVariant_CET_DST = "2018-06-08T00:00:00+02:00"
const timeVariant_CET_WINTER = "2018-01-17T00:00:00+01:00"

func TestParseTimeInUTC(t *testing.T) {
	//var parse time.Time
	var err error
	if _, err = time.Parse(xsdDateTimeLayout, timeVariant_UTC); err != nil {
		t.Errorf("Failed:%s", err.Error())
	}
}

func TestParseNoTimezoneInDST(t *testing.T) {
	var parse time.Time
	var err error

	verifyString := "2018-06-08T00:00:00+02:00"
	if parse, err = time.ParseInLocation(xsdDateTimeLayoutNoTimeZone, timeVariant_NoTimeZoneDST, time.Now().Location()); err != nil {
		t.Errorf("Failed in parsing time zone object:%s", err.Error())
	}
	if TimeToString(parse, xsdDateTimeLayout) != verifyString {
		t.Errorf("Failed in parsing timezone in DST, expected:%s, got:%s", verifyString, TimeToString(parse, xsdDateTimeLayout))
	}
}
func TestParseNoTimezoneInWinter(t *testing.T) {
	var parse time.Time
	var err error
	verifyString := "2018-01-08T00:00:00+01:00"
	if parse, err = time.ParseInLocation(xsdDateTimeLayoutNoTimeZone, timeVariant_NoTimeZoneWinter, time.Now().Location()); err != nil {
		t.Errorf("Failed in parsing time zone object:%s", err.Error())
	}
	if TimeToString(parse, xsdDateTimeLayout) != verifyString {
		t.Errorf("Failed in parsing timezone in winter, expected:%s, got:%s", verifyString, TimeToString(parse, xsdDateTimeLayout))
	}
}

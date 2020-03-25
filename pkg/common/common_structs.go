/*
* @Author: magsv
* @Date:   2017-11-06 15:42:00
* @Last Modified by:   magsv
* @Last Modified time: 2018-08-07 13:14:46
 */

package common

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

const Rollmonths_default = -999
const Rolldays_default = -999

const xsdDateTimeLayout = time.RFC3339
const xsdDateLayout = "2006-01-02"
const NullFloatValue = -999.99
const xsdDateTimeLayoutNoTimeZone = "2006-01-02T15:04:05"

type xsdDouble struct {
	float64
}

//unmarshals a xsd date element in the form of 2006-01-02 (YYYY-MM-DD)
type xsdDate struct {
	time.Time
}

//unmarshals an xsd datetime element in the form of 2017-05-12T09:18:06
type xsdDateTime struct {
	time.Time
}

//unmarshal function to handle xsd double parsing,
// xsd double e.g. according to IEEE standard allows for lexical space that needs to be handled
func (c *xsdDouble) UnmarshalText(text []byte) error {
	var v string
	var err error
	var number float64
	v = string(text)

	parse := strings.Trim(v, " ")
	if number, err = strconv.ParseFloat(parse, 64); err != nil {
		return err
	}
	*c = xsdDouble{number}
	return err
}

//unmarshal function to handle xsd date time parsing
func (c *xsdDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, _ := time.Parse(xsdDateTimeLayout, v)
	*c = xsdDateTime{parse}
	return nil
}

func (c *xsdDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, _ := time.Parse(xsdDateLayout, v)
	*c = xsdDate{parse}
	return nil
}

const datelayout_out = "2006-01-02T15_04_05"
const datelayout_out_csv = "2006-01-02 15:04:05"

const daylayout_query = "2006-01-02"
const monthlayout_query = "01"
const yearlayout_query = "2006"

type DataSet struct {
	Name        string
	HeadersName []string
	Rows        []RowData
}

type RowData struct {
	Columns []ColumnData
}

type ColumnData struct {
	IntVal        int       `json:",omitempty"`
	FloatVal      float64   `json:",omitempty"`
	StrVal        string    `json:",omitempty"`
	TimeValue     time.Time `json:",omitempty"`
	IsEmptyColumn bool      `json:",omitempty"`
	IsInt         bool      `json:",omitempty"`
	IsFloat       bool      `json:",omitempty"`
	IsStr         bool      `json:",omitempty"`
	IsTime        bool      `json:",omitempty"`
}

type DataSetJson struct {
	Name        string
	HeadersName []string
	Rows        []RowDataJson
}

type RowDataJson struct {
	Columns []ColumnDataJson
}

type ColumnDataJson struct {
	Name  string
	Value interface{}
}

type Config struct {
	ServerUrl       string
	User            string
	Password        string
	OutputFile      string
	SparqlQueryFile string
	Params          []ParamType
}

type ParamType struct {
	Name  string
	Value string
}

func (row *RowData) AddIntValue(value int) {
	cData := ColumnData{}
	cData.IntVal = value
	cData.IsInt = true
	row.Columns = append(row.Columns, cData)
}

func (row *RowData) AddEmptyColumn() {
	cData := ColumnData{}
	cData.IsEmptyColumn = true
	row.Columns = append(row.Columns, cData)
}

func (row *RowData) AddStrValue(value string) {
	cData := ColumnData{}
	cData.StrVal = value
	cData.IsStr = true
	row.Columns = append(row.Columns, cData)
}

func (row *RowData) AddFloatValue(value float64) {
	cData := ColumnData{}
	cData.FloatVal = value
	cData.IsFloat = true
	row.Columns = append(row.Columns, cData)
}

func (row *RowData) AddTimeValue(value time.Time) {
	cData := ColumnData{}
	cData.TimeValue = value
	cData.IsTime = true
	row.Columns = append(row.Columns, cData)
}

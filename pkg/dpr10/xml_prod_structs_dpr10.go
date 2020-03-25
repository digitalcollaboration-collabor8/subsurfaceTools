/*
* @Author: magsv
* @Date:   2018-02-22 10:54:45
* @Last Modified by:   magsv
* @Last Modified time: 2018-08-07 15:39:54
 */
package dpr10

import (
	"encoding/xml"

	"strconv"
	"strings"
	"time"
)

const xsdDateTimeLayout = time.RFC3339

const xsdDateTimeLayoutNoTimeZone = "2006-01-02T15:04:05"
const xsdDateLayout = "2006-01-02"

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
	var parse time.Time
	var err error
	d.DecodeElement(&v, &start)
	if parse, err = time.Parse(xsdDateTimeLayout, v); err != nil {
		//try to parse it without timezone..Read it in using local timezone, just assume that if it is undetermined it should be local timezone
		parse, err = time.ParseInLocation(xsdDateTimeLayoutNoTimeZone, v, time.Now().Location())
	}

	*c = xsdDateTime{parse}
	return nil
}

func (c *xsdDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, _ := time.ParseInLocation(xsdDateLayout, v, time.Now().Location())
	*c = xsdDate{parse}
	return nil
}

type ProcessingData struct {
	UUid           string //uniqueue identifier for this dataset from a given report
	FileName       string
	FilePath       string
	DocumentName   string
	ReadTime       time.Time
	ReportStart    time.Time
	ReportEnd      time.Time
	ReportTimStart time.Time
	ReportTimEnd   time.Time
}

type FacilityData struct {
	DataIdentification ProcessingData
	Installation       Installation
}

type Objects struct {
	XMLName            xml.Name `xml:"objects" json:"-"`
	DataIdentification ProcessingData
	DocumentInfo       DocumentInfo
	Context            ReportContext
	ProdObjects        []ProdVolume `xml:"object"`
}

type ReportContext struct {
	XMLName       xml.Name    `xml:"context"`
	Kind          string      `xml:"kind"`
	Title         NameElement `xml:"title"`
	Month         string      `xml:"month"`
	Year          string      `xml:"year"`
	ReportVersion float64     `xml:"reportVersion"`
	ReportStatus  string      `xml:"reportStatus"`
	Installation  Installation
}

type DocumentInfo struct {
	XMLName          xml.Name    `xml:"documentInfo" json:"-"`
	DocumentName     NameElement `xml:"DocumentName"`
	DocumentDate     xsdDateTime `xml:"DocumentDate"`
	FileCreationInfo FileCreationInfo
	AuditTrail       AuditTrail
}

type FileCreationInfo struct {
	XMLName          xml.Name    `xml:"FileCreationInformation" json:"-"`
	FileCreationDate xsdDateTime `xml:"FileCreationDate"`
	SoftwareName     string      `xml:"SoftwareName"`
	FileCreator      string      `xml:"FileCreator"`
	Comment          string      `xml:"Comment"`
}

type AuditTrail struct {
	XMLName xml.Name        `xml:"AuditTrail" json:"-"`
	Events  []DocumentEvent `xml:"Event"`
}

type DocumentEvent struct {
	EventDate        xsdDateTime `xml:"EventDate"`
	ResponsibleParty string      `xml:"ResponsibleParty"`
	Comment          string      `xml:"Comment"`
}

type ProdVolume struct {
	Name         string `xml:"name"`
	Installation Installation
	DateStart    xsdDate    `xml:"dateStart"`
	DateEnd      xsdDate    `xml:"dateEnd"`
	Facilities   []Facility `xml:"facility"`
}

type Installation struct {
	XMLName      xml.Name `xml:"installation" json:"-"`
	Kind         string   `xml:"kind,attr"`
	UidRef       string   `xml:"uidRef,attr"`
	Name         string   `xml:",chardata"`
	NamingSystem string   `xml:"namingSystem,attr"`
}

type Facility struct {
	Uid                string `xml:"uid,attr"`
	DataIdentification FacilityData
	Name               FacilityName
	FacilityParent1    FacilityParent1
	FacilityParent2    FacilityParent2
	ContextFacility    ContextFacility
	WellStatus         string `xml:"wellStatus" json:"omitEmpty"`
	WellProducing      string `xml:"wellProducing"`
	WellInjecting      string `xml:"wellInjecting"`
	OperationTime      Value  `xml:"operationTime" json:"omitEmpty"`
	Flow               []Flow `xml:"flow"`
}

type FacilityParent1 struct {
	XMLName      xml.Name `xml:"facilityParent1" json:"-"`
	Kind         string   `xml:"kind,attr"`
	UidRef       string   `xml:"uidRef,attr"`
	Name         string   `xml:",chardata"`
	NamingSystem string   `xml:"namingSystem,attr"`
}

type FacilityParent2 struct {
	XMLName      xml.Name `xml:"facilityParent2" json:"-"`
	Kind         string   `xml:"kind,attr"`
	UidRef       string   `xml:"uidRef,attr"`
	Name         string   `xml:",chardata"`
	NamingSystem string   `xml:"namingSystem,attr"`
}

type FacilityName struct {
	XMLName      xml.Name `xml:"name" json:"-"`
	Kind         string   `xml:"kind,attr"`
	UidRef       string   `xml:"uidRef,attr"`
	Name         string   `xml:",chardata"`
	NamingSystem string   `xml:"namingSystem,attr"`
}

type Flow struct {
	UID       string `xml:"uid,attr"`
	Name      string `xml:"name" json:"omitEmpty"`
	Kind      string `xml:"kind"`
	Qualifier string `xml:"qualifier"`

	PortDiff        FlowPortDiff  `xml:"portDiff"`
	Temp            Value         `xml:"temp"`
	Pres            Value         `xml:"pres"`
	GoR             Value         `xml:"gor"`
	WaterConcVol    Value         `xml:"waterConcVol"`
	Product         []FlowProduct `xml:"product"`
	InternalFlowUid string
}

type FlowPortDiff struct {
	Port          string `xml:"port"`
	ChokeRelative Value  `xml:"chokeRelative"`
}

type FlowProduct struct {
	Name                string             `xml:"name"`
	Kind                string             `xml:"kind"`
	Rvp                 Value              `xml:"rvp"`
	BsW                 Value              `xml:"bsw"`
	DensityStd          Value              `xml:"densityStd"`
	WobbeIndex          Value              `xml:"wobbeIndex"`
	GrossCalorificValue Value              `xml:"grossCalorificValueStd"`
	Period              []FlowPeriod       `xml:"period"`
	ComponentContent    []ComponentContent `xml:componentContent"`
	InternalFlowUid     string
}

type ComponentContent struct {
	Kind          string `xml:"kind"`
	Concentration Value  `xml:"concentration"`
}

type FlowPeriod struct {
	Kind            string       `xml:"kind"`
	DateStart       xsdDate      `xml:"dateStart"`
	DateEnd         xsdDate      `xml:"dateEnd"`
	DTimStart       xsdDateTime  `xml:"dTimStart"`
	DTimEnd         xsdDateTime  `xml:"dTimEnd"`
	BalanceSets     []BalanceSet `xml:"balanceSet"`
	VolumeOnly      Value        `xml:"volume"`
	Volume          VolumeValue
	VolumeStd       Value `xml:"volumeStd"`
	Temp            Value `xml:"temp"`
	Pres            Value `xml:"pres"`
	PortDiff        PortDiff
	Mass            Value `xml:"mass"`
	Work            Value `xml:"work"`
	Density         DensityValue
	InternalFlowUid string
}

type PortDiff struct {
	XMLName       xml.Name `xml:"portDiff"`
	Port          string   `xml:"port"`
	PresDiff      Value    `xml:"presDiff"`
	TempDiff      Value    `xml:"tempDiff"`
	ChokeRelative Value    `xml:"chokeRelative"`
	ChokeSize     Value    `xml:"chokeSize"`
}

type VolumeValue struct {
	XMLName  xml.Name `xml:"volumeValue" json:"-"`
	Volume   Value    `xml:"volume"`
	Temp     Value    `xml:"temp"`
	Pressure Value    `xml:"pres"`
}

type DensityValue struct {
	XMLName  xml.Name `xml:"densityValue" json:"-"`
	Density  Value    `xml:"density"`
	Temp     Value    `xml:"temp"`
	Pressure Value    `xml:"pres"`
}

type BalanceSet struct {
	Kind           string `xml:"kind"`
	CargoNumber    string `xml:"cargoNumber"`
	Destination    Destination
	Event          Event
	Volume         VolumeValue
	Mass           Value `xml:"mass"`
	Density        DensityValue
	BalanceDetails []BalanceDetail `xml:"balanceDetail"`
}

type BalanceDetail struct {
	Owner   string `xml:"owner"`
	Share   Value  `xml:"share"`
	Volume  VolumeValue
	Mass    Value `xml:"mass"`
	Density DensityValue
}

type Destination struct {
	XMLName xml.Name `xml:"destination" json:"-"`
	Name    string   `xml:"name"`
	Country string   `xml:"country"`
}

type Event struct {
	XMLName xml.Name `xml:"event" json:"-"`
	Date    xsdDate  `xml:"date"`
	Kind    string   `xml:"kind"`
}

type Value struct {
	Value float64 `xml:",chardata"`
	Uom   string  `xml:"uom,attr"`
}

type DoubleValue struct {
	Value xsdDouble `xml:",chardata"`
	Uom   string    `xml:"uom,attr"`
}

type NameElement struct {
	Txt          string `xml:",chardata"`
	NamingSystem string `xml:"namingSystem,attr"`
}

//start DPR 1.0 mappings

type WITSMLComposite struct {
	XMLName            xml.Name `xml:"WITSMLComposite"`
	DataIdentification ProcessingData
	DocumentInfo       DocumentInfo
	ProdOperationSet   ProdOperationSet
	ProdVolumeSet      ProdVolumeSet
	Wellset            WellSet
}

type ProdOperationSet struct {
	XMLName       xml.Name        `xml:"productionOperationSet"`
	ProdOperation []ProdOperation `xml:"productionOperation"`
}

type ProdOperation struct {
	Name            string `xml:"name"`
	Installation    Installation
	ContextFacility ContextFacility
	Kind            string `xml:"kind"`
	PeriodKind      string `xml:"periodKind"`

	DateStart          xsdDate              `xml:"dateStart"`
	DateEnd            xsdDate              `xml:"dateEnd"`
	DTimStart          xsdDateTime          `xml:"dTimStart"`
	DTimEnd            xsdDateTime          `xml:"dTimEnd"`
	InstallationReport []InstallationReport `xml:"installationReport"`
}

type ProdVolumeSet struct {
	XMLName      xml.Name `xml:"productVolumeSet"`
	Name         string   `xml:"name"`
	Installation string   `xml:"installation"`

	ProdVolumes []ProdVolume `xml:"productVolume"`
}

type ContextFacility struct {
	XMLName      xml.Name `xml:"contextFacility" json:"-"`
	Kind         string   `xml:"kind,attr"`
	UidRef       string   `xml:"uidRef,attr"`
	Name         string   `xml:",chardata"`
	NamingSystem string   `xml:"namingSystem,attr"`
}

type InstallationReport struct {
	Installation       Installation
	OperationalHSE     []OperationalHSE `xml:"operationalHSE"`
	ProductionActivity ProductionActivity
	BedsAvailable      int         `xml:"bedsAvailable"`
	CrewCounts         []CrewCount `xml:"crewCount"`
	Work               DoubleValue `xml:"work"`
}

type CrewCount struct {
	Type  string `xml:"type,attr"`
	Count int    `xml:",chardata"`
}

type ProductionActivity struct {
	XMLName             xml.Name `xml:"productionActivity" json:"-"`
	LostProduction      LostProduction
	WaterCleaning       []WaterCleaning      `xml:"waterCleaningQuality"`
	ShutDown            []ShutDown           `xml:"shutdown"`
	OperationalComments []OperationalComment `xml:"operationalComment"`
}

type WaterCleaning struct {
	Uid                string `xml:"uid,attr"`
	SamplePoint        string `xml:"samplePoint"`
	OilInWaterProduced Value  `xml:"oilInWaterProduced"`
}

type OperationalHSE struct {
	IncidentCount     int       `xml:"IncidentCount"`
	SinceLostTime     Value     `xml:"sinceLostTime"`
	AlarmCount        int       `xml:"alarmCount"`
	SafetyIntroCount  int       `xml:"safetyIntroCount"`
	SafetyDescription string    `xml:"safetyDescription"`
	Safety            []Safety  `xml:"safety"`
	Weather           []Weather `xml:"weather"`
}

type Safety struct {
	MeantimeIncident Value         `xml:"meantimeIncident"`
	SafetyCount      []SafetyCount `xml:"safetyCount"`
	Comment          []Comment     `xml:"comment"`
}

type SafetyCount struct {
	Type   string `xml:"type,attr"`
	Period string `xml:"period,attr"`
	Count  int    `xml:",chardata"`
}

type LostProduction struct {
	XMLName              xml.Name               `xml:"lostProduction" json:"-"`
	Reasons              []VolumeAndReason      `xml:"volumeAndReason"`
	ThirdPartyProcessing []ThirdPartyProcessing `xml:"thirdPartyProcessing"`
}

type VolumeAndReason struct {
	UoM        string  `xml:"uom,attr"`
	ReasonLost string  `xml:"reasonLost,attr"`
	Value      float64 `xml:",chardata"`
}

type ThirdPartyProcessing struct {
	Installation   Installation `xml:"installation"`
	OilStdTempPres Value        `xml:"oilStdTempPres"`
	GasStdTempPres Value        `xml:"gasStdTempPres"`
}

type ShutDown struct {
	Installation       Installation `xml:"installation"`
	Description        string       `xml:"description"`
	DTimStart          xsdDateTime  `xml:"dTimStart"`
	DTimEnd            xsdDateTime  `xml:"dTimEnd"`
	VolumetricDownTime Value        `xml:"volumetricDownTime"`
	LossOilStdTempPres Value        `xml:"lossOilStdTempPres"`
	LossGasStdTempPres Value        `xml:"lossGasStdTempPres"`
	Activity           []Comment    `xml:"activity"`
}

type OperationalComment struct {
	Type      string      `xml:"type"`
	DTimStart xsdDateTime `xml:"dTimStart"`
	DTimEnd   xsdDateTime `xml:"dTimEnd"`
	Comment   string      `xml:"comment"`
}

type Comment struct {
	Who       string      `xml:"who"`
	Role      string      `xml:"role"`
	DTimStart xsdDateTime `xml:"dTimStart"`
	DTimEnd   xsdDateTime `xml:"dTimeEnd"`
	Comment   string      `xml:"comment"`
}

type Weather struct {
	DTim                xsdDateTime `xml:"dTim"` //dtim
	Agency              string      `xml:"agency"`
	BarometricPressure  Value       `xml:"barometricPressure"`
	BeaufortScaleNumber int         `xml:"beaufortScaleNumber"`
	TempSurfaceMn       Value       `xml:"tempSurfaceMn"`
	TempSurfaceMx       Value       `xml:"tempSurfaceMx"`
	TempWindChill       Value       `xml:"tempWindChill"`
	TempSea             Value       `xml:"tempSea"`
	Visibility          Value       `xml:"visibility"`
	AziWave             Value       `xml:"aziWave"`
	HtWave              Value       `xml:"htWave"`
	SignificantWave     Value       `xml:"significantWave"`
	MaxWave             Value       `xml:"maxWave"`
	PeriodWave          Value       `xml:"periodWave"`
	AziWind             Value       `xml:"aziWind"`
	VelWind             Value       `xml:"velWind"`
	TypePrecip          string      `xml:"typePrecip"`
	AmtPrecip           Value       `xml:"amtPrecip"`
	CoverCloud          string      `xml:"coverCloud"`
	CeilingCloud        Value       `xml:"ceilingCloud"`
	CurrentSea          Value       `xml:"currentSea"`
	AziCurrentSea       Value       `xml:"aziCurrentSea"`
	Comments            string      `xml:"comments"`
}

//add welltests
type WellSet struct {
	XMLName xml.Name `xml:"wellSet"`
	Wells   []Well   `xml:"well"`
}

type Well struct {
	Uid       string     `xml:uid,attr`
	Name      string     `xml:"name"`
	Wellbores []Wellbore `xml:"wellboreSet>wellbore"`
}

type Wellbore struct {
	Name      string     `xml:"name"`
	WellTests []WellTest `xml:"wellTestSet>wellTest"`
}

type WellTest struct {
	Name             string         `xml:"name"`
	TestDate         xsdDateTime    `xml:"testDate"`
	TestType         string         `xml:"testType"`
	ChokeSize        Value          `xml:"chokeOrificeSize"`
	StandardTempPres string         `xml:"standardTempPres"`
	ProductionTest   ProductionTest `xml:"productionTest"`
}

type ProductionTest struct {
	TestDuration Value `xml:"testDuration"`
	WHT          Value `xml:"wellheadData>temperature"`
	WHP          Value `xml:"wellheadData>flowingPressure"`
	ChokeSize    Value `xml:"wellheadData>chokeOrificeSize"`
	SepPress     Value `xml:"separatorData>separatorPressure"`
	SepTemp      Value `xml:"separatorData>separatorTemperature"`
	OilRate      Value `xml:"productionTestResults>oilRate>fluidRateStdTempPres"`
	GasRate      Value `xml:"productionTestResults>gasRate>fluidRateStdTempPres"`
	WatRate      Value `xml:"productionTestResults>waterRate>fluidRate"`
	GoR          Value `xml:"productionTestResults>gasOilRatio"`
}

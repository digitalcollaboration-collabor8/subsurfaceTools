package ddrml

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

const xsdDateTimeLayoutNoTimeZone = "2006-01-02T15:04:05"

type drillTimestamp struct {
	time.Time
}

//const xsdDateTimeLayout = "2006-01-02T15:04:05"
const xsdDateTimeLayout = time.RFC3339
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
		//try to parse it without timezone..
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

type Value struct {
	Value float64 `xml:",chardata"`
	Uom   string  `xml:"uom,attr"`
}

type ProcessingData struct {
	UUid         string //uniqueue identifier for this dataset from a given report
	FileName     string
	FilePath     string
	DocumentName string
	ReadTime     time.Time
}

//unmarshal function to handle xsd date time parsing
func (c *drillTimestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, _ := time.Parse(time.RFC3339, v)
	*c = drillTimestamp{parse}
	return nil
}

type DrillReports struct {
	XMLName            xml.Name      `xml:"drillReports" json:"-"`
	DrillReports       []DrillReport `xml:"drillReport"`
	DataIdentification ProcessingData
}

type DrillReport struct {
	Uid                  string         `xml:"uid,attr"`
	UidWell              string         `xml:"uidWell,attr"`
	UidWellbore          string         `xml:"uidWellbore,attr"`
	NameWell             string         `xml:"nameWell"`
	NameWellbore         string         `xml:"nameWellbore"`
	Name                 string         `xml:"name"`
	DTimStart            drillTimestamp `xml:"dTimStart"`
	DTimEnd              drillTimestamp `xml:"dTimEnd"`
	VersionKind          string         `xml:"versionKind"`
	CreatedDate          drillTimestamp `xml:"createDate"`
	WellAlias            WellAlias
	WellboreAliases      []WellboreAlias `xml:"wellboreAlias"`
	WellboreInfo         WellboreInfo
	StatusInfo           StatusInfo
	BitRecords           []BitRecord           `xml:"bitRecord"`
	CasingLinerTubings   []CasingLinerTubing   `xml:"casing_liner_tubing"`
	CementStages         []CementStage         `xml:"cementStage"`
	Fluids               []Fluid               `xml:"fluid"`
	PorePressures        []PorePressure        `xml:"porePressure"`
	ExtendedReport       ValueDTim             `xml:"extendedReport"`
	SurveyStations       []SurveyStation       `xml:"surveyStation"`
	Activities           []Activity            `xml:"activity"`
	LogInfos             []LogInfo             `xml:"logInfo"`
	CoreInfos            []CoreInfo            `xml:"coreInfo"`
	WellTestInfos        []WellTestInfo        `xml:"wellTestInfo"`
	FormTestInfos        []FormTestInfo        `xml:"formTestInfo"`
	LithShowInfos        []LithShowInfo        `xml:"lithShowInfo"`
	EquipFailureInfos    []EquipFailureInfo    `xml:"equipFailureInfo"`
	ControlIncidentInfos []ControlIncidentInfo `xml:"controlIncidentInfo"`
	StratInfos           []StratInfo           `xml:"stratInfo"`
	PerfInfos            []PerfInfo            `xml:"perfInfo"`
	GasReadingInfos      []GasReadingInfo      `xml:"gasReadingInfo"`
	Weathers             []Weather             `xml:"weather"`
	Comment              string                `xml:"comment"`
	DataIdentification   ProcessingData
}

type WellAlias struct {
	XMLName      xml.Name `xml:"wellAlias" json:"-"`
	Name         string   `xml:"name"`
	NamingSystem string   `xml:"namingSystem"`
}

type WellboreAlias struct {
	Name         string `xml:"name"`
	NamingSystem string `xml:"namingSystem"`
}

type WellboreInfo struct {
	XMLName           xml.Name       `xml:"wellboreInfo" json:"-"`
	DTimSpud          drillTimestamp `xml:"dTimSpud"`
	DTimPreSpud       drillTimestamp `xml:"dTimPreSpud"`
	DateDrillComplete xsdDate        `xml:"dateDrillComplete"`
	DaysAhead         xsdDouble      `xml:"daysAhead"`
	DaysBehind        xsdDouble      `xml:"daysBehind"`
	Operator          string         `xml:"operator"`
	DrillContractor   string         `xml:"drillContractor"`
	RigAliases        []RigAlias     `xml:"rigAlias"`
}

type RigAlias struct {
	XMLName      xml.Name `xml:"rigAlias" json:"-"`
	Name         string   `xml:"name"`
	NamingSystem string   `xml:"namingSystem"`
}

type StatusInfo struct {
	XMLName           xml.Name       `xml:"statusInfo" json:"-"`
	Reportnumber      int            `xml:"reportNo"`
	DTim              drillTimestamp `xml:"dTim"`
	Md                Value          `xml:"md"`
	Tvd               Value          `xml:"tvd"`
	MdPlugTop         Value          `xml:"mdPlugTop"`
	DiaHole           Value          `xml:"diaHole"`
	DTimDiaHoleStart  drillTimestamp `xml:"dTimDiaHoleStart"` //DTim
	MdDiaHoleStart    Value          `xml:"mdDiaHoleStart"`
	DiaPilot          Value          `xml:"diaPilot"`
	MdDiaPilotPlan    Value          `xml:"mdDiaPilotPlan"`
	TVDDiaPilotPlan   Value          `xml:"tvdDiaPilotPlan"`
	TypeWellBore      string         `xml:"typeWellbore"`
	PrimaryConveyance string         `xml:"primaryConveyance"`
	MdKickoff         Value          `xml:"mdKickoff"`
	TvdKickoff        Value          `xml:"tvdKickoff"`
	StrengthForm      Value          `xml:"strengthForm"`
	MdStrengthForm    Value          `xml:"mdStrengthForm"`
	TvdStrengthForm   Value          `xml:"tvdStrengthForm"`
	DiaCasingLast     Value          `xml:"diaCsgLast"`
	MdCasingLast      Value          `xml:"mdCsgLast"`
	TvdCasingLast     Value          `xml:"tvdCsgLast"`
	PressTestType     string         `xml:"presTestType"`
	MdPlanned         Value          `xml:"mdPlanned"`
	DistDrilled       Value          `xml:"distDrill"`
	ElevKelly         Value          `xml:"elevKelly"`
	WellheadElevation Value          `xml:"wellheadElevation"`
	WaterDepth        Value          `xml:"waterDepth"`
	Sum24Hr           string         `xml:"sum24Hr"`
	Forecast24Hr      string         `xml:"forecast24Hr"`
	RopCurrent        Value          `xml:"ropCurrent"`
	TightWell         bool           `xml:"tightWell"`
	HPHT              bool           `xml:"hpht"`
	AvgPresBH         Value          `xml:"avgPresBH"`
	AvgTempBH         Value          `xml:"avgTempBH"`
	FixedRig          bool           `xml:"fixedRig"`
}

type BitRecord struct {
	NumBitRun         string `xml:"numBitRun"`
	NumBit            string `xml:"numBit"`
	DiaBit            Value  `xml:"diaBit"`
	Manufacturer      string `xml:"manufacturer"`
	CodeMfg           string `xml:"codeMfg"`
	DullGrade         string `xml:"dullGrade"`
	CodeIADC          string `xml:"codeIADC"`
	CondFinalInner    int    `xml:"condFinalInner"`
	CondFinalOuter    int    `xml:"condFinalOuter"`
	CondFinalDull     string `xml:"condFinalDull"`
	CondFinalLocation string `xml:"condFinalLocation"`
	CondFinalBearing  string `xml:"condFinalBearing"`
	CondFinalGauge    string `xml:"condFinalGauge"`
	CondFinalOther    string `xml:"condFinalOther"`
	CondFinalReason   string `xml:"condFinalReason"`
	BitRun            BitRun
	Nozzles           []Nozzle `xml:"nozzle"`
}

type BitRun struct {
	XMLName         xml.Name `xml:"bitRun" json:"-"`
	ETimOpBit       Value    `xml:"eTimOpBit"`
	MDHoleStart     Value    `xml:"mdHoleStart"`
	MDHoleStop      Value    `xml:"mdHoleStop"`
	RopAv           Value    `xml:"ropAv"`
	MDHoleMadeRun   Value    `xml:"mdHoleMadeRun"`
	HrsDrilled      Value    `xml:"hrsDrilled"`
	HrsDrilledRun   Value    `xml:"hrsDrilledRun"`
	MdTotalHoleMade Value    `xml:"mdTotHoleMade"`
	TotHrsDrilled   Value    `xml:"totHrsDrilled"`
	TotRop          Value    `xml:"totRop"`
}

type Nozzle struct {
	NumNozzle int   `xml:"numNozzle"`
	DiaNozzle Value `xml:"diaNozzle"`
}

type CasingLinerTubing struct {
	Type                 string `xml:"type"`
	Id                   Value  `xml:"id"`
	Od                   Value  `xml:"od"`
	Weight               Value  `xml:"weight"`
	Grade                string `xml:"grade"`
	Connection           string `xml:"connection"`
	Length               Value  `xml:"length"`
	MdTop                Value  `xml:"mdTop"`
	MdBottom             Value  `xml:"mdBottom"`
	CasingLinerTubingRun CasingLinerTubingRun
	Comment              string `xml:"comment"`
}

type CasingLinerTubingRun struct {
	XMLName     xml.Name       `xml:"casing_liner_tubing_run" json:"-"`
	CasingType  string         `xml:"casingType"`
	Description string         `xml:"description"`
	DTimStart   drillTimestamp `xml:"dTimStart"`
	DTimEnd     drillTimestamp `xml:"dTimEnd"`
}

type CementStage struct {
	DTimPumpStart    drillTimestamp   `xml:"dTimPumpStart"`
	DTimPumpEnd      drillTimestamp   `xml:"dTimPumpEnd"`
	JobType          string           `xml:"jobType"`
	CasingStrDia     Value            `xml:"casingStrDia"`
	Comments         string           `xml:"comments"`
	VolReturns       Value            `xml:"volReturns"`
	CementingFluids  []CementingFluid `xml:"cementingFluid"`
	DTimPresReleased drillTimestamp   `xml:"dTimPresReleased"`
	AnnFlowAfter     bool             `xml:"annFlowAfter"`
	TopPlug          bool             `xml:"topPlug"`
	BotPlug          bool             `xml:"botPlug"`
	PlugBumped       bool             `xml:"plugBumped"`
	PresBump         Value            `xml:"presBump"`
	FloatHeld        bool             `xml:"floatHeld"`
	Reciprocated     bool             `xml:"reciprocated"`
	Rotated          bool             `xml:"rotated"`
}

type CementingFluid struct {
	TypeFluid      string `xml:"typeFluid"`
	DescFluid      string `xml:"descFluid"`
	RatioMixWater  Value  `xml:"ratioMixWater"`
	Density        Value  `xml:"density"`
	VolPumped      Value  `xml:"volPumped"`
	Yp             Value  `xml:"yp"`
	ETimThickening Value  `xml:"eTimThickening"`
	PCFreeWater    Value  `xml:"pcFreeWater"`
	Comments       string `xml:"comments"`
}

type Fluid struct {
	Type            string         `xml:"type"`
	LocationSample  string         `xml:"locationSample"`
	DTim            drillTimestamp `xml:"dTim"` //dtim
	Md              Value          `xml:"md"`
	Tvd             Value          `xml:"tvd"`
	PresBopRating   Value          `xml:"presBopRating"`
	MudClass        string         `xml:"mudClass"`
	Density         Value          `xml:"density"`
	VisFunnel       Value          `xml:"visFunnel"`
	PV              Value          `xml:"pv"`
	YP              Value          `xml:"yp"`
	Gel10Sec        Value          `xml:"gel10Sec"`
	Gel10Min        Value          `xml:"gel10Min"`
	Gel30Min        Value          `xml:"gel30Min"`
	FilterCakeLtlp  Value          `xml:"filterCakeLtlp"`
	FiltrateLtlp    Value          `xml:"filtrateLtlp"`
	TempHtHp        Value          `xml:"tempHthp"`
	FiltrateHtHp    Value          `xml:"filtrateHthp"`
	FilterCakeHtHp  Value          `xml:"filterCakeHthp"`
	SolidsPc        Value          `xml:"solidsPc"`
	WaterPc         Value          `xml:"waterPc"`
	OilPc           Value          `xml:"oilPc"`
	SandPc          Value          `xml:"sandPc"`
	SolidsLowGravPc Value          `xml:"solidsLowGravPc"`
	PH              float64        `xml:"ph"`
	PM              Value          `xml:"pm"`
	PMFiltrate      Value          `xml:"pmFiltrate"`
	MF              Value          `xml:"mf"`
	Chloride        Value          `xml:"chloride"`
	Calcium         Value          `xml:"calcium"`
	Magnesium       Value          `xml:"magnesium"`
	Rheometers      []Rheometer    `xml:"rheometer"`
	Lime            Value          `xml:"lime"`
	SolidsHiGravPc  Value          `xml:"solidsHiGravPc"`
	SolCorPc        Value          `xml:"solCorPc"`
	Comments        string         `xml:"comments"`
}

type Rheometer struct {
	TempRheom  Value   `xml:"tempRheom"`
	PressRheom Value   `xml:"presRheom"`
	Vis3Rpm    float64 `xml:"vis3Rpm"`
	Vis6Rpm    float64 `xml:"vis6Rpm"`
	Vis30Rpm   float64 `xml:"vis30Rpm"`
	Vis60Rpm   float64 `xml:"vis60Rpm"`
	Vis100Rpm  float64 `xml:"vis100Rpm"`
	Vis200Rpm  float64 `xml:"vis200Rpm"`
	Vis300Rpm  float64 `xml:"vis300Rpm"`
	Vis600Rpm  float64 `xml:"vis600Rpm"`
}

type PorePressure struct {
	ReadingKind         string         `xml:"readingKind"`
	EquivalentMudWeight Value          `xml:"equivalentMudWeight"`
	DTim                drillTimestamp `xml:"dtim"` //dtim
	Md                  Value          `xml:"md"`
	Tvd                 Value          `xml:"tvd"`
	Comment             string         `xml:"comment"`
}

//use to tackle an element with a dtim attribute e.g. <extendedReport dTim="2013-09-08T00:00:00Z">
type ValueDTim struct {
	Value string         `xml:",chardata"`
	DTim  drillTimestamp `xml:"dTim,attr"` //dtim
}

type SurveyStation struct {
	DTim    drillTimestamp `xml:"dTim"`
	Md      Value          `xml:"md"`
	Tvd     Value          `xml:"tvd"`
	Incl    Value          `xml:"incl"`
	Azi     Value          `xml:"azi"`
	Comment string         `xml:"comment"`
}

type Activity struct {
	DTimStart           drillTimestamp `xml:"dTimStart"` //dtim
	DTimeEnd            drillTimestamp `xml:"dTimEnd"`   //dtim
	Md                  Value          `xml:"md"`
	Tvd                 Value          `xml:"tvd"`
	Phase               string         `xml:"phase"`
	ProprietaryCode     string         `xml:"proprietaryCode"`
	Conveyance          string         `xml:"conveyance"`
	MdHoleStart         Value          `xml:"mdHoleStart"`
	State               string         `xml:"state"`
	StateDetailActivity string         `xml:"stateDetailActivity"`
	Comment             string         `xml:"comments"`
}

type LogInfo struct {
	DTim           drillTimestamp `xml:"dTim"` //dtim
	RunNumber      string         `xml:"runNumber"`
	ServiceCompany string         `xml:"serviceCompany"`
	Service        string         `xml:"service"`
	MdTop          Value          `xml:"mdTop"`
	MdBottom       Value          `xml:"mdBottom"`
	TvdTop         Value          `xml:"tvdTop"`
	TvdBottom      Value          `xml:"tvdBottom"`
	Tool           string         `xml:"tool"`
	TempBHCt       Value          `xml:"tempBHCT"`
	TempBHST       Value          `xml:"tempBHST"`
	ETimStatic     Value          `xml:"eTimStatic"`
	MdTempTool     Value          `xml:"mdTempTool"`
	TvdTempTool    Value          `xml:"tvdTempTool"`
	Comment        string         `xml:"comment"`
}

type CoreInfo struct {
	DTim            drillTimestamp `xml:"dTim"` //dtim
	CoreNumber      string         `xml:"coreNumber"`
	MDTop           Value          `xml:"mdTop"`
	MDBottom        Value          `xml:"mdBottom"`
	TvdTop          Value          `xml:"tvdTop"`
	TvdBottom       Value          `xml:"tvdBottom"`
	LenRecovered    Value          `xml:"lenRecovered"`
	RecoverPC       Value          `xml:"recoverPc"`
	LenBarrel       Value          `xml:"lenBarrel"`
	InnerBarrelType string         `xml:"innerBarrelType"`
	CoreDescription string         `xml:"coreDescription"`
}

type WellTestInfo struct {
	DTim            drillTimestamp `xml:"dTim"` //dtim
	TestType        string         `xml:"testType"`
	TestNumber      int            `xml:"testNumber"`
	MdTop           Value          `xml:"mdTop"`
	MdBottom        Value          `xml:"mdBottom"`
	TvdTop          Value          `xml:"tvdTop"`
	TvdBottom       Value          `xml:"tvdBottom"`
	ChokeSize       Value          `xml:"chokeOrificeSize"`
	DensityOil      Value          `xml:"densityOil"`
	DensityWater    Value          `xml:"densityWater"`
	DensityGas      Value          `xml:"densityGas"`
	FlowRateOil     Value          `xml:"flowRateOil"`
	FlowRateWater   Value          `xml:"flowRateWater"`
	FlowRateGas     Value          `xml:"flowRateGas"`
	PresShutIn      Value          `xml:"presShutIn"`
	PresFlowing     Value          `xml:"presFlowing"`
	PresBottom      Value          `xml:"presBottom"`
	GoR             Value          `xml:"gasOilRatio"`
	WaterOilRatio   Value          `xml:"waterOilRatio"`
	Chloride        Value          `xml:"chloride"`
	CarbonDioxide   Value          `xml:"carbonDioxide"`
	HydrogenSulfide Value          `xml:"hydrogenSulfide"`
	VolOilTotal     Value          `xml:"volOilTotal"`
	VolGasTotal     Value          `xml:"volGasTotal"`
	VolWaterTotal   Value          `xml:"volWaterTotal"`
	VolOilStored    Value          `xml:"volOilStored"`
	Comment         string         `xml:"comment"`
}

type FormTestInfo struct {
	DTim                  drillTimestamp `xml:"dTim"` //dtim
	RunNumber             string         `xml:"runNumber"`
	TestNumber            int            `xml:"testNUmber"`
	Md                    Value          `xml:"md"`
	Tvd                   Value          `xml:"tvd"`
	PresPore              Value          `xml:"presPore"`
	FluidDensity          Value          `xml:"fluidDensity"`
	HydrostaticPresBefore Value          `xml:"hydrostaticPresBefore"`
	LeakOffPressure       Value          `xml:"leakOffPressure"`
	GoodSeal              bool           `xml:"goodSeal"`
	MdSample              Value          `xml:"mdSample"`
	DominateComponent     string         `xml:"dominateComponent"`
	DensityHC             Value          `xml:"densityHC"`
	VolumeSample          Value          `xml:"volumeSample"`
	Description           string         `xml:"description"`
}

type LithShowInfo struct {
	DTim      drillTimestamp `xml:"dTim"` //dTim
	MdTop     Value          `xml:"mdTop"`
	MdBottom  Value          `xml:"mdBottom"`
	TvdTop    Value          `xml:"tvdTop"`
	TvdBottom Value          `xml:"tvdBottom"`
	Show      string         `xml:"show"`
	Lithology string         `xml:"lithology"`
}

type EquipFailureInfo struct {
	DTim               drillTimestamp `xml:"dTim"` //dtim
	Md                 Value          `xml:"md"`
	Tvd                Value          `xml:"tvd"`
	EquipClass         string         `xml:"equipClass"`
	ETimMissProduction Value          `xml:"eTimMissProduction"`
	DTimRepair         drillTimestamp `xml:"dTimRepair"`
	Description        string         `xml:"description"`
}

type ControlIncidentInfo struct {
	DTim             drillTimestamp `xml:"dTim"` //dTim
	MdInflow         Value          `xml:"mdInflow"`
	TvdInflow        Value          `xml:"tvdInflow"`
	Phase            string         `xml:"phase"`
	ProprietaryCode  string         `xml:"proprietaryCode"`
	ETimLost         Value          `xml:"eTimLost"`
	DTimRegained     drillTimestamp `xml:"dTimRegained"` //dTim
	DiaBit           Value          `xml:"diaBit"`
	MdBit            Value          `xml:"mdBit"`
	WtMud            Value          `xml:"wtMud"`
	PorePressure     Value          `xml:"porePressure"`
	DiaCsgLast       Value          `xml:"diaCsgLast"`
	MdCsgLast        Value          `xml:"mdCsgLast"`
	VolMudGained     Value          `xml:"volMudGained"`
	PresShutinCasing Value          `xml:"presShutinCasing"`
	PresShutInDrill  Value          `xml:"presShutInDrill"`
	IncidentType     string         `xml:"incidentType"`
	KillingType      string         `xml:"killingType"`
	Formation        string         `xml:"formation"`
	TempBottom       Value          `xml:"tempBottom"`
	PresMaxChoke     Value          `xml:"presMaxChoke"`
	Description      string         `xml:"description"`
}

type StratInfo struct {
	DTim          drillTimestamp `xml:"dTim"`
	MdTopPlanned  Value          `xml:"mdTopPlanned"`
	TvdTopPlanned Value          `xml:"tvdTopPlanned"`
	MdTop         Value          `xml:"mdTop"`
	TvdTop        Value          `xml:"tvdTop"`
	Description   string         `xml:"description"`
}

type PerfInfo struct {
	DTimOpen  drillTimestamp `xml:"dTimOpen"` //dtim
	DTimClose drillTimestamp `xml:"dTimClose"`
	MdTop     Value          `xml:"mdTop"`
	MdBottom  Value          `xml:"mdBottom"`
	TvdTop    Value          `xml:"tvdTop"`
	TvdBottom Value          `xml:"tvdBottom"`
	Comment   string         `xml:"comment"`
}

type GasReadingInfo struct {
	DTim        drillTimestamp `xml:"dTim"` //dtim
	ReadingType string         `xml:"readingType"`
	MdTop       Value          `xml:"mdTop"`
	MdBottom    Value          `xml:"mdBottom"`
	TvdTop      Value          `xml:"tvdTop"`
	TvdBottom   Value          `xml:"tvdBottom"`
	GasHigh     Value          `xml:"gasHigh"`
	GasLow      Value          `xml:"gasLow"`
	Meth        Value          `xml:"meth"`
	Eth         Value          `xml:"eth"`
	Prop        Value          `xml:"prop"`
	Ibut        Value          `xml:"ibut"`
	NBut        Value          `xml:"nbut"`
	IPent       Value          `xml:"ipent"`
	NPent       Value          `xml:"npent"`
	Comment     string         `xml:comment`
}

type Weather struct {
	DTim                drillTimestamp `xml:"dTim"` //dtim
	Agency              string         `xml:"agency"`
	BarometricPressure  Value          `xml:"barometricPressure"`
	BeaufortScaleNumber int            `xml:"beaufortScaleNumber"`
	TempSurfaceMn       Value          `xml:"tempSurfaceMn"`
	TempSurfaceMx       Value          `xml:"tempSurfaceMx"`
	TempWindChill       Value          `xml:"tempWindChill"`
	TempSea             Value          `xml:"tempSea"`
	Visibility          Value          `xml:"visibility"`
	AziWave             Value          `xml:"aziWave"`
	HtWave              Value          `xml:"htWave"`
	SignificantWave     Value          `xml:"significantWave"`
	MaxWave             Value          `xml:"maxWave"`
	PeriodWave          Value          `xml:"periodWave"`
	AziWind             Value          `xml:"aziWind"`
	VelWind             Value          `xml:"velWind"`
	TypePrecip          string         `xml:"typePrecip"`
	AmtPrecip           Value          `xml:"amtPrecip"`
	CoverCloud          string         `xml:"coverCloud"`
	CeilingCloud        Value          `xml:"ceilingCloud"`
	CurrentSea          Value          `xml:"currentSea"`
	AziCurrentSea       Value          `xml:"aziCurrentSea"`
	Comments            string         `xml:"comments"`
}

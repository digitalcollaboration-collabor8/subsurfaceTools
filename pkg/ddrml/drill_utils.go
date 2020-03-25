/*
* @Author: magsv
* @Date:   2017-08-02 10:47:54
* @Last Modified by:   magsv
* @Last Modified time: 2019-01-22 12:42:53
 */

package ddrml

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

//Builds an excel output file for a given list of mprml objects
func BuildXLSFileForDrilling(path string, dReports []DrillReports, oneFilePerSheet bool, appendTimeInName bool) error {
	//just need to join everything
	var drillReports []DrillReport
	for i := 0; i < len(dReports); i++ {
		drillReports = append(drillReports, dReports[i].DrillReports...)
	}
	dataSets := BuildDDRDataset(drillReports)

	return common.CreateWorkbookFromDataSet(path, dataSets, oneFilePerSheet, appendTimeInName)
}

func BuildCsvFileForDrilling(path string, dReports []DrillReports) error {
	var drillReports []DrillReport
	for i := 0; i < len(dReports); i++ {
		drillReports = append(drillReports, dReports[i].DrillReports...)
	}
	dataSets := BuildDDRDataset(drillReports)
	return common.DatasetsToCsv(dataSets, path, ";")
}

//builds a json output file for drillind data
func BuildJsonFileForDrilling(path string, dReports []DrillReports) error {
	var err error
	var data []byte
	var drillReports []DrillReport
	for i := 0; i < len(dReports); i++ {
		drillReports = append(drillReports, dReports[i].DrillReports...)
	}
	dataSets := BuildDDRDataset(drillReports)

	if data, err = common.DatasetsToJson(dataSets); err != nil {
		return err
	}
	//need to write it to a file
	return common.Write2File(path, data)
}

//Reads a set of xml ddr files from a given folder path and
//parses them into a list of struct ddr objects to be used for further processing
func ReadDDRXMLFiles2Struct(folderPath string) ([]DrillReports, error) {
	var err error
	var dReports []DrillReports
	var files []string
	var dReport DrillReports
	start := time.Now()

	if files, err = common.GetFilesWithExtension(folderPath, "*.xml"); err != nil {
		zap.S().Error("Failed in reading file in folder:", err.Error())
		return dReports, err
	}
	fileSearchTook := time.Since(start)
	xmlStart := time.Now()
	for i := 0; i < len(files); i++ {
		zap.S().Info("Processing ddr xml file:", files[i])
		if dReport, err = ParseDDRFile(files[i]); err != nil {
			zap.S().Errorf("Failed in parsing ddr xml file:%s\n", err.Error())
			return dReports, err
		}
		dReports = append(dReports, dReport)
	}
	xmlParseTook := time.Since(xmlStart)

	zap.S().Infof("File search took:%s", fileSearchTook)
	zap.S().Infof("XML parse took:%s", xmlParseTook)
	return dReports, nil
}

func ParseDDRFile(inputFile string) (DrillReports, error) {
	var drillReports DrillReports
	var err error
	//var data []byte
	var data io.Reader

	/*if data, err = common.ReadFile(inputFile); err != nil {
		return drillReports, err
	}*/
	if data, err = os.Open(inputFile); err != nil {
		return drillReports, err
	}

	//make sure we can read iso-8859-1
	decoder := xml.NewDecoder(data)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&drillReports); err != nil {
		return drillReports, err
	}
	/*if err = xml.Unmarshal(data, &drillReports); err != nil {
		return drillReports, err
	}*/
	//just add the uid and file data for this set of objects
	drillReports.DataIdentification = createDataId(inputFile)

	for i := 0; i < len(drillReports.DrillReports); i++ {
		drillReports.DrillReports[i].DataIdentification = drillReports.DataIdentification
	}
	return drillReports, nil

}

//Creates uuid, fileinformation for a given file
func createDataId(inputFile string) ProcessingData {
	pData := ProcessingData{}
	pData.FileName = common.GetFileName(inputFile)
	pData.FilePath = inputFile
	pData.UUid = common.CreateUUID()
	pData.ReadTime = time.Now()
	return pData

}

func BuildDDRDataset(dReports []DrillReport) []common.DataSet {
	var dSets []common.DataSet
	dSets = append(dSets, buildDDRReportFileInfo(dReports, "REPORT_FILE_INFO"))
	dSets = append(dSets, buildDDRDrillReportInfo(dReports, "REPORT_INFO"))
	dSets = append(dSets, buildDDRWellboreInfo(dReports, "WELLBORE_INFO"))
	dSets = append(dSets, buildDDRStatusInfo(dReports, "STATUS_INFO"))
	dSets = append(dSets, buildDDRBitRecords(dReports, "BIT_RECORDS"))
	dSets = append(dSets, buildDDRCasingLinerTubings(dReports, "CASING_LINER_TUBING"))
	dSets = append(dSets, buildDDRCementStages(dReports, "CEMENT_STAGES"))
	dSets = append(dSets, buildDDRFluid(dReports, "FLUIDS"))
	dSets = append(dSets, buildDDRPorePressure(dReports, "PORE_PRESSURES"))
	dSets = append(dSets, buildDDRSurveyStation(dReports, "SURVEY_STATIONS"))
	dSets = append(dSets, buildDDRActivities(dReports, "ACTIVITIES"))
	dSets = append(dSets, buildDDRLoginfo(dReports, "LOG_INFO"))
	dSets = append(dSets, buildDDRCoreInfo(dReports, "CORE_INFO"))
	dSets = append(dSets, buildDDRWellTestInfo(dReports, "WELLTEST_INFO"))
	dSets = append(dSets, buildDDRFormTestInfo(dReports, "FORMTEST_INFO"))
	dSets = append(dSets, buildDDRLithShowInfo(dReports, "LITHSHOW_INFO"))
	dSets = append(dSets, buildEquipFailureInfo(dReports, "EQUIPFAILURE_INFO"))
	dSets = append(dSets, buildControlIncidentInfo(dReports, "CONTROLINCIDENT_INFO"))
	dSets = append(dSets, buildStratInfo(dReports, "STRAT_INFO"))
	dSets = append(dSets, buildPerfInfo(dReports, "PERF_INFO"))
	dSets = append(dSets, buildGasReadingInfo(dReports, "GASREADING_INFO"))
	dSets = append(dSets, buildWeather(dReports, "WEATHER"))
	return dSets
}

func buildDDRReportFileInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	headers = []string{"DataUUID", "FileName", "FilePath"}

	for i := 0; i < len(dReports); i++ {
		row := common.RowData{}
		row.AddStrValue(dReports[i].DataIdentification.UUid)
		row.AddStrValue(dReports[i].DataIdentification.FileName)
		row.AddStrValue(dReports[i].DataIdentification.FilePath)
		rows = append(rows, row)
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRFluid(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "Type", "LocationSample",
		"DTim", "MD", "MD_UoM", "TVD", "TVD_UoM", "PresBopRating", "PresBopRating_UoM",
		"MudClass", "Density", "Density_UoM", "VisFunnel", "VisFunnel_UoM", "PV", "PV_UoM",
		"YP", "YP_UoM", "Gel10Sec", "Gel10Sec_UoM", "Gel10Min", "Gel10Min_UoM",
		"Gel30Min", "Gel30Min_UoM", "FilterCakeLtlp", "FilterCakeLtlp_UoM",
		"FiltrateLtlp", "FiltrateLtlp_UoM", "TempHthp", "TempHthp_UoM",
		"FiltrateHthp", "FiltrateHthp_UoM", "FilterCakeHthp", "FilterCakegthp_UoM",
		"SolidsPc", "SolidsPc_UoM", "WaterPc", "WaterPc_UoM", "OilPc", "OilPc_UoM",
		"SandPc", "SandPc_UoM", "SolidsLowGravPc", "SolidsLowGravPc_UoM",
		"PH", "PM", "PM_UoM", "PMFiltrate", "PMFiltrate_UoM", "MF", "MF_UoM", "Chloride", "Chloride_UoM",
		"Calcium", "Calcium_UoM", "Magnesium", "Magnesium_UoM",
		"TempRheom", "TempRheom_UoM", "PresRheom", "PresRheom_UoM",
		"Vis3Rpm", "Vis6Rpm", "Vis30Rpm", "Vis60Rpm", "Vis100Rpm", "Vis200Rpm",
		"Vis300Rpm", "Vis600Rpm", "Lime", "Lime_UoM", "SolidsHiGravPc", "SolidsHiGravPc_UoM",
		"SolCorPc", "SolCorPc_UoM", "Comments"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].Fluids); s++ {
			fluid := dReports[i].Fluids[s]
			//need to add one row per rheometer
			for x := 0; x < len(fluid.Rheometers); x++ {
				rheom := fluid.Rheometers[x]
				row := idSet.Rows[0]
				row.AddStrValue(recordId)
				row.AddStrValue(fluid.Type)
				row.AddStrValue(fluid.LocationSample)
				row.AddTimeValue(fluid.DTim.Time)
				row.AddFloatValue(fluid.Md.Value)
				row.AddStrValue(fluid.Md.Uom)
				row.AddFloatValue(fluid.Tvd.Value)
				row.AddStrValue(fluid.Tvd.Uom)
				row.AddFloatValue(fluid.PresBopRating.Value)
				row.AddStrValue(fluid.PresBopRating.Uom)
				row.AddStrValue(fluid.MudClass)
				row.AddFloatValue(fluid.Density.Value)
				row.AddStrValue(fluid.Density.Uom)
				row.AddFloatValue(fluid.VisFunnel.Value)
				row.AddStrValue(fluid.VisFunnel.Uom)
				row.AddFloatValue(fluid.PV.Value)
				row.AddStrValue(fluid.PV.Uom)
				row.AddFloatValue(fluid.YP.Value)
				row.AddStrValue(fluid.YP.Uom)
				row.AddFloatValue(fluid.Gel10Sec.Value)
				row.AddStrValue(fluid.Gel10Sec.Uom)
				row.AddFloatValue(fluid.Gel10Min.Value)
				row.AddStrValue(fluid.Gel10Min.Uom)
				row.AddFloatValue(fluid.Gel30Min.Value)
				row.AddStrValue(fluid.Gel30Min.Uom)
				row.AddFloatValue(fluid.FilterCakeLtlp.Value)
				row.AddStrValue(fluid.FilterCakeLtlp.Uom)
				row.AddFloatValue(fluid.FiltrateLtlp.Value)
				row.AddStrValue(fluid.FiltrateLtlp.Uom)
				row.AddFloatValue(fluid.TempHtHp.Value)
				row.AddStrValue(fluid.TempHtHp.Uom)
				row.AddFloatValue(fluid.FiltrateHtHp.Value)
				row.AddStrValue(fluid.FiltrateHtHp.Uom)
				row.AddFloatValue(fluid.FilterCakeHtHp.Value)
				row.AddStrValue(fluid.FilterCakeHtHp.Uom)
				row.AddFloatValue(fluid.SolidsPc.Value)
				row.AddStrValue(fluid.SolidsPc.Uom)
				row.AddFloatValue(fluid.WaterPc.Value)
				row.AddStrValue(fluid.WaterPc.Uom)
				row.AddFloatValue(fluid.OilPc.Value)
				row.AddStrValue(fluid.OilPc.Uom)
				row.AddFloatValue(fluid.SandPc.Value)
				row.AddStrValue(fluid.SandPc.Uom)
				row.AddFloatValue(fluid.SolidsLowGravPc.Value)
				row.AddStrValue(fluid.SolidsLowGravPc.Uom)
				row.AddFloatValue(fluid.PH)
				row.AddFloatValue(fluid.PM.Value)
				row.AddStrValue(fluid.PM.Uom)
				row.AddFloatValue(fluid.PMFiltrate.Value)
				row.AddStrValue(fluid.PMFiltrate.Uom)
				row.AddFloatValue(fluid.MF.Value)
				row.AddStrValue(fluid.MF.Uom)
				row.AddFloatValue(fluid.Chloride.Value)
				row.AddStrValue(fluid.Chloride.Uom)
				row.AddFloatValue(fluid.Calcium.Value)
				row.AddStrValue(fluid.Calcium.Uom)
				row.AddFloatValue(fluid.Magnesium.Value)
				row.AddStrValue(fluid.Magnesium.Uom)
				row.AddFloatValue(rheom.TempRheom.Value)
				row.AddStrValue(rheom.TempRheom.Uom)
				row.AddFloatValue(rheom.PressRheom.Value)
				row.AddStrValue(rheom.PressRheom.Uom)
				row.AddFloatValue(rheom.Vis3Rpm)
				row.AddFloatValue(rheom.Vis6Rpm)
				row.AddFloatValue(rheom.Vis30Rpm)
				row.AddFloatValue(rheom.Vis60Rpm)
				row.AddFloatValue(rheom.Vis100Rpm)
				row.AddFloatValue(rheom.Vis200Rpm)
				row.AddFloatValue(rheom.Vis300Rpm)
				row.AddFloatValue(rheom.Vis600Rpm)
				row.AddFloatValue(fluid.Lime.Value)
				row.AddStrValue(fluid.Lime.Uom)
				row.AddFloatValue(fluid.SolidsHiGravPc.Value)
				row.AddStrValue(fluid.SolidsHiGravPc.Uom)
				row.AddFloatValue(fluid.SolCorPc.Value)
				row.AddStrValue(fluid.SolCorPc.Uom)
				row.AddStrValue(fluid.Comments)

				rows = append(rows, row)
			}
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildWeather(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "Agency", "BarometricPressure",
		"BarometricPressure_UoM", "BeaufortScaleNumber", "TempSurfaceMn", "TempSurfaceMn_UoM",
		"TempSurfaceMx", "TempSurfaceMx_UoM", "TempWindChill", "TempWindChill_UoM",
		"TempSea", "TempSea_UoM", "Visibility", "Visibility_UoM", "AziWave",
		"AziWave_UoM", "HtWave", "HtWave_UoM",
		"SignificantWave", "SignificantWave_UoM", "MaxWave", "MaxWave_UoM", "PeriodWave",
		"PeriodWave_UoM", "AziWind", "AziWind_UoM", "VelWind", "VelWind_UoM", "TypePrecip",
		"AmtPrecip", "AmtPrecip_UoM", "CoverCloud", "CeilingCloud", "CeilingCloud_UoM",
		"CurrentSea", "CurrentSea_UoM", "AziCurrentSea", "AziCurrentSea_UoM", "Comments"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].Weathers); s++ {
			wInfo := dReports[i].Weathers[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(wInfo.DTim.Time)
			row.AddStrValue(wInfo.Agency)
			row.AddFloatValue(wInfo.BarometricPressure.Value)
			row.AddStrValue(wInfo.BarometricPressure.Uom)
			row.AddIntValue(wInfo.BeaufortScaleNumber)
			row.AddFloatValue(wInfo.TempSurfaceMn.Value)
			row.AddStrValue(wInfo.TempSurfaceMn.Uom)
			row.AddFloatValue(wInfo.TempSurfaceMx.Value)
			row.AddStrValue(wInfo.TempSurfaceMx.Uom)
			row.AddFloatValue(wInfo.TempWindChill.Value)
			row.AddStrValue(wInfo.TempWindChill.Uom)
			row.AddFloatValue(wInfo.TempSea.Value)
			row.AddStrValue(wInfo.TempSea.Uom)
			row.AddFloatValue(wInfo.Visibility.Value)
			row.AddStrValue(wInfo.Visibility.Uom)
			row.AddFloatValue(wInfo.AziWave.Value)
			row.AddStrValue(wInfo.AziWave.Uom)
			row.AddFloatValue(wInfo.HtWave.Value)
			row.AddStrValue(wInfo.HtWave.Uom)
			row.AddFloatValue(wInfo.SignificantWave.Value)
			row.AddStrValue(wInfo.SignificantWave.Uom)
			row.AddFloatValue(wInfo.SignificantWave.Value)
			row.AddStrValue(wInfo.SignificantWave.Uom)
			row.AddFloatValue(wInfo.MaxWave.Value)
			row.AddStrValue(wInfo.MaxWave.Uom)
			row.AddFloatValue(wInfo.PeriodWave.Value)
			row.AddStrValue(wInfo.PeriodWave.Uom)
			row.AddFloatValue(wInfo.AziWind.Value)
			row.AddStrValue(wInfo.AziWind.Uom)
			row.AddFloatValue(wInfo.VelWind.Value)
			row.AddStrValue(wInfo.VelWind.Uom)
			row.AddStrValue(wInfo.TypePrecip)
			row.AddFloatValue(wInfo.AmtPrecip.Value)
			row.AddStrValue(wInfo.AmtPrecip.Uom)
			row.AddStrValue(wInfo.CoverCloud)
			row.AddFloatValue(wInfo.CeilingCloud.Value)
			row.AddStrValue(wInfo.CeilingCloud.Uom)
			row.AddFloatValue(wInfo.CurrentSea.Value)
			row.AddStrValue(wInfo.CurrentSea.Uom)
			row.AddFloatValue(wInfo.AziCurrentSea.Value)
			row.AddStrValue(wInfo.AziCurrentSea.Uom)
			row.AddStrValue(wInfo.Comments)

			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildGasReadingInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "ReadingType", "MdTop",
		"MdTop_UoM", "MdBottom", "MdBottom_UoM", "TvdTop", "TvdTop_UoM",
		"TvdBottom", "TvdBottom_UoM", "GasHigh", "GasHigh_UoM", "GasLow",
		"GasLow_UoM", "Meth", "Meth_UoM", "Eth", "Eth_UoM", "Prop", "Prop_UoM",
		"IBut", "IBut_UoM", "NBut", "NBut_UoM", "IPent", "IPent_UoM", "NPent",
		"NPent_UoM", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].GasReadingInfos); s++ {
			gInfo := dReports[i].GasReadingInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(gInfo.DTim.Time)
			row.AddStrValue(gInfo.ReadingType)
			row.AddFloatValue(gInfo.MdTop.Value)
			row.AddStrValue(gInfo.MdTop.Uom)
			row.AddFloatValue(gInfo.MdBottom.Value)
			row.AddStrValue(gInfo.MdBottom.Uom)
			row.AddFloatValue(gInfo.TvdTop.Value)
			row.AddStrValue(gInfo.TvdTop.Uom)
			row.AddFloatValue(gInfo.TvdBottom.Value)
			row.AddStrValue(gInfo.TvdBottom.Uom)
			row.AddFloatValue(gInfo.GasHigh.Value)
			row.AddStrValue(gInfo.GasHigh.Uom)
			row.AddFloatValue(gInfo.GasLow.Value)
			row.AddStrValue(gInfo.GasLow.Uom)
			row.AddFloatValue(gInfo.Meth.Value)
			row.AddStrValue(gInfo.Meth.Uom)
			row.AddFloatValue(gInfo.Eth.Value)
			row.AddStrValue(gInfo.Eth.Uom)
			row.AddFloatValue(gInfo.Prop.Value)
			row.AddStrValue(gInfo.Prop.Uom)
			row.AddFloatValue(gInfo.Ibut.Value)
			row.AddStrValue(gInfo.Ibut.Uom)
			row.AddFloatValue(gInfo.NBut.Value)
			row.AddStrValue(gInfo.NBut.Uom)
			row.AddFloatValue(gInfo.IPent.Value)
			row.AddStrValue(gInfo.IPent.Uom)
			row.AddFloatValue(gInfo.NPent.Value)
			row.AddStrValue(gInfo.NPent.Uom)
			row.AddStrValue(gInfo.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildPerfInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTimOpen", "DTimClose", "MdTop", "MdTop_UoM",
		"MdBottom", "MdBottom_UoM", "TvdTop", "TvdTop_UoM", "TvdBottom",
		"TvdBottom_UoM", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].PerfInfos); s++ {
			pInfo := dReports[i].PerfInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(pInfo.DTimOpen.Time)
			row.AddTimeValue(pInfo.DTimClose.Time)
			row.AddFloatValue(pInfo.MdTop.Value)
			row.AddStrValue(pInfo.MdTop.Uom)
			row.AddFloatValue(pInfo.MdBottom.Value)
			row.AddStrValue(pInfo.MdBottom.Uom)
			row.AddFloatValue(pInfo.TvdTop.Value)
			row.AddStrValue(pInfo.TvdTop.Uom)
			row.AddFloatValue(pInfo.TvdBottom.Value)
			row.AddStrValue(pInfo.TvdBottom.Uom)
			row.AddStrValue(pInfo.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildStratInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "MDTopPlanned", "MDTopPlanned_UoM",
		"TvdTopPlanned", "TvdTopPlanned_UoM", "MdTop", "MdTop_UoM", "TvdTop", "TvdTop_UoM", "Description"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].StratInfos); s++ {
			sInfo := dReports[i].StratInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(sInfo.DTim.Time)
			row.AddFloatValue(sInfo.MdTopPlanned.Value)
			row.AddStrValue(sInfo.MdTopPlanned.Uom)
			row.AddFloatValue(sInfo.TvdTopPlanned.Value)
			row.AddStrValue(sInfo.TvdTopPlanned.Uom)
			row.AddFloatValue(sInfo.MdTop.Value)
			row.AddStrValue(sInfo.MdTop.Uom)
			row.AddFloatValue(sInfo.TvdTop.Value)
			row.AddStrValue(sInfo.TvdTop.Uom)
			row.AddStrValue(sInfo.Description)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildControlIncidentInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "MdInflow", "MdInflow_UoM",
		"TvdInflow", "TvdInflow_UoM", "Phase", "ProprietaryCode", "ETimLost", "ETimLost_UoM",
		"DTimRegained", "DiaBit", "DiaBit_UoM", "MdBit", "MdBit_UoM", "WTMud", "WTMud_UoM",
		"PorePressure", "PorePressure_UoM", "DiaCsgLast", "DiaCsgLast_UoM", "MDCsgLast", "MdCsgLast_UoM",
		"VolMudGained", "VolMudGained_UoM", "PresShutinCasing", "PresShutinCasing_UoM",
		"PresShutInDrill", "PresShutinDrill_UoM", "IncidentType", "KillingType", "Formation",
		"TempBottom", "TempBottom_UoM", "PresMaxChoke", "PresMaxChoke_UoM", "Description"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].ControlIncidentInfos); s++ {
			cInfo := dReports[i].ControlIncidentInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(cInfo.DTim.Time)
			row.AddFloatValue(cInfo.MdInflow.Value)
			row.AddStrValue(cInfo.MdInflow.Uom)
			row.AddFloatValue(cInfo.TvdInflow.Value)
			row.AddStrValue(cInfo.TvdInflow.Uom)
			row.AddStrValue(cInfo.Phase)
			row.AddStrValue(cInfo.ProprietaryCode)
			row.AddFloatValue(cInfo.ETimLost.Value)
			row.AddStrValue(cInfo.ETimLost.Uom)
			row.AddTimeValue(cInfo.DTimRegained.Time)
			row.AddFloatValue(cInfo.DiaBit.Value)
			row.AddStrValue(cInfo.DiaBit.Uom)
			row.AddFloatValue(cInfo.MdBit.Value)
			row.AddStrValue(cInfo.MdBit.Uom)
			row.AddFloatValue(cInfo.WtMud.Value)
			row.AddStrValue(cInfo.WtMud.Uom)
			row.AddFloatValue(cInfo.PorePressure.Value)
			row.AddStrValue(cInfo.PorePressure.Uom)
			row.AddFloatValue(cInfo.DiaCsgLast.Value)
			row.AddStrValue(cInfo.DiaCsgLast.Uom)
			row.AddFloatValue(cInfo.MdCsgLast.Value)
			row.AddStrValue(cInfo.MdCsgLast.Uom)
			row.AddFloatValue(cInfo.VolMudGained.Value)
			row.AddStrValue(cInfo.VolMudGained.Uom)
			row.AddFloatValue(cInfo.PresShutinCasing.Value)
			row.AddStrValue(cInfo.PresShutinCasing.Uom)
			row.AddFloatValue(cInfo.PresShutInDrill.Value)
			row.AddStrValue(cInfo.PresShutInDrill.Uom)
			row.AddStrValue(cInfo.IncidentType)
			row.AddStrValue(cInfo.KillingType)
			row.AddStrValue(cInfo.Formation)
			row.AddFloatValue(cInfo.TempBottom.Value)
			row.AddStrValue(cInfo.TempBottom.Uom)
			row.AddFloatValue(cInfo.PresMaxChoke.Value)
			row.AddStrValue(cInfo.PresMaxChoke.Uom)
			row.AddStrValue(cInfo.Description)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}
func buildEquipFailureInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "Md", "Md_UoM",
		"Tvd", "Tvd_UoM", "EquipClass", "ETimMissProduction", "ETimMissProduction_UoM",
		"DTimRepair", "Description"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].EquipFailureInfos); s++ {
			eInfo := dReports[i].EquipFailureInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(eInfo.DTim.Time)
			row.AddFloatValue(eInfo.Md.Value)
			row.AddStrValue(eInfo.Md.Uom)
			row.AddFloatValue(eInfo.Tvd.Value)
			row.AddStrValue(eInfo.Tvd.Uom)
			row.AddStrValue(eInfo.EquipClass)
			row.AddFloatValue(eInfo.ETimMissProduction.Value)
			row.AddStrValue(eInfo.ETimMissProduction.Uom)
			row.AddTimeValue(eInfo.DTimRepair.Time)
			row.AddStrValue(eInfo.Description)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRLithShowInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "MdTop", "MdTop_UoM",
		"MdBottom", "MdBottom_UoM", "TvdTop", "TvdTop_UoM", "TvdBottom", "TvdBottom_UoM",
		"Show", "Lithology"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].LithShowInfos); s++ {
			lInfo := dReports[i].LithShowInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(lInfo.DTim.Time)
			row.AddFloatValue(lInfo.MdTop.Value)
			row.AddStrValue(lInfo.MdTop.Uom)
			row.AddFloatValue(lInfo.MdBottom.Value)
			row.AddStrValue(lInfo.MdBottom.Uom)
			row.AddFloatValue(lInfo.TvdTop.Value)
			row.AddStrValue(lInfo.TvdTop.Uom)
			row.AddFloatValue(lInfo.TvdBottom.Value)
			row.AddStrValue(lInfo.TvdBottom.Uom)
			row.AddStrValue(lInfo.Show)
			row.AddStrValue(lInfo.Lithology)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRFormTestInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "RunNumber",
		"TestNumber", "MD", "MD_UoM", "TVD", "TVD_UoM", "PresPore", "PresPore_UoM",
		"FluidDensity", "FluidDensity_UoM", "HydrostaticPresBefore", "HydrostaticPresBefore_UoM",
		"LeakOffPressure", "LeakOffPressure_UoM", "GoodSeal", "MdSample", "MdSample_UoM",
		"DominateComponent", "DensityHC", "DensityHC_UoM", "VolumeSample", "VolumeSample_UoM", "Description"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].FormTestInfos); s++ {
			fInfo := dReports[i].FormTestInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(fInfo.DTim.Time)
			row.AddStrValue(fInfo.RunNumber)
			row.AddIntValue(fInfo.TestNumber)
			row.AddFloatValue(fInfo.Md.Value)
			row.AddStrValue(fInfo.Md.Uom)
			row.AddFloatValue(fInfo.Tvd.Value)
			row.AddStrValue(fInfo.Tvd.Uom)
			row.AddFloatValue(fInfo.PresPore.Value)
			row.AddStrValue(fInfo.PresPore.Uom)
			row.AddFloatValue(fInfo.FluidDensity.Value)
			row.AddStrValue(fInfo.FluidDensity.Uom)
			row.AddFloatValue(fInfo.HydrostaticPresBefore.Value)
			row.AddStrValue(fInfo.HydrostaticPresBefore.Uom)
			row.AddFloatValue(fInfo.LeakOffPressure.Value)
			row.AddStrValue(fInfo.LeakOffPressure.Uom)
			row.AddStrValue(strconv.FormatBool(fInfo.GoodSeal))
			row.AddFloatValue(fInfo.MdSample.Value)
			row.AddStrValue(fInfo.MdSample.Uom)
			row.AddStrValue(fInfo.DominateComponent)
			row.AddFloatValue(fInfo.DensityHC.Value)
			row.AddStrValue(fInfo.DensityHC.Uom)
			row.AddFloatValue(fInfo.VolumeSample.Value)
			row.AddStrValue(fInfo.VolumeSample.Uom)
			row.AddStrValue(fInfo.Description)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRWellTestInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "TestType", "TestNumber",
		"MDTop", "MDTop_UoM", "MDBottom", "MDBottom_UoM", "TvdTop", "TvdTop_UoM",
		"TvdBottom", "TvdBottom_UoM", "ChokeSize", "ChokeSize_UomM",
		"DensityOil", "DensityOil_UoM", "DensityWater", "DensityWater_UoM",
		"DensityGas", "DensityGas_UoM", "FlowrateOil", "FlowrateOil_UoM",
		"FlowrateWater", "FlowrateWater_UoM", "FlowrateGas", "FlowrateGas_UoM", "PresShutin",
		"PresShutin_UoM", "PresFlowing", "PresFlowing_UoM", "PresBottom", "PresBottom_UoM",
		"GasOilRatio", "GasOilRatio_UoM", "WaterOilRatio", "WaterOilRatio_UoM",
		"Chloride", "Chloride_UoM", "CarbonDioxide", "CarbonDioxide_UoM", "HydrogenSulfide",
		"HydrogenSulfide_UoM", "VolOilTotal", "VolOilTotal_UoM", "VolGasTotal", "VolGasTotal_UoM",
		"VolWaterTotal", "VolWaterTotal_UoM", "VolOilStored", "VolOilStored_UoM", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].WellTestInfos); s++ {
			wInfo := dReports[i].WellTestInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(wInfo.DTim.Time)
			row.AddStrValue(wInfo.TestType)
			row.AddIntValue(wInfo.TestNumber)
			row.AddFloatValue(wInfo.MdTop.Value)
			row.AddStrValue(wInfo.MdTop.Uom)
			row.AddFloatValue(wInfo.MdBottom.Value)
			row.AddStrValue(wInfo.MdBottom.Uom)
			row.AddFloatValue(wInfo.TvdTop.Value)
			row.AddStrValue(wInfo.TvdTop.Uom)
			row.AddFloatValue(wInfo.TvdBottom.Value)
			row.AddStrValue(wInfo.TvdBottom.Uom)
			row.AddFloatValue(wInfo.ChokeSize.Value)
			row.AddStrValue(wInfo.ChokeSize.Uom)
			row.AddFloatValue(wInfo.DensityOil.Value)
			row.AddStrValue(wInfo.DensityOil.Uom)
			row.AddFloatValue(wInfo.DensityWater.Value)
			row.AddStrValue(wInfo.DensityWater.Uom)
			row.AddFloatValue(wInfo.DensityGas.Value)
			row.AddStrValue(wInfo.DensityGas.Uom)
			row.AddFloatValue(wInfo.FlowRateOil.Value)
			row.AddStrValue(wInfo.FlowRateOil.Uom)
			row.AddFloatValue(wInfo.FlowRateWater.Value)
			row.AddStrValue(wInfo.FlowRateWater.Uom)
			row.AddFloatValue(wInfo.FlowRateGas.Value)
			row.AddStrValue(wInfo.FlowRateGas.Uom)
			row.AddFloatValue(wInfo.PresShutIn.Value)
			row.AddStrValue(wInfo.PresShutIn.Uom)
			row.AddFloatValue(wInfo.PresFlowing.Value)
			row.AddStrValue(wInfo.PresFlowing.Uom)
			row.AddFloatValue(wInfo.PresBottom.Value)
			row.AddStrValue(wInfo.PresBottom.Uom)
			row.AddFloatValue(wInfo.GoR.Value)
			row.AddStrValue(wInfo.GoR.Uom)
			row.AddFloatValue(wInfo.WaterOilRatio.Value)
			row.AddStrValue(wInfo.WaterOilRatio.Uom)
			row.AddFloatValue(wInfo.Chloride.Value)
			row.AddStrValue(wInfo.Chloride.Uom)
			row.AddFloatValue(wInfo.CarbonDioxide.Value)
			row.AddStrValue(wInfo.CarbonDioxide.Uom)
			row.AddFloatValue(wInfo.HydrogenSulfide.Value)
			row.AddStrValue(wInfo.HydrogenSulfide.Uom)
			row.AddFloatValue(wInfo.VolOilTotal.Value)
			row.AddStrValue(wInfo.VolOilTotal.Uom)
			row.AddFloatValue(wInfo.VolGasTotal.Value)
			row.AddStrValue(wInfo.VolGasTotal.Uom)
			row.AddFloatValue(wInfo.VolWaterTotal.Value)
			row.AddStrValue(wInfo.VolWaterTotal.Uom)
			row.AddFloatValue(wInfo.VolOilStored.Value)
			row.AddStrValue(wInfo.VolOilStored.Uom)
			row.AddStrValue(wInfo.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRCoreInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim",
		"CoreNumber", "MDTop", "MDTop_UoM", "MDBottom", "MDBottom_UoM",
		"TvdTop", "TvdTop_UoM", "TvdBottom", "TvdBottom_UoM", "LenRecovered",
		"LenRecovered_UoM", "RecoverPc", "RecoverPc_UoM",
		"LenBarrel", "LenBarrel_UoM", "InnerBarrelType", "CoreDescription"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].CoreInfos); s++ {
			cInfo := dReports[i].CoreInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(cInfo.DTim.Time)
			row.AddStrValue(cInfo.CoreNumber)
			row.AddFloatValue(cInfo.MDTop.Value)
			row.AddStrValue(cInfo.MDTop.Uom)
			row.AddFloatValue(cInfo.MDBottom.Value)
			row.AddStrValue(cInfo.MDBottom.Uom)
			row.AddFloatValue(cInfo.TvdTop.Value)
			row.AddStrValue(cInfo.TvdTop.Uom)
			row.AddFloatValue(cInfo.TvdBottom.Value)
			row.AddStrValue(cInfo.TvdBottom.Uom)
			row.AddFloatValue(cInfo.LenRecovered.Value)
			row.AddStrValue(cInfo.LenRecovered.Uom)
			row.AddFloatValue(cInfo.RecoverPC.Value)
			row.AddStrValue(cInfo.RecoverPC.Uom)
			row.AddFloatValue(cInfo.LenBarrel.Value)
			row.AddStrValue(cInfo.LenBarrel.Uom)
			row.AddStrValue(cInfo.InnerBarrelType)
			row.AddStrValue(cInfo.CoreDescription)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRLoginfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim",
		"RunNumber", "ServiceCompany", "Service", "MDTop", "MDTop_UoM", "MDBottom",
		"MDBottom_UoM", "TvdTop", "TvdTop_UoM", "TvdBottom", "TvdBottom_UoM", "Tool",
		"TempBHCT", "TempBHCT_UoM", "TempBHST", "TempBHST_UoM", "ETimStatic", "ETimStatic_UoM",
		"MDTempTool", "MDTempTool_UoM", "TvdTempTool", "TvdTempTool_UoM", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].LogInfos); s++ {
			lInfo := dReports[i].LogInfos[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(lInfo.DTim.Time)
			row.AddStrValue(lInfo.RunNumber)
			row.AddStrValue(lInfo.ServiceCompany)
			row.AddStrValue(lInfo.Service)
			row.AddFloatValue(lInfo.MdTop.Value)
			row.AddStrValue(lInfo.MdTop.Uom)
			row.AddFloatValue(lInfo.MdBottom.Value)
			row.AddStrValue(lInfo.MdBottom.Uom)
			row.AddFloatValue(lInfo.TvdTop.Value)
			row.AddStrValue(lInfo.TvdTop.Uom)
			row.AddFloatValue(lInfo.TvdBottom.Value)
			row.AddStrValue(lInfo.TvdBottom.Uom)
			row.AddStrValue(lInfo.Tool)
			row.AddFloatValue(lInfo.TempBHCt.Value)
			row.AddStrValue(lInfo.TempBHCt.Uom)
			row.AddFloatValue(lInfo.TempBHST.Value)
			row.AddStrValue(lInfo.TempBHST.Uom)
			row.AddFloatValue(lInfo.ETimStatic.Value)
			row.AddStrValue(lInfo.ETimStatic.Uom)
			row.AddFloatValue(lInfo.MdTempTool.Value)
			row.AddStrValue(lInfo.MdTempTool.Uom)
			row.AddFloatValue(lInfo.TvdTempTool.Value)
			row.AddStrValue(lInfo.TvdTempTool.Uom)
			row.AddStrValue(lInfo.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRActivities(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTimStart", "DTimEnd",
		"MD", "MD_UoM", "Tvd", "Tvd_UoM", "Phase", "ProprietaryCode",
		"Conveyance", "MdHoleStart", "MdHoleStart_UoM", "State", "StateDetailActivity", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].Activities); s++ {
			act := dReports[i].Activities[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(act.DTimStart.Time)
			row.AddTimeValue(act.DTimeEnd.Time)
			row.AddFloatValue(act.Md.Value)
			row.AddStrValue(act.Md.Uom)
			row.AddFloatValue(act.Tvd.Value)
			row.AddStrValue(act.Tvd.Uom)
			row.AddStrValue(act.Phase)
			row.AddStrValue(act.ProprietaryCode)
			row.AddStrValue(act.Conveyance)
			row.AddFloatValue(act.MdHoleStart.Value)
			row.AddStrValue(act.MdHoleStart.Uom)
			row.AddStrValue(act.State)
			row.AddStrValue(act.StateDetailActivity)
			row.AddStrValue(act.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRSurveyStation(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTim", "MD", "MD_UoM", "Tvd", "Tvd_UoM",
		"Incl", "Incl_UoM", "Azi", "Azi_UoM",
		"Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].SurveyStations); s++ {
			sStation := dReports[i].SurveyStations[s]
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(sStation.DTim.Time)
			row.AddFloatValue(sStation.Md.Value)
			row.AddStrValue(sStation.Md.Uom)
			row.AddFloatValue(sStation.Tvd.Value)
			row.AddStrValue(sStation.Tvd.Uom)
			row.AddFloatValue(sStation.Incl.Value)
			row.AddStrValue(sStation.Incl.Uom)
			row.AddFloatValue(sStation.Azi.Value)
			row.AddStrValue(sStation.Azi.Uom)
			row.AddStrValue(sStation.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRPorePressure(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "ReadingKind", "EquivalentMudWeight",
		"EquivalentMudWeight_UoM", "DTim", "MD", "MD_UoM", "Tvd", "Tvd_UoM", "Comment"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].PorePressures); s++ {

			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			porePressure := dReports[i].PorePressures[s]
			row.AddStrValue(porePressure.ReadingKind)
			row.AddFloatValue(porePressure.EquivalentMudWeight.Value)
			row.AddStrValue(porePressure.EquivalentMudWeight.Uom)
			row.AddTimeValue(porePressure.DTim.Time)
			row.AddFloatValue(porePressure.Md.Value)
			row.AddStrValue(porePressure.Md.Uom)
			row.AddFloatValue(porePressure.Tvd.Value)
			row.AddStrValue(porePressure.Tvd.Uom)
			row.AddStrValue(porePressure.Comment)

			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRDrillReportInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	headers = []string{"DataUUID", "DrillReportUid",
		"DrillReportUidWell", "DrillReportUidWellbore",
		"NameWell", "NameWellNamingSystem",
		"NameWellbore", "NameWellboreNamingSystem", "ReportName", "ReportDTimStart",
		"ReportDTimEnd", "VersionKind",
		"CreatedDate",
		"ExtendedReportTime", "ExtendedReport", "Comment"}
	dSet.HeadersName = headers
	dSet.Name = datasetName
	for i := 0; i < len(dReports); i++ {
		//need to add one row entry per wellbore alias
		for s := 0; s < len(dReports[i].WellboreAliases); s++ {
			row := common.RowData{}
			row.AddStrValue(dReports[i].DataIdentification.UUid)
			row.AddStrValue(dReports[i].Uid)
			row.AddStrValue(dReports[i].UidWell)
			row.AddStrValue(dReports[i].UidWellbore)
			row.AddStrValue(dReports[i].WellAlias.Name)
			row.AddStrValue(dReports[i].WellAlias.NamingSystem)
			row.AddStrValue(dReports[i].WellboreAliases[s].Name)
			row.AddStrValue(dReports[i].WellboreAliases[s].NamingSystem)
			row.AddStrValue(dReports[i].Name)
			row.AddTimeValue(dReports[i].DTimStart.Time)
			row.AddTimeValue(dReports[i].DTimEnd.Time)
			row.AddStrValue(dReports[i].VersionKind)
			row.AddTimeValue(dReports[i].CreatedDate.Time)
			row.AddTimeValue(dReports[i].ExtendedReport.DTim.Time)
			row.AddStrValue(dReports[i].ExtendedReport.Value)
			row.AddStrValue(dReports[i].Comment)
			rows = append(rows, row)
		}

	}
	dSet.Rows = rows
	return dSet
}

func buildDDRWellboreInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTimSpud", "DTimPreSpud", "DateDrillComplete",
		"DaysAhead", "Operator", "DrillContractor", "RigAlias", "RigAliasNamingSystem", "DaysBehind"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		for x := 0; x < len(dReports[i].WellboreInfo.RigAliases); x++ {
			idSet = buildDDRRowIdentification(dReports[i], datasetName)
			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddTimeValue(dReports[i].WellboreInfo.DTimSpud.Time)
			row.AddTimeValue(dReports[i].WellboreInfo.DTimPreSpud.Time)
			row.AddTimeValue(dReports[i].WellboreInfo.DateDrillComplete.Time)
			row.AddFloatValue(dReports[i].WellboreInfo.DaysAhead.float64)
			row.AddStrValue(dReports[i].WellboreInfo.Operator)
			row.AddStrValue(dReports[i].WellboreInfo.DrillContractor)
			row.AddStrValue(dReports[i].WellboreInfo.RigAliases[x].Name)
			row.AddStrValue(dReports[i].WellboreInfo.RigAliases[x].NamingSystem)
			row.AddFloatValue(dReports[i].WellboreInfo.DaysBehind.float64)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRCementStages(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "DTimPumpStart", "DTimPumpEnd",
		"JobType", "CasingStrDia", "CasingStrDia_UoM", "Comments", "VolReturns", "VolReturns_UoM",
		"TypeFluid", "DescFluid", "RatioMixWater", "RatioMixWater_UoM", "Density", "Density_UoM",
		"VolPumped", "VolPumped_UoM", "YP", "YP_UoM", "ETimThickening", "ETimThickening_UoM",
		"PCFreeWater", "PCFreeWater_UoM", "FluidComments", "DTimPresReleased",
		"AnnFlowAfter", "TopPlug", "BotPlug", "PlugBumped", "PresBump", "PresBump_UoM", "FloatHeld",
		"Reciprocated", "Rotated"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for x := 0; x < len(dReports[i].CementStages); x++ {
			cementStage := dReports[i].CementStages[x]
			//need to generate one row per cementing fluid..
			for s := 0; s < len(cementStage.CementingFluids); s++ {

				cementingFluid := cementStage.CementingFluids[s]
				row := idSet.Rows[0]
				row.AddStrValue(recordId)
				row.AddTimeValue(cementStage.DTimPumpStart.Time)
				row.AddTimeValue(cementStage.DTimPumpEnd.Time)
				row.AddStrValue(cementStage.JobType)
				row.AddFloatValue(cementStage.CasingStrDia.Value)
				row.AddStrValue(cementStage.CasingStrDia.Uom)
				row.AddStrValue(cementStage.Comments)
				row.AddFloatValue(cementStage.VolReturns.Value)
				row.AddStrValue(cementStage.VolReturns.Uom)
				row.AddStrValue(cementingFluid.TypeFluid)
				row.AddStrValue(cementingFluid.DescFluid)
				row.AddFloatValue(cementingFluid.RatioMixWater.Value)
				row.AddStrValue(cementingFluid.RatioMixWater.Uom)
				row.AddFloatValue(cementingFluid.Density.Value)
				row.AddStrValue(cementingFluid.Density.Uom)
				row.AddFloatValue(cementingFluid.VolPumped.Value)
				row.AddStrValue(cementingFluid.VolPumped.Uom)
				row.AddFloatValue(cementingFluid.Yp.Value)
				row.AddStrValue(cementingFluid.Yp.Uom)
				row.AddFloatValue(cementingFluid.ETimThickening.Value)
				row.AddStrValue(cementingFluid.ETimThickening.Uom)
				row.AddFloatValue(cementingFluid.PCFreeWater.Value)
				row.AddStrValue(cementingFluid.PCFreeWater.Uom)
				row.AddStrValue(cementingFluid.Comments)
				row.AddTimeValue(cementStage.DTimPresReleased.Time)
				row.AddStrValue(strconv.FormatBool(cementStage.AnnFlowAfter))
				row.AddStrValue(strconv.FormatBool(cementStage.TopPlug))
				row.AddStrValue(strconv.FormatBool(cementStage.BotPlug))
				row.AddStrValue(strconv.FormatBool(cementStage.PlugBumped))
				row.AddFloatValue(cementStage.PresBump.Value)
				row.AddStrValue(cementStage.PresBump.Uom)
				row.AddStrValue(strconv.FormatBool(cementStage.FloatHeld))
				row.AddStrValue(strconv.FormatBool(cementStage.Reciprocated))
				row.AddStrValue(strconv.FormatBool(cementStage.Rotated))
				rows = append(rows, row)
			}
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRCasingLinerTubings(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "Type", "ID", "ID_UoM",
		"OD", "OD_UoM", "Weight", "Weight_UoM", "Grade",
		"Connection", "Length", "Length_UoM", "MdTop", "MdTop_UoM",
		"MdBottom", "MdBottom_UoM", "CasingType", "Description",
		"DTimStart", "DTimEnd", "Comment"}
	headers = append(idSet.HeadersName, headers...)
	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for s := 0; s < len(dReports[i].CasingLinerTubings); s++ {
			casingL := dReports[i].CasingLinerTubings[s]

			row := idSet.Rows[0]
			row.AddStrValue(recordId)
			row.AddStrValue(casingL.Type)
			row.AddFloatValue(casingL.Id.Value)
			row.AddStrValue(casingL.Id.Uom)
			row.AddFloatValue(casingL.Od.Value)
			row.AddStrValue(casingL.Od.Uom)
			row.AddFloatValue(casingL.Weight.Value)
			row.AddStrValue(casingL.Weight.Uom)
			row.AddStrValue(casingL.Grade)
			row.AddStrValue(casingL.Connection)
			row.AddFloatValue(casingL.Length.Value)
			row.AddStrValue(casingL.Length.Uom)
			row.AddFloatValue(casingL.MdTop.Value)
			row.AddStrValue(casingL.MdTop.Uom)
			row.AddFloatValue(casingL.MdBottom.Value)
			row.AddStrValue(casingL.MdBottom.Uom)
			row.AddStrValue(casingL.CasingLinerTubingRun.CasingType)
			row.AddStrValue(casingL.CasingLinerTubingRun.Description)
			row.AddTimeValue(casingL.CasingLinerTubingRun.DTimStart.Time)
			row.AddTimeValue(casingL.CasingLinerTubingRun.DTimEnd.Time)
			row.AddStrValue(casingL.Comment)
			rows = append(rows, row)
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRBitRecords(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"RecordId", "NumBitRun", "Numbit",
		"DiaBit", "DiaBit_UoM", "Manufacturer", "CodeMfg", "DullGrade",
		"CodelADC", "CondFinalInner", "CondFinalOuter", "CondFinalDull",
		"CondFinalLocation", "CondFinalBearing", "CondFinalGauge", "CondFinalOther",
		"CondFinalReason", "BitRun_ETimOpBit", "BitRun_ETimOpBit_UoM",
		"MdHoleStart", "MdHoleStart_UoM", "MdHoleStop", "MdHoleStop_UoM",
		"RopAv", "RopAV_UoM", "MdHoleMadeRun", "MdHoleMadeRun_UoM", "HrsDrilled",
		"HrsDrilled_UoM", "HrsDrilledRun", "HrsDrilledRun_UoM", "MdTotHoleMade", "MdTotHoleMade_UoM",
		"TotHrsDrilled", "TotHrsDrilled_UoM", "TotRop", "TotRop_UoM",
		"NumNozzle", "DiaNozzle", "DiaNozzle_UoM"}
	headers = append(idSet.HeadersName, headers...)

	for i := 0; i < len(dReports); i++ {
		recordId := common.CreateUUID()
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		for x := 0; x < len(dReports[i].BitRecords); x++ {
			bitRecord := dReports[i].BitRecords[x]
			//need to loop through the noozles to create one entry per row
			for s := 0; s < len(bitRecord.Nozzles); s++ {
				row := idSet.Rows[0]

				row.AddStrValue(recordId)
				row.AddStrValue(bitRecord.NumBitRun)
				row.AddStrValue(bitRecord.NumBit)
				row.AddFloatValue(bitRecord.DiaBit.Value)
				row.AddStrValue(bitRecord.DiaBit.Uom)
				row.AddStrValue(bitRecord.Manufacturer)
				row.AddStrValue(bitRecord.CodeMfg)
				row.AddStrValue(bitRecord.DullGrade)
				row.AddStrValue(bitRecord.CodeIADC)
				row.AddIntValue(bitRecord.CondFinalInner)
				row.AddIntValue(bitRecord.CondFinalOuter)
				row.AddStrValue(bitRecord.CondFinalDull)
				row.AddStrValue(bitRecord.CondFinalLocation)
				row.AddStrValue(bitRecord.CondFinalBearing)
				row.AddStrValue(bitRecord.CondFinalGauge)
				row.AddStrValue(bitRecord.CondFinalOther)
				row.AddStrValue(bitRecord.CondFinalReason)
				row.AddFloatValue(bitRecord.BitRun.ETimOpBit.Value)
				row.AddStrValue(bitRecord.BitRun.ETimOpBit.Uom)
				row.AddFloatValue(bitRecord.BitRun.MDHoleStart.Value)
				row.AddStrValue(bitRecord.BitRun.MDHoleStart.Uom)
				row.AddFloatValue(bitRecord.BitRun.MDHoleStop.Value)
				row.AddStrValue(bitRecord.BitRun.MDHoleStop.Uom)
				row.AddFloatValue(bitRecord.BitRun.RopAv.Value)
				row.AddStrValue(bitRecord.BitRun.RopAv.Uom)
				row.AddFloatValue(bitRecord.BitRun.MDHoleMadeRun.Value)
				row.AddStrValue(bitRecord.BitRun.MDHoleMadeRun.Uom)
				row.AddFloatValue(bitRecord.BitRun.HrsDrilled.Value)
				row.AddStrValue(bitRecord.BitRun.HrsDrilled.Uom)
				row.AddFloatValue(bitRecord.BitRun.HrsDrilledRun.Value)
				row.AddStrValue(bitRecord.BitRun.HrsDrilledRun.Uom)
				row.AddFloatValue(bitRecord.BitRun.MdTotalHoleMade.Value)
				row.AddStrValue(bitRecord.BitRun.MdTotalHoleMade.Uom)
				row.AddFloatValue(bitRecord.BitRun.TotHrsDrilled.Value)
				row.AddStrValue(bitRecord.BitRun.TotHrsDrilled.Uom)
				row.AddFloatValue(bitRecord.BitRun.TotRop.Value)
				row.AddStrValue(bitRecord.BitRun.TotRop.Uom)
				row.AddIntValue(bitRecord.Nozzles[s].NumNozzle)
				row.AddFloatValue(bitRecord.Nozzles[s].DiaNozzle.Value)
				row.AddStrValue(bitRecord.Nozzles[s].DiaNozzle.Uom)
				rows = append(rows, row)
			}
		}
	}
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

func buildDDRStatusInfo(dReports []DrillReport, datasetName string) common.DataSet {
	var dSet common.DataSet
	var headers []string
	var rows []common.RowData
	dSet.Name = datasetName
	idSet := buildDDRRowIdentification(dReports[0], datasetName)
	headers = []string{"ReportNo", "DTim", "MD", "MD_UoM",
		"TVD", "TVD_UoM",
		"MDPlugTop", "MDPlugTop_UoM",
		"DiaHole", "DiaHole_UoM",
		"DTimDiaHoleStart",
		"MDDiaHoleStart", "MDDiaHoleStart_UoM",
		"DiaPilot", "DiaPilot_UoM",
		"MDDiaPilotPlan", "MDDiaPilotPlan_UoM",
		"TVDDiaPilotPlan", "TVDDiaPilotPlan_Uom",
		"TypeWellbore", "PrimaryConveyance",
		"MDKickOff", "MDKickOff_UoM",
		"TVDKickOff", "TVDKickOff_UoM",
		"StrengthForm", "StrengthForm_UoM",
		"MDStrengthForm", "MDStrengthForm_UoM",
		"TVDStrengthForm", "TVDStrengthForm_UoM",
		"DiaCsgLast", "DiaCsgLast_UoM",
		"MDCsgLast", "MDCsgLast_UoM",
		"TVDCsgLast", "TVDCsgLast_UoM",
		"PresTestType",
		"MDPlanned", "MDPlanned_UoM",
		"DistDrill", "DistDrill_UoM",
		"ElevKelly", "ElevKelly_UoM",
		"WellheadElevation", "WellheadElevation_UoM",
		"WaterDepth", "WaterDepth_UoM",
		"Sum24Hr",
		"Forecast24Hr",
		"RopCurrent", "RopCurrent_UoM",
		"TightWell", "HPHT",
		"AvgPresBH", "AvgPresBH_UoM",
		"AvgTempBH", "AvgTempBH_UoM",
		"FixedRig"}
	headers = append(idSet.HeadersName, headers...)
	for i := 0; i < len(dReports); i++ {
		idSet = buildDDRRowIdentification(dReports[i], datasetName)
		row := idSet.Rows[0]
		sInfo := dReports[i].StatusInfo
		row.AddIntValue(sInfo.Reportnumber)
		row.AddTimeValue(sInfo.DTim.Time)
		row.AddFloatValue(sInfo.Md.Value)
		row.AddStrValue(sInfo.Md.Uom)
		row.AddFloatValue(sInfo.Tvd.Value)
		row.AddStrValue(sInfo.Tvd.Uom)
		row.AddFloatValue(sInfo.MdPlugTop.Value)
		row.AddStrValue(sInfo.MdPlugTop.Uom)
		row.AddFloatValue(sInfo.DiaHole.Value)
		row.AddStrValue(sInfo.DiaHole.Uom)
		row.AddTimeValue(sInfo.DTimDiaHoleStart.Time)
		row.AddFloatValue(sInfo.MdDiaHoleStart.Value)
		row.AddStrValue(sInfo.MdDiaHoleStart.Uom)
		row.AddFloatValue(sInfo.DiaPilot.Value)
		row.AddStrValue(sInfo.DiaPilot.Uom)
		row.AddFloatValue(sInfo.MdDiaPilotPlan.Value)
		row.AddStrValue(sInfo.MdDiaPilotPlan.Uom)
		row.AddFloatValue(sInfo.TVDDiaPilotPlan.Value)
		row.AddStrValue(sInfo.TVDDiaPilotPlan.Uom)
		row.AddStrValue(sInfo.TypeWellBore)
		row.AddStrValue(sInfo.PrimaryConveyance)
		row.AddFloatValue(sInfo.MdKickoff.Value)
		row.AddStrValue(sInfo.MdKickoff.Uom)
		row.AddFloatValue(sInfo.TvdKickoff.Value)
		row.AddStrValue(sInfo.TvdKickoff.Uom)
		row.AddFloatValue(sInfo.StrengthForm.Value)
		row.AddStrValue(sInfo.StrengthForm.Uom)
		row.AddFloatValue(sInfo.MdStrengthForm.Value)
		row.AddStrValue(sInfo.MdStrengthForm.Uom)
		row.AddFloatValue(sInfo.TvdStrengthForm.Value)
		row.AddStrValue(sInfo.TvdStrengthForm.Uom)
		row.AddFloatValue(sInfo.DiaCasingLast.Value)
		row.AddStrValue(sInfo.DiaCasingLast.Uom)
		row.AddFloatValue(sInfo.MdCasingLast.Value)
		row.AddStrValue(sInfo.MdCasingLast.Uom)
		row.AddFloatValue(sInfo.TvdCasingLast.Value)
		row.AddStrValue(sInfo.TvdCasingLast.Uom)
		row.AddStrValue(sInfo.PressTestType)
		row.AddFloatValue(sInfo.MdPlanned.Value)
		row.AddStrValue(sInfo.MdPlanned.Uom)
		row.AddFloatValue(sInfo.DistDrilled.Value)
		row.AddStrValue(sInfo.DistDrilled.Uom)
		row.AddFloatValue(sInfo.ElevKelly.Value)
		row.AddStrValue(sInfo.ElevKelly.Uom)
		row.AddFloatValue(sInfo.WellheadElevation.Value)
		row.AddStrValue(sInfo.WellheadElevation.Uom)
		row.AddFloatValue(sInfo.WaterDepth.Value)
		row.AddStrValue(sInfo.WaterDepth.Uom)
		row.AddStrValue(sInfo.Sum24Hr)
		row.AddStrValue(sInfo.Forecast24Hr)
		row.AddFloatValue(sInfo.RopCurrent.Value)
		row.AddStrValue(sInfo.RopCurrent.Uom)
		row.AddStrValue(strconv.FormatBool(sInfo.TightWell))
		row.AddStrValue(strconv.FormatBool(sInfo.HPHT))
		row.AddFloatValue(sInfo.AvgPresBH.Value)
		row.AddStrValue(sInfo.AvgPresBH.Uom)
		row.AddFloatValue(sInfo.AvgTempBH.Value)
		row.AddStrValue(sInfo.AvgTempBH.Uom)
		row.AddStrValue(strconv.FormatBool(sInfo.FixedRig))
		rows = append(rows, row)
	}
	//start off with the data that is generic
	dSet.HeadersName = headers
	dSet.Rows = rows
	return dSet
}

//Function will build a base dataset extrating the essential information from a drillreport such as
//wellbore name, uid, times and so on. Other dataset functions can then use this base dataset to append additional
//data, hence this will be used as indentification columns
func buildDDRRowIdentification(dReport DrillReport, dataSetName string) common.DataSet {
	var dSet common.DataSet
	dSet = common.DataSet{}
	dSet.Name = dataSetName
	dSet.HeadersName = []string{"DataUUID", "DrillReportUid",
		"DrillReportUidWell", "DrillReportUidWellbore", "NameWell",
		"NameWellbore", "ReportName", "ReportDTimStart",
		"ReportDTimEnd"}
	row := common.RowData{}
	row.AddStrValue(dReport.DataIdentification.UUid)
	row.AddStrValue(dReport.Uid)
	row.AddStrValue(dReport.UidWell)
	row.AddStrValue(dReport.UidWellbore)
	row.AddStrValue(dReport.NameWell)
	row.AddStrValue(dReport.NameWellbore)
	row.AddStrValue(dReport.Name)
	row.AddTimeValue(dReport.DTimStart.Time)
	row.AddTimeValue(dReport.DTimEnd.Time)
	dSet.Rows = append(dSet.Rows, row)
	return dSet
}

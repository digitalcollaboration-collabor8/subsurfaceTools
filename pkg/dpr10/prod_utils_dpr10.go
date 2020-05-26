/*
* @Author: magsv
* @Date:   2018-02-22 10:46:22
* @Last Modified by:   magsv
* @Last Modified time: 2018-08-07 14:35:15
 */
package dpr10

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"go.uber.org/zap"
)

//Parses a MPRML xml file to a struct for further processing
func ParseProdXMLFile(inputFile string) (WITSMLComposite, error) {
	var witsmlComposite WITSMLComposite
	var err error
	var data []byte
	if data, err = common.ReadFile(inputFile); err != nil {
		return witsmlComposite, err
	}
	if err = xml.Unmarshal(data, &witsmlComposite); err != nil {
		return witsmlComposite, err
	}
	//just add the uid and file data for this set of objects
	witsmlComposite.DataIdentification = createDataId(inputFile)
	if len(witsmlComposite.ProdVolumeSet.ProdVolumes) > 0 {
		witsmlComposite.DataIdentification.ReportStart = witsmlComposite.ProdVolumeSet.ProdVolumes[0].DateStart.Time
		witsmlComposite.DataIdentification.ReportEnd = witsmlComposite.ProdVolumeSet.ProdVolumes[0].DateEnd.Time
	}

	return witsmlComposite, nil
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

//Reads a set of xml mprml files from a given folder path and
//parses them into a list of struct objects to be used for further processing
func ReadProdXMLFiles2Struct(folderPath string) ([]WITSMLComposite, error) {
	var err error
	var objects []WITSMLComposite
	var files []string
	var data WITSMLComposite
	start := time.Now()

	if files, err = common.GetFilesWithExtension(folderPath, "*.xml"); err != nil {
		zap.S().Error("Failed in reading file in folder:", err.Error())
		return objects, err
	}
	fileSearchTook := time.Since(start)
	xmlStart := time.Now()
	for i := 0; i < len(files); i++ {
		zap.S().Info("Processing xml file:", files[i])
		if data, err = ParseProdXMLFile(files[i]); err != nil {
			zap.S().Errorf("Failed in parsing xml file:%s\n", err.Error())
			return objects, err
		}
		objects = append(objects, data)
	}
	xmlParseTook := time.Since(xmlStart)

	zap.S().Infof("File search took:%s", fileSearchTook)
	zap.S().Infof("XML parse took:%s", xmlParseTook)
	return objects, nil
}

//Parses through a set of mprml objects and organises them by facility kind returning
//a map of facilities where the key is the facility kind
func OrganiseFacilitiesByKind_DPR10(objects []WITSMLComposite) map[string][]Facility {
	zap.S().Infof("Organizing facilities by kind. length of objects:%d", len(objects))
	facilities := make(map[string][]Facility)
	for a := 0; a < len(objects); a++ {
		for i := 0; i < len(objects[a].ProdVolumeSet.ProdVolumes); i++ {
			for x := 0; x < len(objects[a].ProdVolumeSet.ProdVolumes[i].Facilities); x++ {
				//add internal flowuids to all flow entities
				for s := 0; s < len(objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x].Flow); s++ {
					//generate the internal flowUid
					internalUid := common.CreateUUID()
					objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x].Flow[s].InternalFlowUid = internalUid
				}
				//check if the facility type exists in the map..
				facilityKind := objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x].Name.Kind
				//zap.S().Infof("Processing facility kind:%s", facilityKind)
				objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x].DataIdentification.DataIdentification = objects[a].DataIdentification
				objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x].DataIdentification.Installation = objects[a].ProdVolumeSet.ProdVolumes[i].Installation

				if len(facilities[facilityKind]) == 0 {
					//not there add it
					facilities[facilityKind] = []Facility{objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x]}
				} else {
					//it is there we just append it to the existing list
					facilities[facilityKind] = append(facilities[facilityKind], objects[a].ProdVolumeSet.ProdVolumes[i].Facilities[x])
				}
			}
		}
	}
	return facilities
}

//Extracts data from the report files as read from the disk
func extractFileInformation(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "FileName", "FilePath", "ReadTime", "DocumentName"}
	for i := 0; i < len(objects); i++ {
		row := common.RowData{}
		row.AddStrValue(objects[i].DataIdentification.UUid)
		row.AddStrValue(objects[i].DataIdentification.FileName)
		row.AddStrValue(objects[i].DataIdentification.FilePath)
		row.AddTimeValue(objects[i].DataIdentification.ReadTime)
		row.AddStrValue(objects[i].DocumentInfo.DocumentName.Txt)
		rows = append(rows, row)

	}
	dataSet.Rows = rows
	return dataSet

}

func BuildCSVFileForProduction(path string, objects []WITSMLComposite, discriminator string) error {
	dataSets := BuildDataSetsDPR_10(objects)
	return common.DatasetsToCsv(dataSets, path, discriminator)
}

//Builds an excel output file for a given list of mprml objects
func BuildXLSFileForProduction(path string, objects []WITSMLComposite, oneFilePerSheet bool, appendTimeInName bool) error {
	dataSets := BuildDataSetsDPR_10(objects)
	return common.CreateWorkbookFromDataSet(path, dataSets, oneFilePerSheet, appendTimeInName)
}

func BuildDataSetsDPR_10(objects []WITSMLComposite) []common.DataSet {
	var datasets []common.DataSet
	var groupedFacilities map[string][]Facility
	var dataSetMap map[string]*common.DataSet
	groupedFacilities = OrganiseFacilitiesByKind_DPR10(objects)

	//build the report file info
	datasets = append(datasets, extractFileInformation("REPORT_FILE_INFO", objects))
	datasets = append(datasets, extractFacilityInfo("FACILITIES", groupedFacilities))
	datasets = append(datasets, extractFlowInformation("FLOW_INFO", objects))
	//add the document info data
	datasets = append(datasets, extractDocumentInfo_DPR10("DOCUMENT_INFO", objects))
	datasets = append(datasets, extractOpHSE_DPR10("OP_HSE", objects))
	datasets = append(datasets, extractPersonell_DPR10("CREW", objects))
	datasets = append(datasets, extract_OpComments_DPR10("OP_COMMENTS", objects))
	datasets = append(datasets, extractLostProduction_DPR10("LOST_PRODUCTION", objects))
	datasets = append(datasets, extractWaterCleaning_DPR10("WATER_CLEANING", objects))
	datasets = append(datasets, extractWellTests_DPR10("WELLTESTS", objects))
	for kind, facilities := range groupedFacilities {
		zap.S().Infof("Building sheet for facility kind:%s", kind)
		if kind != "wellhead" && kind != "bottomhole" {
			dataSetMap = extractFacilityProdVolumes(strings.ToUpper(kind), facilities)
			//now have a map with keys where each key is a combination of facility kind + period kind, add each entity
		} else {
			dataSetMap = extractFacilityProdVolumesBHPAndWHP(strings.ToUpper(kind), facilities)
		}
		for key, value := range dataSetMap {
			zap.S().Infof("Adding dataset with key:%s", key)
			datasets = append(datasets, *value)
		}
	}
	return datasets

}

func BuildJsonFileForProduction(path string, objects []WITSMLComposite) error {
	var err error
	var data []byte
	dataSets := BuildDataSetsDPR_10(objects)
	if data, err = common.DatasetsToJson(dataSets); err != nil {
		return err
	}
	//need to write it to a file
	return common.Write2File(path, data)
}

func extractFacilityInfo(dataSetName string, groupedFacilities map[string][]Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"FacilityKind", "FacilityName", "UidRef", "NameKind", "FileName"}
	for kind, facilities := range groupedFacilities {
		for i := 0; i < len(facilities); i++ {
			row := common.RowData{}
			//add the top identification information so that it can be bound back to the original file
			row.AddStrValue(kind)
			row.AddStrValue(facilities[i].Name.Name)
			row.AddStrValue(facilities[i].Name.UidRef)
			row.AddStrValue(facilities[i].Name.Kind)
			row.AddStrValue(facilities[i].DataIdentification.DataIdentification.FileName)
			rows = append(rows, row)
		}
	}
	dataSet.Rows = rows
	return dataSet
}

func extractFlowInformation(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "FileName", "FilePath",
		"FacilityName", "FacilityKind", "FlowName", "FlowKind", "FlowQualifier"}
	for i := 0; i < len(objects); i++ {
		for x := 0; x < len(objects[i].ProdVolumeSet.ProdVolumes); x++ {
			for s := 0; s < len(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities); s++ {
				for y := 0; y < len(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Flow); y++ {
					row := common.RowData{}
					//add the top identification information so that it can be bound back to the original file
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].DataIdentification.DataIdentification.UUid)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].DataIdentification.DataIdentification.FileName)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].DataIdentification.DataIdentification.FilePath)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Name.Name)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Name.Kind)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Flow[y].Name)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Flow[y].Kind)
					row.AddStrValue(objects[i].ProdVolumeSet.ProdVolumes[x].Facilities[s].Flow[y].Qualifier)
					rows = append(rows, row)
				}
			}
		}
	}
	dataSet.Rows = rows
	return dataSet
}

func extractWellTests_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "FileName", "FilePath",
		"WellUid", "WellName", "WellboreName", "WellTestName",
		"TestDate", "TestType", "ChokeSize", "ChokeSizeUoM",
		"StandardTempPres", "TestDuration", "TestDurationUoM", "WHT", "WHTUoM",
		"WHP", "WHPUoM", "ChokeSize", "ChokeSizeUoM", "SepPress", "SepPressUoM",
		"SepTemp", "SepTempUoM", "OilRate", "OilRateUoM", "GasRate", "GasRateUoM",
		"WatRate", "WatRateUoM", "GoR", "GoRUoM"}
	for i := 0; i < len(objects); i++ {
		dIdentification := objects[i].DataIdentification
		wSet := objects[i].Wellset
		for s := 0; s < len(wSet.Wells); s++ {
			for y := 0; y < len(wSet.Wells[s].Wellbores); y++ {
				for x := 0; x < len(wSet.Wells[s].Wellbores[y].WellTests); x++ {

					row := common.RowData{}
					//add the top identification information so that it can be bound back to the original file
					row.AddStrValue(dIdentification.UUid)
					row.AddStrValue(dIdentification.FileName)
					row.AddStrValue(dIdentification.FilePath)
					row.AddStrValue(wSet.Wells[s].Uid)
					row.AddStrValue(wSet.Wells[s].Name)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].Name)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].Name)
					row.AddTimeValue(wSet.Wells[s].Wellbores[y].WellTests[x].TestDate.Time)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].TestType)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ChokeSize.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ChokeSize.Uom)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].StandardTempPres)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.TestDuration.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.TestDuration.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WHT.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WHT.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WHP.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WHP.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.ChokeSize.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.ChokeSize.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.SepPress.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.SepPress.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.SepTemp.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.SepTemp.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.OilRate.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.OilRate.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.GasRate.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.GasRate.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WatRate.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.WatRate.Uom)
					row.AddFloatValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.GoR.Value)
					row.AddStrValue(wSet.Wells[s].Wellbores[y].WellTests[x].ProductionTest.GoR.Uom)
					rows = append(rows, row)

				}
			}
		}

	}

	dataSet.Rows = rows
	return dataSet
}

func extractFacilityProdVolumesBHPAndWHP(dataSetName string, facilities []Facility) map[string]*common.DataSet {
	returnData := make(map[string]*common.DataSet)
	//dataSet := common.DataSet{}
	//dataSet.Name = dataSetName
	zap.S().Infof("Building dataset with name:%s", dataSetName)
	headersName := []string{"DataUUID", "FileName", "FilePath", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"FacilityParent1", "FacilityParent1Kind", "FacilityParent1UID",
		"FacilityParent2", "FacilityParent2Kind", "FacilityParent2UID",
		"ContextFacility", "ContextFacilityKind", "ContextFacilityUID",
		"FlowUid", "InteralFlowUid", "FlowKind", "FlowQualifier", "FlowName",
		"ReportStart", "ReportEnd",
		"chokeRelative", "chokeUoM",
		"Pressure", "PressureUoM", "Temp", "TempUoM"}
	//operationTime, wellProducing, wellInjecting, chokeRelative reserved for well/wellbores

	for i := 0; i < len(facilities); i++ {

		for x := 0; x < len(facilities[i].Flow); x++ {
			row := common.RowData{}
			//add the top identification information so that it can be bound back to the original file
			row.AddStrValue(facilities[i].DataIdentification.DataIdentification.UUid)
			row.AddStrValue(facilities[i].DataIdentification.DataIdentification.FileName)
			row.AddStrValue(facilities[i].DataIdentification.DataIdentification.FilePath)
			row.AddStrValue(facilities[i].DataIdentification.Installation.Kind)
			row.AddStrValue(facilities[i].DataIdentification.Installation.UidRef)
			row.AddStrValue(facilities[i].DataIdentification.Installation.Name)
			//extract top facility data
			row.AddStrValue(facilities[i].Name.Kind)   //add facility kind
			row.AddStrValue(facilities[i].Name.UidRef) //add facility uid ref
			row.AddStrValue(facilities[i].Name.Name)   //add facility name
			//add facilityparent1
			row.AddStrValue(facilities[i].FacilityParent1.Name)
			row.AddStrValue(facilities[i].FacilityParent1.Kind)
			row.AddStrValue(facilities[i].FacilityParent1.UidRef)
			//add facilityparent2
			row.AddStrValue(facilities[i].FacilityParent2.Name)
			row.AddStrValue(facilities[i].FacilityParent2.Kind)
			row.AddStrValue(facilities[i].FacilityParent2.UidRef)
			//add contextFacility
			row.AddStrValue(facilities[i].ContextFacility.Name)
			row.AddStrValue(facilities[i].ContextFacility.Kind)
			row.AddStrValue(facilities[i].ContextFacility.UidRef)
			//extract flow data
			row.AddStrValue(facilities[i].Flow[x].UID)             //add flow uid
			row.AddStrValue(facilities[i].Flow[x].InternalFlowUid) //add the internal flow uid created during parsing
			row.AddStrValue(facilities[i].Flow[x].Kind)            //flow kind
			row.AddStrValue(facilities[i].Flow[x].Qualifier)       //add flow qualifier
			row.AddStrValue(facilities[i].Flow[x].Name)
			//add the dates
			row.AddTimeValue(facilities[i].DataIdentification.DataIdentification.ReportStart)
			row.AddTimeValue(facilities[i].DataIdentification.DataIdentification.ReportEnd)
			if facilities[i].Flow[x].PortDiff.ChokeRelative.Uom != "" {
				//add the choke
				row.AddFloatValue(facilities[i].Flow[x].PortDiff.ChokeRelative.Value)
				row.AddStrValue(facilities[i].Flow[x].PortDiff.ChokeRelative.Uom)
			} else {
				//add empty columns
				row.AddFloatValue(common.NullFloatValue)
				row.AddStrValue("")
			}
			if facilities[i].Flow[x].Pres.Uom != "" {
				//add wellhead pressure
				row.AddFloatValue(facilities[i].Flow[x].Pres.Value)
				row.AddStrValue(facilities[i].Flow[x].Pres.Uom)
			} else {
				//add empty columns
				row.AddFloatValue(common.NullFloatValue)
				row.AddStrValue("")
			}
			if facilities[i].Flow[x].Temp.Uom != "" {
				//add wellhead pressure
				row.AddFloatValue(facilities[i].Flow[x].Temp.Value)
				row.AddStrValue(facilities[i].Flow[x].Temp.Uom)
			} else {
				//add empty columns
				row.AddFloatValue(common.NullFloatValue)
				row.AddStrValue("")
			}

			key := strings.ToUpper(facilities[i].Name.Kind)
			if _, ok := returnData[key]; ok {
				//key is there just add data to that row set
				returnData[key].Rows = append(returnData[key].Rows, row)
			} else {
				//not there add it
				rows := []common.RowData{}
				rows = append(rows, row)
				returnData[key] = &common.DataSet{HeadersName: headersName, Name: key, Rows: rows}
			}
		}
	}
	return returnData
}

func extractFacilityProdVolumes(dataSetName string, facilities []Facility) map[string]*common.DataSet {
	//var rows []common.RowData
	returnData := make(map[string]*common.DataSet)
	//dataSet := common.DataSet{}
	//dataSet.Name = dataSetName
	zap.S().Infof("Building dataset with name:%s", dataSetName)
	headersName := []string{"DataUUID", "FileName", "FilePath", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"FacilityParent1", "FacilityParent1Kind", "FacilityParent1UID",
		"FacilityParent2", "FacilityParent2Kind", "FacilityParent2UID",
		"ContextFacility", "ContextFacilityKind", "ContextFacilityUID",
		"FlowUid", "InteralFlowUid", "FlowKind", "FlowQualifier", "FlowName", "ProductKind", "ProductName",
		"PeriodKind", "dateStart", "dateEnd", "dTimStart", "dTimEnd",
		"volume", "volumeUoM", "volumeStd", "volumeStdUoM", "volumeValue", "volumeValueUoM",
		"volumeValueTempCondValue", "volumeValueTempCondUom", "volumeValuePresCondValue",
		"volumeValuePresCondUoM", "densityValue", "densityValueUoM",
		"densityValueTempCondValue", "densityValueTempCondUom",
		"densityValuePresCondValue", "densityValuePresCondUoM",
		"densityStd", "densityStdUoM", "bsw", "bswUoM",
		"rvp", "rvpUoM",
		"mass", "massUoM",
		//"GoR", "GoRUoM", "WaterConcVol", "WaterConcVolUoM",
		"operationTime", "operationTimeUoM", "wellProducing", "wellInjecting",
		"chokeRelative", "chokeUoM",
		"WellheadPressure", "WellheadPressureUoM", "WellheadTemp", "WellheadTempUoM",
		"BottomholePressure", "BottomholePressureUoM", "BottomholeTemp", "BottomholeTempUoM"}
	//operationTime, wellProducing, wellInjecting, chokeRelative reserved for well/wellbores

	for i := 0; i < len(facilities); i++ {

		for x := 0; x < len(facilities[i].Flow); x++ {
			//need to take out all of the base things relating to flow

			for z := 0; z < len(facilities[i].Flow[x].Product); z++ {

				for s := 0; s < len(facilities[i].Flow[x].Product[z].Period); s++ {

					row := common.RowData{}
					//add the top identification information so that it can be bound back to the original file
					row.AddStrValue(facilities[i].DataIdentification.DataIdentification.UUid)
					row.AddStrValue(facilities[i].DataIdentification.DataIdentification.FileName)
					row.AddStrValue(facilities[i].DataIdentification.DataIdentification.FilePath)
					row.AddStrValue(facilities[i].DataIdentification.Installation.Kind)
					row.AddStrValue(facilities[i].DataIdentification.Installation.UidRef)
					row.AddStrValue(facilities[i].DataIdentification.Installation.Name)
					//extract top facility data
					row.AddStrValue(facilities[i].Name.Kind)   //add facility kind
					row.AddStrValue(facilities[i].Name.UidRef) //add facility uid ref
					row.AddStrValue(facilities[i].Name.Name)   //add facility name
					//add facilityparent1
					row.AddStrValue(facilities[i].FacilityParent1.Name)
					row.AddStrValue(facilities[i].FacilityParent1.Kind)
					row.AddStrValue(facilities[i].FacilityParent1.UidRef)
					//add facilityparent2
					row.AddStrValue(facilities[i].FacilityParent2.Name)
					row.AddStrValue(facilities[i].FacilityParent2.Kind)
					row.AddStrValue(facilities[i].FacilityParent2.UidRef)
					//add contextFacility
					row.AddStrValue(facilities[i].ContextFacility.Name)
					row.AddStrValue(facilities[i].ContextFacility.Kind)
					row.AddStrValue(facilities[i].ContextFacility.UidRef)
					//extract flow data
					row.AddStrValue(facilities[i].Flow[x].UID)             //add flow uid
					row.AddStrValue(facilities[i].Flow[x].InternalFlowUid) //add the internal flow uid created during parsing
					row.AddStrValue(facilities[i].Flow[x].Kind)            //flow kind
					row.AddStrValue(facilities[i].Flow[x].Qualifier)       //add flow qualifier
					row.AddStrValue(facilities[i].Flow[x].Name)            //add flow name
					//extract period and product data
					row.AddStrValue(facilities[i].Flow[x].Product[z].Kind)                      // add product kind
					row.AddStrValue(facilities[i].Flow[x].Product[z].Name)                      //add product kind
					row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Kind)            //add period kind
					row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateStart.Time) //add date start
					row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateEnd.Time)
					row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DTimStart.Time) //add dtim start and end
					row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DTimEnd.Time)   //add dtime end
					if facilities[i].Flow[x].Product[z].Period[s].VolumeOnly.Uom != "" {
						//just add it if there is data
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].VolumeOnly.Value) //add volume only value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].VolumeOnly.Uom)

					} else {
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
					}
					if facilities[i].Flow[x].Product[z].Period[s].VolumeStd.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].VolumeStd.Value) //add volume std value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].VolumeStd.Uom)
					} else {
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
					}

					if facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Value) //add volume value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Uom)
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Value)     //add volume value temp cond
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Uom)         //add volume value temp uom
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Value) //add volume value pres cond
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Uom)     //add volume value pres uom
					} else {
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("") //add volume value pres uom
					}
					if facilities[i].Flow[x].Product[z].Period[s].Density.Density.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Value)  //add density value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Uom)      //add density uom
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Value)     //add density temp cond value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Uom)         //add density temp cond uom
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Value) //add density pressure cond value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Uom)     //add density pressure cond uom

					} else {
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
					}
					//need to add density std element if evident
					if facilities[i].Flow[x].Product[z].DensityStd.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].DensityStd.Value)
						row.AddStrValue(facilities[i].Flow[x].Product[z].DensityStd.Uom)
					} else {
						//add null values
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")

					}
					//add bsw if evident
					if facilities[i].Flow[x].Product[z].BsW.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].BsW.Value)
						row.AddStrValue(facilities[i].Flow[x].Product[z].BsW.Uom)
					} else {
						//add null values
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")

					}
					//add rvp if evident
					if facilities[i].Flow[x].Product[z].Rvp.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Rvp.Value)
						row.AddStrValue(facilities[i].Flow[x].Product[z].Rvp.Uom)
					} else {
						//add null values
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")

					}
					if facilities[i].Flow[x].Product[z].Period[s].Mass.Uom != "" {
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Value) //add mass value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Uom)
					} else {
						row.AddFloatValue(common.NullFloatValue)
						row.AddStrValue("")
					}

					//need to process additional well/wellbore data if it is of that kind...
					if facilities[i].Name.Kind == "well" || facilities[i].Name.Kind == "wellbore" {
						row.AddFloatValue(facilities[i].OperationTime.Value)
						row.AddStrValue(facilities[i].OperationTime.Uom)
						row.AddStrValue(facilities[i].WellProducing)
						row.AddStrValue(facilities[i].WellInjecting)
						if facilities[i].Flow[x].PortDiff.ChokeRelative.Uom != "" {
							row.AddFloatValue(facilities[i].Flow[x].PortDiff.ChokeRelative.Value)
							row.AddStrValue(facilities[i].Flow[x].PortDiff.ChokeRelative.Uom)
						} else if facilities[i].Flow[x].PortDiff.ChokeRelative.Uom == "" {
							//just see if we can locate the choke in the flow
							chokeValue := getChokeData(facilities[i].Flow)
							if chokeValue.Uom != "" {
								row.AddFloatValue(chokeValue.Value)
								row.AddStrValue(chokeValue.Uom)
							}
						} else {
							row.AddFloatValue(common.NullFloatValue)
							row.AddStrValue("")
						}

					} else {
						//just add empty columns
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
					}
					if facilities[i].Name.Kind == "wellhead" {
						//need to add whp and wht
						if facilities[i].Flow[x].Pres.Uom != "" {
							row.AddFloatValue(facilities[i].Flow[x].Pres.Value)
							row.AddStrValue(facilities[i].Flow[x].Pres.Uom)
						} else {
							row.AddEmptyColumn()
							row.AddEmptyColumn()
						}
						if facilities[i].Flow[x].Temp.Uom != "" {
							row.AddFloatValue(facilities[i].Flow[x].Temp.Value)
							row.AddStrValue(facilities[i].Flow[x].Temp.Uom)
						} else {
							row.AddEmptyColumn()
							row.AddEmptyColumn()
						}

					} else {
						//just check if we can locate the flow using the uid
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
					}
					if facilities[i].Name.Kind == "bottomhole" {
						//need to add whp and wht
						if facilities[i].Flow[x].Pres.Uom != "" {
							row.AddFloatValue(facilities[i].Flow[x].Pres.Value)
							row.AddStrValue(facilities[i].Flow[x].Pres.Uom)
						} else {
							row.AddEmptyColumn()
							row.AddEmptyColumn()
						}
						if facilities[i].Flow[x].Temp.Uom != "" {
							row.AddFloatValue(facilities[i].Flow[x].Temp.Value)
							row.AddStrValue(facilities[i].Flow[x].Temp.Uom)
						} else {
							row.AddEmptyColumn()
							row.AddEmptyColumn()
						}

					} else {
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
						row.AddEmptyColumn()
					}

					//need to check if we have the combination of facility kind and period in the return map otherwise we just add it
					key := strings.ToUpper(facilities[i].Name.Kind + "_" + strings.Replace(facilities[i].Flow[x].Product[z].Period[s].Kind, " ", "_", -1))
					if _, ok := returnData[key]; ok {
						//key is there just add data to that row set
						returnData[key].Rows = append(returnData[key].Rows, row)
					} else {
						//not there add it
						rows := []common.RowData{}
						rows = append(rows, row)
						returnData[key] = &common.DataSet{HeadersName: headersName, Name: key, Rows: rows}
					}
					//rows = append(rows, row)

				}
			}
		}
	}
	//dataSet.Rows = rows
	return returnData
}

func getWellheadOrBottomHoleFlow(uid string, facilities []Facility, facilityKind string) Flow {
	returnFlow := Flow{}
	breakOut := false

	for i := 0; i < len(facilities) && !breakOut; i++ {
		zap.S().Infof("Checking facility:%s, uid:%s, kind:%s against uid:%s", facilities[i].Name.Name, facilities[i].Name.UidRef,
			facilities[i].Name.Kind, uid)
		if facilities[i].Name.UidRef == uid && facilities[i].Name.Kind == facilityKind {
			//get the flow that contains temp and pres
			for s := 0; s < len(facilities[i].Flow) && !breakOut; s++ {
				if facilities[i].Flow[s].Temp.Uom != "" || facilities[i].Flow[s].Pres.Uom != "" {
					zap.S().Infof("LOcated facilitykind:%s, with uidref:%s", facilityKind, uid)
					breakOut = true
					returnFlow = facilities[i].Flow[s]
					if facilities[i].Flow[s].Temp.Uom == "" {
						//just blank out the temp with default null value
						returnFlow.Temp.Value = common.NullFloatValue
					}
					if facilities[i].Flow[s].Pres.Uom == "" {
						//just blank out the temp with default null value
						returnFlow.Pres.Value = common.NullFloatValue
					}
				}
			}
		}
	}
	return returnFlow
}

func getChokeData(flows []Flow) Value {
	returnValue := Value{Uom: "", Value: common.NullFloatValue}
	breakOut := false
	for i := 0; i < len(flows) && !breakOut; i++ {
		if flows[i].PortDiff.ChokeRelative.Uom != "" {
			returnValue = Value{Uom: flows[i].PortDiff.ChokeRelative.Uom,
				Value: flows[i].PortDiff.ChokeRelative.Value}
			breakOut = true
		}
	}
	return returnValue
}

//Extracts data from the document info in the mprml struct objects
func extractDocumentInfo_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "FileName", "FilePath", "DocumentName", "DocumentDate",
		"AuditEventDate", "AuditEventRespParty", "AuditEventComment"}

	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DocumentInfo.DocumentName.Txt
		documentDate := objects[i].DocumentInfo.DocumentDate
		for x := 0; x < len(objects[i].DocumentInfo.AuditTrail.Events); x++ {
			eventDate := objects[i].DocumentInfo.AuditTrail.Events[x].EventDate
			eventRespParty := objects[i].DocumentInfo.AuditTrail.Events[x].ResponsibleParty
			eventComment := objects[i].DocumentInfo.AuditTrail.Events[x].Comment
			row := common.RowData{}
			row.AddStrValue(dataUUID)
			row.AddStrValue(objects[i].DataIdentification.FileName)
			row.AddStrValue(objects[i].DataIdentification.FilePath)
			row.AddStrValue(documentName)
			row.AddTimeValue(documentDate.Time)
			row.AddTimeValue(eventDate.Time)
			row.AddStrValue(eventRespParty)
			row.AddStrValue(eventComment)
			rows = append(rows, row)
		}

	}
	dataSet.Rows = rows
	return dataSet
}

func extractPersonell_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "FileName", "FilePath",
		"PeriodKind", "DateStart", "DateEnd", "DTimStart", "DTimEnd",
		"Installation_Name", "Installation_Kind",
		"BedsAvailable", "CrewType", "CrewCount", "Work", "WorkUoM"}
	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DataIdentification.DocumentName
		fileName := objects[i].DataIdentification.FileName
		filePath := objects[i].DataIdentification.FilePath

		for x := 0; x < len(objects[i].ProdOperationSet.ProdOperation); x++ {
			dateStart := objects[i].ProdOperationSet.ProdOperation[x].DateStart
			dateEnd := objects[i].ProdOperationSet.ProdOperation[x].DateEnd
			dtimStart := objects[i].ProdOperationSet.ProdOperation[x].DTimStart
			dtimEnd := objects[i].ProdOperationSet.ProdOperation[x].DTimEnd
			periodKind := objects[i].ProdOperationSet.ProdOperation[x].Kind
			for y := 0; y < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport); y++ {
				bedsAvailable := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].BedsAvailable
				work := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].Work.Value
				workUoM := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].Work.Uom
				for z := 0; z < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].CrewCounts); z++ {

					row := common.RowData{}
					row.AddStrValue(dataUUID)
					row.AddStrValue(documentName)
					row.AddStrValue(fileName)
					row.AddStrValue(filePath)
					row.AddStrValue(periodKind)
					row.AddTimeValue(dateStart.Time)
					row.AddTimeValue(dateEnd.Time)
					row.AddTimeValue(dtimStart.Time)
					row.AddTimeValue(dtimEnd.Time)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Name)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Kind)
					row.AddIntValue(bedsAvailable)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].CrewCounts[z].Type)
					row.AddIntValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].CrewCounts[z].Count)
					row.AddFloatValue(work.float64)
					row.AddStrValue(workUoM)
					rows = append(rows, row)

				}
				//
			}

		}
	}

	dataSet.Rows = rows
	return dataSet
}

func extract_OpComments_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "FileName", "FilePath",
		"PeriodKind", "DateStart", "DateEnd", "DTimStart", "DTimEnd",
		"Installation_Name", "Installation_Kind", "CommentType", "CommentDTimStart", "CommentDTimEnd", "Comment"}
	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DataIdentification.DocumentName
		fileName := objects[i].DataIdentification.FileName
		filePath := objects[i].DataIdentification.FilePath

		for x := 0; x < len(objects[i].ProdOperationSet.ProdOperation); x++ {
			dateStart := objects[i].ProdOperationSet.ProdOperation[x].DateStart
			dateEnd := objects[i].ProdOperationSet.ProdOperation[x].DateEnd
			dtimStart := objects[i].ProdOperationSet.ProdOperation[x].DTimStart
			dtimEnd := objects[i].ProdOperationSet.ProdOperation[x].DTimEnd
			periodKind := objects[i].ProdOperationSet.ProdOperation[x].Kind
			for y := 0; y < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport); y++ {
				for z := 0; z < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.OperationalComments); z++ {
					commentType := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.OperationalComments[z].Type
					commentDTimStart := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.OperationalComments[z].DTimStart
					commentDTimEnd := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.OperationalComments[z].DTimEnd
					comment := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.OperationalComments[z].Comment
					row := common.RowData{}
					row.AddStrValue(dataUUID)
					row.AddStrValue(documentName)
					row.AddStrValue(fileName)
					row.AddStrValue(filePath)
					row.AddStrValue(periodKind)
					row.AddTimeValue(dateStart.Time)
					row.AddTimeValue(dateEnd.Time)
					row.AddTimeValue(dtimStart.Time)
					row.AddTimeValue(dtimEnd.Time)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Name)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Kind)
					row.AddStrValue(commentType)
					row.AddTimeValue(commentDTimStart.Time)
					row.AddTimeValue(commentDTimEnd.Time)
					row.AddStrValue(comment)
					rows = append(rows, row)
				}
			}
		}
	}
	dataSet.Rows = rows
	return dataSet
}

func extractLostProduction_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "FileName", "FilePath",
		"PeriodKind", "DateStart", "DateEnd", "DTimStart", "DTimEnd",
		"Installation_Name", "Installation_Kind",
		"ReasonLost", "Volume", "VolumeUoM"}
	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DataIdentification.DocumentName
		fileName := objects[i].DataIdentification.FileName
		filePath := objects[i].DataIdentification.FilePath

		for x := 0; x < len(objects[i].ProdOperationSet.ProdOperation); x++ {
			dateStart := objects[i].ProdOperationSet.ProdOperation[x].DateStart
			dateEnd := objects[i].ProdOperationSet.ProdOperation[x].DateEnd
			dtimStart := objects[i].ProdOperationSet.ProdOperation[x].DTimStart
			dtimEnd := objects[i].ProdOperationSet.ProdOperation[x].DTimEnd
			periodKind := objects[i].ProdOperationSet.ProdOperation[x].Kind
			for y := 0; y < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport); y++ {
				for z := 0; z < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.LostProduction.Reasons); z++ {

					row := common.RowData{}
					row.AddStrValue(dataUUID)
					row.AddStrValue(documentName)
					row.AddStrValue(fileName)
					row.AddStrValue(filePath)
					row.AddStrValue(periodKind)
					row.AddTimeValue(dateStart.Time)
					row.AddTimeValue(dateEnd.Time)
					row.AddTimeValue(dtimStart.Time)
					row.AddTimeValue(dtimEnd.Time)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Name)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Kind)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.LostProduction.Reasons[z].ReasonLost)
					row.AddFloatValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.LostProduction.Reasons[z].Value)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.LostProduction.Reasons[z].UoM)
					rows = append(rows, row)
				}
			}
		}
	}
	dataSet.Rows = rows
	return dataSet
}

func extractWaterCleaning_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "FileName", "FilePath",
		"PeriodKind", "DateStart", "DateEnd", "DTimStart", "DTimEnd",
		"Installation_Name", "Installation_Kind",
		"WaterCleaningUid", "SamplePoint", "OilInWaterProduced", "OilInWaterProducedUoM"}
	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DataIdentification.DocumentName
		fileName := objects[i].DataIdentification.FileName
		filePath := objects[i].DataIdentification.FilePath

		for x := 0; x < len(objects[i].ProdOperationSet.ProdOperation); x++ {
			dateStart := objects[i].ProdOperationSet.ProdOperation[x].DateStart
			dateEnd := objects[i].ProdOperationSet.ProdOperation[x].DateEnd
			dtimStart := objects[i].ProdOperationSet.ProdOperation[x].DTimStart
			dtimEnd := objects[i].ProdOperationSet.ProdOperation[x].DTimEnd
			periodKind := objects[i].ProdOperationSet.ProdOperation[x].Kind
			for y := 0; y < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport); y++ {
				for z := 0; z < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.WaterCleaning); z++ {

					row := common.RowData{}
					row.AddStrValue(dataUUID)
					row.AddStrValue(documentName)
					row.AddStrValue(fileName)
					row.AddStrValue(filePath)
					row.AddStrValue(periodKind)
					row.AddTimeValue(dateStart.Time)
					row.AddTimeValue(dateEnd.Time)
					row.AddTimeValue(dtimStart.Time)
					row.AddTimeValue(dtimEnd.Time)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Name)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Kind)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.WaterCleaning[z].Uid)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.WaterCleaning[z].SamplePoint)
					row.AddFloatValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.WaterCleaning[z].OilInWaterProduced.Value)
					row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].ProductionActivity.WaterCleaning[z].OilInWaterProduced.Uom)
					rows = append(rows, row)
				}
			}
		}
	}
	dataSet.Rows = rows
	return dataSet
}

func extractOpHSE_DPR10(dataSetName string, objects []WITSMLComposite) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "FileName", "FilePath",
		"PeriodKind", "DateStart", "DateEnd", "DTimStart", "DTimEnd",
		"Installation_Name", "Installation_Kind",
		"IncidentCount", "SafetyCountType", "SafetyCountPeriodType", "SafetyCount",
		"SafetyIntroCount", "SinceLostTime", "SinceLostTimeUoM"}
	for i := 0; i < len(objects); i++ {
		dataUUID := objects[i].DataIdentification.UUid
		documentName := objects[i].DataIdentification.DocumentName
		fileName := objects[i].DataIdentification.FileName
		filePath := objects[i].DataIdentification.FilePath

		for x := 0; x < len(objects[i].ProdOperationSet.ProdOperation); x++ {
			dateStart := objects[i].ProdOperationSet.ProdOperation[x].DateStart
			dateEnd := objects[i].ProdOperationSet.ProdOperation[x].DateEnd
			dtimStart := objects[i].ProdOperationSet.ProdOperation[x].DTimStart
			dtimEnd := objects[i].ProdOperationSet.ProdOperation[x].DTimEnd
			periodKind := objects[i].ProdOperationSet.ProdOperation[x].Kind
			for y := 0; y < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport); y++ {
				for z := 0; z < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE); z++ {
					incidentCount := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].IncidentCount
					safetyIntroCount := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].SafetyIntroCount
					sinceLostTimeValue := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].SinceLostTime.Value
					sinceLostTimeUoM := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].SinceLostTime.Uom
					for w := 0; w < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety); w++ {
						//safetyCount:=objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety[w].SafetyCount
						for r := 0; r < len(objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety[w].SafetyCount); r++ {
							safetyCount := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety[w].SafetyCount[r].Count
							safetyPeriod := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety[w].SafetyCount[r].Period
							safetyType := objects[i].ProdOperationSet.ProdOperation[x].InstallationReport[y].OperationalHSE[z].Safety[w].SafetyCount[r].Type
							row := common.RowData{}
							row.AddStrValue(dataUUID)
							row.AddStrValue(documentName)
							row.AddStrValue(fileName)
							row.AddStrValue(filePath)
							row.AddStrValue(periodKind)
							row.AddTimeValue(dateStart.Time)
							row.AddTimeValue(dateEnd.Time)
							row.AddTimeValue(dtimStart.Time)
							row.AddTimeValue(dtimEnd.Time)
							row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Name)
							row.AddStrValue(objects[i].ProdOperationSet.ProdOperation[x].Installation.Kind)
							row.AddIntValue(incidentCount)
							row.AddStrValue(safetyType)
							row.AddStrValue(safetyPeriod)
							row.AddIntValue(safetyCount)
							row.AddIntValue(safetyIntroCount)
							row.AddFloatValue(sinceLostTimeValue)
							row.AddStrValue(sinceLostTimeUoM)
							rows = append(rows, row)
						}
					}
				}
				//
			}

		}
	}

	dataSet.Rows = rows
	return dataSet
}

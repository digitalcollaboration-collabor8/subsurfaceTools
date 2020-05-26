/*
* @Author: magsv
* @Date:   2019-02-25 08:36:59
* @Last Modified by:   magsv
* @Last Modified time: 2019-02-26 11:01:45
 */
package dpr20

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"go.uber.org/zap"
)

type ProdObjSliceType []Objects

//Parses a MPRML xml file to a struct for further processing
func ParseProdXMLFile(inputFile string) (Objects, error) {
	var objects Objects
	var err error
	var data []byte
	if data, err = common.ReadFile(inputFile); err != nil {
		return objects, err
	}
	if err = xml.Unmarshal(data, &objects); err != nil {
		return objects, err
	}
	//just add the uid and file data for this set of objects
	objects.DataIdentification = createDataId(inputFile)
	return objects, nil
}

//Converts a parsed mprml objects struc to JSON
func (data *ProdObjSliceType) ToJson() ([]byte, error) {

	return json.Marshal(data)
}

//Parses through a set of mprml objects and organises them by facility kind returning
//a map of facilities where the key is the facility kind
//the function makes sure that only production objects are included and not operational
func OrganiseProdFacilitiesByKind(objects []Objects) map[string][]Facility {
	facilities := make(map[string][]Facility)
	for a := 0; a < len(objects); a++ {
		for i := 0; i < len(objects[a].ProdObjects); i++ {
			if strings.ToLower(objects[a].ProdObjects[i].ObjType) == "obj_productvolume" {
				for x := 0; x < len(objects[a].ProdObjects[i].Facilities); x++ {
					//check if the facility type exists in the map..
					facilityKind := objects[a].ProdObjects[i].Facilities[x].Name.Kind
					objects[a].ProdObjects[i].Facilities[x].DataIdentification.DataIdentification = objects[a].DataIdentification
					objects[a].ProdObjects[i].Facilities[x].DataIdentification.Installation = objects[a].ProdObjects[i].Installation
					if len(facilities[facilityKind]) == 0 {
						//not there add it
						facilities[facilityKind] = []Facility{objects[a].ProdObjects[i].Facilities[x]}
					} else {
						//it is there we just append it to the existing list
						facilities[facilityKind] = append(facilities[facilityKind], objects[a].ProdObjects[i].Facilities[x])
					}
				}
			}
		}
	}
	return facilities
}

//Builds an excel output file for a given list of mprml objects
func BuildXLSFileForProduction(path string, objects []Objects, oneFilePerSheet bool, appendTimeInName bool) error {
	dataSets := BuildDataSets(objects)
	return common.CreateWorkbookFromDataSet(path, dataSets, oneFilePerSheet, appendTimeInName)
}

func BuildCSVFileForProduction(path string, objects []Objects, discriminator string) error {
	dataSets := BuildDataSets(objects)
	return common.DatasetsToCsv(dataSets, path, discriminator)
}

func BuildJsonFileForProduction(path string, objects []Objects) error {
	var err error
	var data []byte
	dataSets := BuildDataSets(objects)
	if data, err = common.DatasetsToJson(dataSets); err != nil {
		return err
	}
	//need to write it to a file
	return common.Write2File(path, data)
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
func ReadProdXMLFiles2Struct(folderPath string) ([]Objects, error) {
	var err error
	var objects []Objects
	var files []string
	var data Objects
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

//Parses through a list of mprml struct objects and generates a flattened
//version of the data with data organized into datasets with a new, e.g.
//the DOCUMENT_INFO dataset holds all document info objects from the parsed xml files
//All data is bound together with the UUID as found in the REPORT_FILE_INFO dataset that
//gives every parsed xml file a UUID to which data is associated
func BuildDataSets(objects []Objects) []common.DataSet {
	var datasets []common.DataSet
	var groupedProdFacilities map[string][]Facility
	groupedProdFacilities = OrganiseProdFacilitiesByKind(objects)
	//build the report file info
	datasets = append(datasets, extractFileInformation("REPORT_FILE_INFO", objects))
	datasets = append(datasets, extractFacilityInfo("FACILITIES", groupedProdFacilities))
	//add the document info data
	datasets = append(datasets, extractDocumentInfo("DOCUMENT_INFO", objects))
	//add the reportcontext
	rContext := extractReportContext("REPORT_CONTEXT", objects)
	datasets = append(datasets, rContext)
	//add the actual prod volume data
	for kind, facilities := range groupedProdFacilities {
		zap.S().Infof("Building sheet for facility kind:%s", kind)
		datasets = append(datasets, extractFacilityProdData(strings.ToUpper(kind), facilities))

	}
	for kind, facilities := range groupedProdFacilities {
		zap.S().Infof("Building sheet for facility kind:%s", kind)
		datasets = append(datasets, extractFacilityParamsData(strings.ToUpper(kind)+"_PARAMS", facilities))

	}

	//process the operational objects
	opObjects := extractOperationalObjects(objects)
	datasets = append(datasets, extract_OpComments_DPR20("OP_COMMENTS", opObjects))
	//add the reportcontext to all instances
	addReportContextToAllInstances(rContext, datasets)
	return datasets

}

func extractFacilityInfo(dataSetName string, groupedFacilities map[string][]Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"FacilityKind", "FacilityName", "UidRef", "NameKind", "FileName", "Installation"}
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

func addReportContextToAllInstances(rContext common.DataSet, dataSets []common.DataSet) {
	for i := 0; i < len(dataSets); i++ {
		if dataSets[i].Name != "REPORT_CONTEXT" {
			dataSets[i].HeadersName = append(dataSets[i].HeadersName, rContext.HeadersName...)
			for x := 0; x < len(dataSets[i].Rows); x++ {
				dataSets[i].Rows[x].Columns = append(dataSets[i].Rows[x].Columns, rContext.Rows[0].Columns...)
			}
		}
	}
}

func extractOperationalObjects(objects []Objects) []Objects {
	var retObjects []Objects
	var opObjects []ProdVolume
	retObject := Objects{}
	retObject.DataIdentification = objects[0].DataIdentification
	retObject.Context = objects[0].Context
	for i := 0; i < len(objects); i++ {
		for x := 0; x < len(objects[i].ProdObjects); x++ {
			if strings.ToLower(objects[i].ProdObjects[x].ObjType) == "obj_productionoperation" {
				opObjects = append(opObjects, objects[i].ProdObjects[x])
			}
		}
	}
	retObject.ProdObjects = opObjects
	retObjects = append(retObjects, retObject)
	return retObjects
}

func extract_OpComments_DPR20(dataSetName string, objects []Objects) common.DataSet {
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

		for x := 0; x < len(objects[i].ProdObjects); x++ {

			dateStart := objects[i].ProdObjects[x].StartDate
			dateEnd := objects[i].ProdObjects[x].EndDate
			dtimStart := objects[i].ProdObjects[x].StartTime
			dtimEnd := objects[i].ProdObjects[x].EndTime
			periodKind := objects[i].ProdObjects[x].PeriodKind

			for y := 0; y < len(objects[i].ProdObjects[x].InstallationReport); y++ {
				for z := 0; z < len(objects[i].ProdObjects[x].InstallationReport[y].ProductionActivity.OperationalComments); z++ {
					commentType := objects[i].ProdObjects[x].InstallationReport[y].ProductionActivity.OperationalComments[z].Type
					commentDTimStart := objects[i].ProdObjects[x].InstallationReport[y].ProductionActivity.OperationalComments[z].DTimStart
					commentDTimEnd := objects[i].ProdObjects[x].InstallationReport[y].ProductionActivity.OperationalComments[z].DTimEnd
					comment := objects[i].ProdObjects[x].InstallationReport[y].ProductionActivity.OperationalComments[z].Comment
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
					row.AddStrValue(objects[i].ProdObjects[x].Installation.Name)
					row.AddStrValue(objects[i].ProdObjects[x].Installation.Kind)
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

//Extracts data from the report context in the mprml struct objects
func extractReportContext(dataSetName string, objects []Objects) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "ReportKind",
		"ReportTitle", "ReportMonth", "ReportYear",
		"ReportVersion", "ReportStatus", "ReportInstallationKind",
		"ReportInstallationUidRef", "ReportInstallationName", "ReportStartDate", "ReportEndDate"}
	for i := 0; i < len(objects); i++ {

		dataUUID := objects[i].DataIdentification.UUid
		reportKind := objects[i].Context.Kind
		reportTitle := objects[i].Context.Title.Txt
		reportMonth := objects[i].Context.Month
		reportYear := objects[i].Context.Year
		reportStartDate := objects[i].Context.StartDate
		reportEndDate := objects[i].Context.EndDate
		reportVersion := objects[i].Context.ReportVersion
		reportStatus := objects[i].Context.ReportStatus
		reportInstallationKind := objects[i].Context.Installation.Kind
		reportInstallationUidRef := objects[i].Context.Installation.UidRef
		reportInstallationName := objects[i].Context.Installation.Name
		row := common.RowData{}
		row.AddStrValue(dataUUID)
		row.AddStrValue(reportKind)
		row.AddStrValue(reportTitle)
		row.AddStrValue(reportMonth)
		row.AddStrValue(reportYear)
		row.AddFloatValue(reportVersion)
		row.AddStrValue(reportStatus)
		row.AddStrValue(reportInstallationKind)
		row.AddStrValue(reportInstallationUidRef)
		row.AddStrValue(reportInstallationName)
		row.AddTimeValue(reportStartDate.Time)
		row.AddTimeValue(reportEndDate.Time)
		rows = append(rows, row)

	}
	dataSet.Rows = rows
	return dataSet
}

//Extracts data from the document info in the mprml struct objects
func extractDocumentInfo(dataSetName string, objects []Objects) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "DocumentName", "DocumentDate",
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

//Extracts data from the report files as read from the disk
func extractFileInformation(dataSetName string, objects []Objects) common.DataSet {
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
		row.AddStrValue(objects[i].DataIdentification.DocumentName)
		rows = append(rows, row)

	}
	dataSet.Rows = rows
	return dataSet

}

//Extracts data from the prod volume facility setup in the MPRML files, note that this
//function will extract entries that includes a balanceset as this will be processed by the cargo/installation processing
func extractFacilityProdData(dataSetName string, facilities []Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"FlowUid", "FlowKind", "FlowQualifier", "FlowName", "ProductKind", "ProductName", "PeriodKind", "dateStart",
		"dateEnd", "Volume", "VolumeUoM", "VolumeStd", "VolumeStdUoM", "VolumeValue", "VolumeValueUoM", "VolumeValueConditions",
		"Mass", "MassUoM", "Density", "DensityUoM", "DensityConditions"}
	if strings.ToLower(dataSetName) == "wellbore" || strings.ToLower(dataSetName) == "well" {
		dataSet.HeadersName = append(dataSet.HeadersName, []string{"OperationTime", "OperationTimeUoM"}...)

	}
	for i := 0; i < len(facilities); i++ {

		for x := 0; x < len(facilities[i].Flow); x++ {

			for z := 0; z < len(facilities[i].Flow[x].Product); z++ {

				for s := 0; s < len(facilities[i].Flow[x].Product[z].Period); s++ {
					//just check that we do not include balance breakdowns..these should be included on separate sheets
					if len(facilities[i].Flow[x].Product[z].Period[s].BalanceSets) == 0 {

						row := common.RowData{}
						//add the top identification information so that it can be bound back to the original file
						row.AddStrValue(facilities[i].DataIdentification.DataIdentification.UUid)
						row.AddStrValue(facilities[i].DataIdentification.Installation.Kind)
						row.AddStrValue(facilities[i].DataIdentification.Installation.UidRef)
						row.AddStrValue(facilities[i].DataIdentification.Installation.Name)
						//extract top facility data
						row.AddStrValue(facilities[i].Name.Kind)   //add facility kind
						row.AddStrValue(facilities[i].Name.UidRef) //add facility uid ref
						row.AddStrValue(facilities[i].Name.Name)   //add facility name
						//extract flow data
						row.AddStrValue(facilities[i].Flow[x].UID)       //add flow uid
						row.AddStrValue(facilities[i].Flow[x].Kind)      //flow kind
						row.AddStrValue(facilities[i].Flow[x].Qualifier) //add flow qualifier
						row.AddStrValue(facilities[i].Flow[x].Name)      //add flow name
						//extract period and product data
						row.AddStrValue(facilities[i].Flow[x].Product[z].Kind)                         // add product kind
						row.AddStrValue(facilities[i].Flow[x].Product[z].Name)                         //add product kind
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Kind)               //add period kind
						row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateStart.Time)    //add date start
						row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateEnd.Time)      //add date end
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].VolumeOnly.Value) //add volume only value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].VolumeOnly.Uom)     //add volume uom
						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].VolumeStd.Value)  //add volumestd value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].VolumeStd.Uom)      //add volumestd uom

						row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Value) //add volume value
						row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Uom)     //add volume uom
						row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Value, 'f', 0, 64) +
							facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Uom + "/" +
							strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Value, 'f', 0, 64) +
							facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Uom) //add volume cond
						//just exlude data here if it is a wellbore
						if facilities[i].Name.Kind != "wellbore" && facilities[i].Name.Kind != "well" {
							//just include non empty mass elements
							if facilities[i].Flow[x].Product[z].Period[s].Mass.Uom != "" {
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Value) //add mass
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Uom)     //add mass uom
							} else {
								row.AddEmptyColumn()
								row.AddEmptyColumn()
							}
							//just include non empty density elements
							if facilities[i].Flow[x].Product[z].Period[s].Density.Density.Uom != "" {
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Value.float64) //add density value
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Uom)             //add density uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Uom + "/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Uom) //add density cond
							} else {
								row.AddEmptyColumn()
								row.AddEmptyColumn()
								row.AddEmptyColumn()
							}
						} else {
							row.AddEmptyColumn()
							row.AddEmptyColumn()
							row.AddEmptyColumn()
							row.AddEmptyColumn()
							row.AddEmptyColumn()
							//add the operation time
							row.AddFloatValue(facilities[i].OperationTime.Value.float64)
							row.AddStrValue(facilities[i].OperationTime.Uom)
						}

						rows = append(rows, row)
					}
				}

			}
		}
	}
	dataSet.Rows = rows
	return dataSet
}

//extractFacilityParams extracts facility parameters
func extractFacilityParamsData(dataSetName string, facilities []Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"ParameterName", "DTimStart", "DTimEnd", "MeasureValue", "MeasureValueUoM"}
	for i := 0; i < len(facilities); i++ {
		for x := 0; x < len(facilities[i].ParameterSets); x++ {
			for y := 0; y < len(facilities[i].ParameterSets[x].Parameters); y++ {
				row := common.RowData{}
				//add the top identification information so that it can be bound back to the original file
				row.AddStrValue(facilities[i].DataIdentification.DataIdentification.UUid)
				row.AddStrValue(facilities[i].DataIdentification.Installation.Kind)
				row.AddStrValue(facilities[i].DataIdentification.Installation.UidRef)
				row.AddStrValue(facilities[i].DataIdentification.Installation.Name)
				//extract top facility data
				row.AddStrValue(facilities[i].Name.Kind)   //add facility kind
				row.AddStrValue(facilities[i].Name.UidRef) //add facility uid ref
				row.AddStrValue(facilities[i].Name.Name)   //add facility name
				//extract flow data
				row.AddStrValue(facilities[i].ParameterSets[x].Name)
				row.AddTimeValue(facilities[i].ParameterSets[x].Parameters[y].StartDate.Time)
				row.AddTimeValue(facilities[i].ParameterSets[x].Parameters[y].EndDate.Time)
				row.AddFloatValue(facilities[i].ParameterSets[x].Parameters[y].MeasureValue.Value)
				row.AddStrValue(facilities[i].ParameterSets[x].Parameters[y].MeasureValue.Uom)
				rows = append(rows, row)
			}
		}
	}

	dataSet.Rows = rows
	return dataSet
}

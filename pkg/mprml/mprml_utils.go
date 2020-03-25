/*
* @Author: magsv
* @Date:   2018-02-22 11:42:38
* @Last Modified by:   magsv
* @Last Modified time: 2018-04-11 12:19:11
 */
package mprml

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"go.uber.org/zap"
)

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
func ProdObjectsToJson(objects []Objects) ([]byte, error) {

	return json.MarshalIndent(objects, "", "    ")
}

//Parses through a set of mprml objects and organises them by facility kind returning
//a map of facilities where the key is the facility kind
func OrganiseFacilitiesByKind(objects []Objects) map[string][]Facility {
	facilities := make(map[string][]Facility)
	for a := 0; a < len(objects); a++ {
		for i := 0; i < len(objects[a].ProdObjects); i++ {
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
	var groupedFacilities map[string][]Facility
	var flatFacilities []Facility
	groupedFacilities = OrganiseFacilitiesByKind(objects)
	//build the report file info
	datasets = append(datasets, extractFileInformation("REPORT_FILE_INFO", objects))
	//add the document info data
	datasets = append(datasets, extractDocumentInfo("DOCUMENT_INFO", objects))
	//add the reportcontext
	datasets = append(datasets, extractReportContext("REPORT_CONTEXT", objects))
	//add the actual prod volume data
	for kind, facilities := range groupedFacilities {
		zap.S().Infof("Building sheet for facility kind:%s", kind)
		datasets = append(datasets, extractFacilityProdData(strings.ToUpper(kind), facilities))

	}
	//add the cargo sheet
	for _, facility := range groupedFacilities {
		flatFacilities = append(flatFacilities, facility...)
	}
	datasets = append(datasets, extractCargoData("CARGO", flatFacilities))
	datasets = append(datasets, extractInventoryData("INVENTORY", flatFacilities))
	return datasets

}

//Extracts data from the report context in the mprml struct objects
func extractReportContext(dataSetName string, objects []Objects) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = dataSetName
	dataSet.HeadersName = []string{"DataUUID", "ReportKind",
		"ReportTitle", "ReportMonth", "ReportYear",
		"ReportVersion", "ReportStatus", "ReportInstallationKind",
		"ReportInstallationUidRef", "ReportInstallationName"}
	for i := 0; i < len(objects); i++ {

		dataUUID := objects[i].DataIdentification.UUid
		reportKind := objects[i].Context.Kind
		reportTitle := objects[i].Context.Title.Txt
		reportMonth := objects[i].Context.Month
		reportYear := objects[i].Context.Year
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
	if strings.ToLower(dataSetName) == "wellbore" {
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
						if facilities[i].Name.Kind != "wellbore" {
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

//Extracts data from the report files under the product volume facility with a flow
//kind of hydrocarbon accounting and that has associated balancesets. This should
//be equal to just extracting the cargo information in the given set of files
func extractCargoData(datasetName string, facilities []Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = datasetName
	dataSet.HeadersName = []string{"DataUUID", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"FlowUid", "FlowKind", "FlowQualifier", "FlowName", "ProductKind",
		"ProductName", "PeriodKind", "dateStart",
		"dateEnd", "BalanceSetKind", "CargoNumber", "Destination", "Country",
		"EventKind", "EventDate",
		"TotalVolume", "TotalVolumeUoM", "TotalVolumeConditions",
		"TotalMass", "TotalMassUoM",
		"TotalDensity", "TotalDensityUoM", "TotalDensityConditions",
		"Owner", "Share", "ShareUoM", "OwnerVolume", "OwnerVolumeUoM", "OwnerVolumeConditions",
		"OwnerMass", "OwnerMassUoM"}
	for i := 0; i < len(facilities); i++ {

		for x := 0; x < len(facilities[i].Flow); x++ {

			for z := 0; z < len(facilities[i].Flow[x].Product); z++ {

				for s := 0; s < len(facilities[i].Flow[x].Product[z].Period); s++ {

					for w := 0; w < len(facilities[i].Flow[x].Product[z].Period[s].BalanceSets); w++ {

						for y := 0; y < len(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails); y++ {
							//just add data for cargo and nothing else
							if facilities[i].Flow[x].Kind == "hydrocarbon accounting" {
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
								row.AddStrValue(facilities[i].Flow[x].Product[z].Kind)                      // add product kind
								row.AddStrValue(facilities[i].Flow[x].Product[z].Name)                      //add product kind
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Kind)            //add period kind
								row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateStart.Time) //add date start
								row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateEnd.Time)   //add date end
								//process the balancesets
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Kind)                //add balanceset kind
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].CargoNumber)         //add cargo number
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Destination.Name)    //destination
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Destination.Country) //add country
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Event.Kind)          //add event kind
								row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Event.Date.Time)    //add event date

								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Volume.Value) //add total volume
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Volume.Uom)     //add volume uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Temp.Value, 'f', 0, 64) +

									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Volume.Pressure.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Mass.Value)
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Mass.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Density.Value.float64)
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Density.Uom)
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Density.Pressure.Uom)

								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Owner)                 //add the owner
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Share.Value)         //add owner share
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Share.Uom)             //add share uom
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Volume.Value) //add share volume
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Volume.Uom)     //add share volume uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Pressure.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Mass.Value) //add share mass
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Mass.Uom)     //add share mass
								rows = append(rows, row)
							}
						}
					}

				}
			}

		}
	}
	dataSet.Rows = rows
	return dataSet
}

//Extracts data from the report files from the prodVolume section with a flow kind type of inventory
//and where you have an associated balance set. This should be equal to extrating the closing inventory
//breakdown as associated with each owner of the stock
func extractInventoryData(datasetName string, facilities []Facility) common.DataSet {
	var rows []common.RowData
	dataSet := common.DataSet{}
	dataSet.Name = datasetName
	dataSet.HeadersName = []string{"DataUUID", "OwningInstallationKind",
		"OwningInstallationUid", "OwningInstallationName",
		"FacilityKind", "FacilityUID", "FacilityName",
		"FlowUid", "FlowKind", "FlowQualifier", "FlowName", "ProductKind",
		"ProductName", "PeriodKind", "dateStart",
		"dateEnd", "BalanceSetKind",
		"TotalVolume", "TotalVolumeUoM", "TotalVolumeConditions",
		"TotalMass", "TotalMassUoM",
		"TotalDensity", "TotalDensityUoM", "TotalDensityConditions",
		"Owner", "Share", "ShareUoM", "OwnerVolume", "OwnerVolumeUoM", "OwnerVolumeConditions",
		"OwnerMass", "OwnerMassUoM",
		"OwnerDensity", "OwnerDensityUoM", "OwnerDensityConditions"}
	for i := 0; i < len(facilities); i++ {

		for x := 0; x < len(facilities[i].Flow); x++ {

			for z := 0; z < len(facilities[i].Flow[x].Product); z++ {

				for s := 0; s < len(facilities[i].Flow[x].Product[z].Period); s++ {

					for w := 0; w < len(facilities[i].Flow[x].Product[z].Period[s].BalanceSets); w++ {

						for y := 0; y < len(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails); y++ {
							//just add data for cargo and nothing else
							if facilities[i].Flow[x].Kind == "inventory" {
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
								row.AddStrValue(facilities[i].Flow[x].Product[z].Kind)                          // add product kind
								row.AddStrValue(facilities[i].Flow[x].Product[z].Name)                          //add product kind
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Kind)                //add period kind
								row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateStart.Time)     //add date start
								row.AddTimeValue(facilities[i].Flow[x].Product[z].Period[s].DateEnd.Time)       //add date end
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].Kind) //add the balanceset kind
								//add the total volume, mass and density
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Value) //add total volume
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Volume.Volume.Uom)     //add volume uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Value, 'f', 0, 64) +

									facilities[i].Flow[x].Product[z].Period[s].Volume.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].Volume.Pressure.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Value)
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Mass.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Value.float64)
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].Density.Density.Uom)
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].Density.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].Density.Pressure.Uom)
								//process the owner splits
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Owner)                 //add the owner
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Share.Value)         //add owner share
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Share.Uom)             //add share uom
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Volume.Value) //add share volume
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Volume.Uom)     //add share volume uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Volume.Pressure.Uom)
								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Mass.Value) //add share mass
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Mass.Uom)     //add share mass

								row.AddFloatValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Density.Value.float64) //add share density
								row.AddStrValue(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Density.Uom)             //add share density uom
								row.AddStrValue(strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Temp.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Temp.Uom +
									"/" +
									strconv.FormatFloat(facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Pressure.Value, 'f', 0, 64) +
									facilities[i].Flow[x].Product[z].Period[s].BalanceSets[w].BalanceDetails[y].Density.Pressure.Uom)
								rows = append(rows, row)
							}
						}
					}

				}
			}

		}
	}
	dataSet.Rows = rows
	return dataSet
}

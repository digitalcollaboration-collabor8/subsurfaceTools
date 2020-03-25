/*
* @Author: magsv
* @Date:   2017-08-18 14:40:52
* @Last Modified by:   magsv
* @Last Modified time: 2018-06-25 10:10:03
 */
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/ddrml"
	"go.uber.org/zap"
)

var (
	Version string
	Build   string
)

func processDDRFiles(inputFolder string, outputFile string, logFile string, moveFolder string,
	outputFormat string, appendTime2Filename bool, oneFilePerSheet bool) {
	var err error
	var drillReports []ddrml.DrillReports
	var files []string
	start := time.Now()
	unmarshalStart := time.Now()
	//just create the folders needed if not existing
	if err = createNeededFolders(outputFile, logFile, moveFolder); err != nil {
		fmt.Printf("Failed in created needed folder:%s", err.Error())
		return
	}
	//initialize the log file

	cfg := zap.Config{

		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{logFile},
		ErrorOutputPaths: []string{logFile},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	if appendTime2Filename {
		//build the new filename
		outputFile = common.AppendTimeAndDateToFile(outputFile)
		zap.S().Infof("Append time 2 filename is true, generated new outputfile name:%s", outputFile)
	}
	if drillReports, err = ddrml.ReadDDRXMLFiles2Struct(inputFolder); err != nil {
		zap.S().Errorf("Failed in processing ddr xml files in folder:%s,error:%s", inputFolder, err.Error())
		return
	}
	unmarshalTook := time.Since(unmarshalStart)
	switch outputFormat {
	case "excel":
		err = outputDDRExcelData(outputFile, drillReports, oneFilePerSheet, appendTime2Filename)
	case "json":
		err = outputDDRJSONData(outputFile, drillReports)
	case "csv":
		err = outputDDRCsvData(outputFile, drillReports)
	default:
		err = outputDDRExcelData(outputFile, drillReports, oneFilePerSheet, appendTime2Filename)
	}
	if err != nil {
		zap.S().Errorf("Failed in outputting DDR data to format:%s, error:%s", outputFormat, err.Error())
		return
	}

	totalTimeTook := time.Since(start)
	if moveFolder != "" {
		zap.S().Infof("Moving DDR files in folder:%s to %s", inputFolder, moveFolder)
		//Need to move the data to the move folder
		//make sure the destination folder exists
		if err = common.CreateAllFolders(moveFolder); err != nil {
			zap.S().Errorf("Failed in creating DDR move folder:%s", err.Error())
			return
		}
		//get the list of the files
		if files, err = common.GetFilesWithExtension(inputFolder, "*.xml"); err != nil {
			zap.S().Error("Failed in reading DDR file in folder:", err.Error())
			return
		}
		if err = common.MoveFiles(files, moveFolder); err != nil {
			zap.S().Errorf("Failed in moving DDR files to folder:%s, error:%s", moveFolder, err.Error())
			return
		}

	}
	zap.S().Info("Finished DDR processing results to:", outputFile)
	zap.S().Infof("XML DDR File parsing took:%s", unmarshalTook)
	zap.S().Infof("Total time used:%s", totalTimeTook)

}

func outputDDRCsvData(outputFile string, drillReports []ddrml.DrillReports) error {

	return ddrml.BuildCsvFileForDrilling(outputFile, drillReports)
}
func outputDDRExcelData(outputFile string, objects []ddrml.DrillReports,
	oneFilePerSheet, appendTime2Filename bool) error {
	var err error
	xlsStart := time.Now()
	if err = ddrml.BuildXLSFileForDrilling(outputFile, objects, oneFilePerSheet, appendTime2Filename); err != nil {
		zap.S().Error("Error:", err.Error())
		return err
	}
	xlsBuildTook := time.Since(xlsStart)
	zap.S().Infof("Excel build took:%s", xlsBuildTook)

	return nil
}

func createNeededFolders(outputFile string, logFile string, moveFolder string) error {
	var err error
	outputFolder := common.GetFolderPathForFile(outputFile)
	logFolder := common.GetFolderPathForFile(logFile)
	if err = common.CreateAllFolders(outputFolder); err != nil {
		return err
	}
	if err = common.CreateAllFolders(logFolder); err != nil {
		return err
	}
	if moveFolder != "" {
		if err = common.CreateAllFolders(moveFolder); err != nil {
			return err
		}
	}
	return nil
}

func outputDDRJSONData(outputFile string, objects []ddrml.DrillReports) error {
	var err error
	jsonStart := time.Now()
	if err = ddrml.BuildJsonFileForDrilling(outputFile, objects); err != nil {
		zap.S().Error("Error:", err.Error())
		return err
	}
	jsonBuildTook := time.Since(jsonStart)
	zap.S().Infof("JSON build took:%s", jsonBuildTook)

	return nil
}

func main() {

	outputFile := flag.String("OUTPUT_FILE", "", "Specifies name and path of result file, e.g. c:\\temp\\output_data.xlsx or c:\\temp\\output_data.json")
	xmlFolder := flag.String("XML_FOLDER", "", "Specifies the path to the folder containing DDR xml files to process")
	logFile := flag.String("LOG_FILE", "", "Path and file name to use for the logging file,e.g. C:\\temp\\processing_zap.S().txt")
	moveFolder := flag.String("MOVE_FOLDER", "", "Specifies a path to a folder where xml files that has been processed should be moved to, if left empty files are not moved")
	outputFormat := flag.String("OUTPUT_FORMAT", "excel", "Specifies output format to use either excel (default), csv or json")
	appendTimeToFileName := flag.Bool("APPEND_TIME2FILENAME", false, "If set to true exection will always add a timestamp to the output file name")
	oneFilePerSheet := flag.Bool("ONE_FILE_PER_SHEET", false, "If set and output is excel, the program will generate one excel file per sheet as output")
	showVersion := flag.Bool("version", false, "If specified will print out the version information and then exit")

	flag.Parse()
	if *showVersion {
		fmt.Println("Version:", Version)
		fmt.Println("Build time:", Build)
		return
	}
	if *outputFile != "" && *xmlFolder != "" && *logFile != "" {
		processDDRFiles(*xmlFolder, *outputFile, *logFile, *moveFolder, *outputFormat,
			*appendTimeToFileName, *oneFilePerSheet)
	} else {
		fmt.Println("Missing required input parameters...")
		flag.PrintDefaults()
		return
	}
}

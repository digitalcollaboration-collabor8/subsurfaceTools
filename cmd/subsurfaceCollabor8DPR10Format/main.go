/*
* @Author: magsv
* @Date:   2017-11-24 09:22:48
* @Last Modified by:   magsv
* @Last Modified time: 2018-06-25 10:13:20
 */

package main

import (
	"flag"
	"fmt"

	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/common"
	"github.com/digitalcollaboration-collabor8/subsurfaceTools/pkg/dpr10"

	"go.uber.org/zap"

	"time"
)

var (
	Version string
	Build   string
)

func processXMLFiles(folderPath string, outputFile string, logFile string, moveFolder string,
	outputFormat string, appendTime2Filename bool, oneFilePerSheet bool) {
	var err error
	var objects []dpr10.WITSMLComposite
	var files []string
	start := time.Now()
	unmarshalStart := time.Now()
	//just create the folders needed if not existing
	if err = createNeededFolders(outputFile, logFile, moveFolder); err != nil {
		fmt.Printf("Failed in created needed folder:%s", err.Error())
		return
	}
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
	//read the xml files in the folder to structs
	if objects, err = dpr10.ReadProdXMLFiles2Struct(folderPath); err != nil {
		zap.S().Errorf("Failed in reading xml files in folder:%s, error:%s", folderPath, err.Error())
		return
	}
	unmarshalTook := time.Since(unmarshalStart)
	switch outputFormat {
	case "excel":
		err = outputExcelData(outputFile, objects, oneFilePerSheet, appendTime2Filename)
	case "json":
		err = outputJSONData(outputFile, objects)
	case "csv":
		err = dpr10.BuildCSVFileForProduction(outputFile, objects, ";")
	default:
		err = outputExcelData(outputFile, objects, oneFilePerSheet, appendTime2Filename)
	}
	if err != nil {
		zap.S().Errorf("Failed in outputting data to format:%s, error:%s", outputFormat, err.Error())
		return
	}

	totalTimeTook := time.Since(start)
	if moveFolder != "" {
		zap.S().Infof("Moving files in folder:%s to %s", folderPath, moveFolder)
		//Need to move the data to the move folder
		//make sure the destination folder exists
		if err = common.CreateAllFolders(moveFolder); err != nil {
			zap.S().Errorf("Failed in creating move folder:%s", err.Error())
			return
		}
		//get the list of the files
		if files, err = common.GetFilesWithExtension(folderPath, "*.xml"); err != nil {
			zap.S().Error("Failed in reading file in folder:", err.Error())
			return
		}
		if err = common.MoveFiles(files, moveFolder); err != nil {
			zap.S().Errorf("Failed in moving files to folder:%s, error:%s", moveFolder, err.Error())
			return
		}

	}
	zap.S().Info("Finished processing results to:", outputFile)
	zap.S().Infof("XML File parsing took:%s", unmarshalTook)
	zap.S().Infof("Total time used:%s", totalTimeTook)
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

func outputJSONData(outputFile string, objects []dpr10.WITSMLComposite) error {
	var err error
	jsonStart := time.Now()
	if err = dpr10.BuildJsonFileForProduction(outputFile, objects); err != nil {
		zap.S().Error("Error:", err.Error())
		return err
	}
	jsonBuildTook := time.Since(jsonStart)
	zap.S().Infof("JSON build took:%s", jsonBuildTook)

	return nil
}

func outputExcelData(outputFile string, objects []dpr10.WITSMLComposite,
	oneFilePerSheet bool, appendTime2Filename bool) error {
	var err error
	xlsStart := time.Now()
	if err = dpr10.BuildXLSFileForProduction(outputFile, objects, oneFilePerSheet, appendTime2Filename); err != nil {
		zap.S().Error("Error:", err.Error())
		return err
	}
	xlsBuildTook := time.Since(xlsStart)
	zap.S().Infof("Excel build took:%s", xlsBuildTook)

	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "If specified will print out the version information and then exit")

	outputFile := flag.String("OUTPUT_FILE", "", "Specifies name and path of result file, e.g. c:\\temp\\output_data.xlsx or c:\\temp\\output_data.json")
	xmlFolder := flag.String("XML_FOLDER", "", "Specifies the path to the folder containing MPMRL Government xml files to process")
	logFile := flag.String("LOG_FILE", "", "Path and file name to use for the logging file,e.g. C:\\temp\\processing_zap.S().txt")
	moveFolder := flag.String("MOVE_FOLDER", "", "Specifies a path to a folder where xml files that has been processed should be moved to, if left empty files are not moved")
	outputFormat := flag.String("OUTPUT_FORMAT", "excel", "Specifies output format to use either excel (default),csv or json")
	appendTimeToFileName := flag.Bool("APPEND_TIME2FILENAME", false, "If set to true exection will always add a timestamp to the output file name")
	oneFilePerSheet := flag.Bool("ONE_FILE_PER_SHEET", false, "If set and output is excel, the program will generate one excel file per sheet as output")
	flag.Parse()
	if *showVersion {
		fmt.Println("Version:", Version)
		fmt.Println("Build time:", Build)
		return
	}
	if *outputFile != "" && *xmlFolder != "" && *logFile != "" {
		processXMLFiles(*xmlFolder, *outputFile, *logFile, *moveFolder, *outputFormat, *appendTimeToFileName, *oneFilePerSheet)
	} else {
		fmt.Println("Missing required input parameters...")
		flag.PrintDefaults()
		return
	}

}

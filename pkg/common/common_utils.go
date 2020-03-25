/*
* @Author: magsv
* @Date:   2016-02-15 10:20:00
* @Last Modified by:   magsv
* @Last Modified time: 2018-11-08 12:04:15
 */
package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func InitializeFolders(outputFolder string, logFile string) error {
	var err error
	if err = CreateAllFolders(outputFolder); err != nil {
		return err
	}

	return nil
}

//Converts a list of datasets to json
func DatasetsToJson(datasets []DataSet) ([]byte, error) {
	var dSetsJson []DataSetJson
	//just move everything over to a column valuebase json set
	for x := 0; x < len(datasets); x++ {
		dSetJson := DataSetJson{HeadersName: datasets[x].HeadersName, Name: datasets[x].Name}
		//process row by row and add the header names
		rows := []RowDataJson{}
		for y := 0; y < len(datasets[x].Rows); y++ {
			row := RowDataJson{}
			for z := 0; z < len(datasets[x].Rows[y].Columns); z++ {
				//check type of value
				column := ColumnDataJson{Name: datasets[x].HeadersName[z]}
				if datasets[x].Rows[y].Columns[z].IsFloat {
					column.Value = datasets[x].Rows[y].Columns[z].FloatVal
				} else if datasets[x].Rows[y].Columns[z].IsInt {
					column.Value = datasets[x].Rows[y].Columns[z].IntVal
				} else if datasets[x].Rows[y].Columns[z].IsStr {
					column.Value = datasets[x].Rows[y].Columns[z].StrVal
				} else if datasets[x].Rows[y].Columns[z].IsTime {
					column.Value = datasets[x].Rows[y].Columns[z].TimeValue
				} else if datasets[x].Rows[y].Columns[z].IsEmptyColumn {
					column.Value = ""
				} else {
					column.Value = ""
				}
				row.Columns = append(row.Columns, column)
				//row.Columns=append(row.Columns,ColumnDataJson{Name:datasets[x].HeadersName[z],Value:datasets[x].Rows[y].Columns[z].})
			}
			rows = append(rows, row)
		}
		dSetJson.Rows = rows
		dSetsJson = append(dSetsJson, dSetJson)
	}

	return json.MarshalIndent(dSetsJson, "", "    ")
}

func TimeToString(tObj time.Time, outputFormat string) string {

	return tObj.Format(outputFormat)
}

func GetTimeZone(tObj *time.Time) string {
	zone, _ := tObj.Zone()
	return zone
}

func DatasetsToCsv(datasets []DataSet, outputFile string, discriminator string) error {
	var err error
	var outputFolder, fileNameOnly, ext, dataOutFile string
	for i := 0; i < len(datasets); i++ {
		//get the outputfolder to store in
		outputFolder = GetFolderPathForFile(outputFile)
		//get the file name without extension
		fileNameOnly, ext = GetFileNameAndExtension(outputFile)
		//build the new filename using header name and removing
		dataOutFile = outputFolder + string(os.PathSeparator) + fileNameOnly + "_" + datasets[i].Name + ext
		if err = DatasetToCsv(datasets[i], dataOutFile, discriminator); err != nil {
			return err
		}
		zap.S().Infof("Wrote csv data for dataset name:%s to location:%s", datasets[i].Name, dataOutFile)
	}
	return nil
}

func DatasetToCsv(dataset DataSet, outputFile string, discriminator string) error {
	var dataBuffer bytes.Buffer
	//first write the headers
	for i := 0; i < len(dataset.HeadersName); i++ {
		if i > 0 {
			dataBuffer.WriteString(discriminator)
		}
		dataBuffer.WriteString(dataset.HeadersName[i])
	}
	//add the break line
	dataBuffer.WriteString("\n")
	//write the data
	for x := 0; x < len(dataset.Rows); x++ {

		for y := 0; y < len(dataset.Rows[x].Columns); y++ {
			//now have the columne write it to the buffer
			if y > 0 {
				dataBuffer.WriteString(discriminator)
			}
			if dataset.Rows[x].Columns[y].IsEmptyColumn {
				dataBuffer.WriteString("\"\"")
			} else if dataset.Rows[x].Columns[y].IsFloat {
				dataBuffer.WriteString(strconv.FormatFloat(dataset.Rows[x].Columns[y].FloatVal, 'f', -1, 64))

			} else if dataset.Rows[x].Columns[y].IsInt {
				dataBuffer.WriteString(strconv.Itoa(dataset.Rows[x].Columns[y].IntVal))

			} else if dataset.Rows[x].Columns[y].IsTime {
				dataBuffer.WriteString(dataset.Rows[x].Columns[y].TimeValue.Format(datelayout_out_csv))

			} else {
				//handle the string value
				dataBuffer.WriteString("\"" + dataset.Rows[x].Columns[y].StrVal + "\"")
			}

		}
		dataBuffer.WriteString("\n")
	}
	//write the databuffer data as a byte array
	return Write2File(outputFile, dataBuffer.Bytes())

}

func CreateUUID() string {

	u1 := uuid.NewV4()
	return u1.String()
}

//function will take the current time and substract the given number of days and return a start and end time
func RollDays(days int) (time.Time, time.Time) {
	var start, end, now time.Time
	now = time.Now()
	end = now.AddDate(0, 0, 1)
	start = now.AddDate(0, 0, 0-days)
	return start, end
}

//function used to take a time object and format it to a string format in the form of yyyy-mm-dd
func FormatTime2QueryDayString(tObj time.Time) string {

	return tObj.Format(daylayout_query)
}

//function used to take a time object and format it to a string format in the form of mm
func FormatTime2QueryMonthString(tObj time.Time) string {

	return tObj.Format(monthlayout_query)
}

//function used to take a time object and format it to a string format in the form of yyyy
func FormatTime2QueryYearString(tObj time.Time) string {

	return tObj.Format(yearlayout_query)
}

//function will take the current time and substract the given number of month and return a start and end time
func RollMonths(months int) (time.Time, time.Time) {
	var start, end, now time.Time
	now = time.Now()
	end = now.AddDate(0, 1, 0)
	start = now.AddDate(0, 0-months, 0)
	return start, end
}

func XSDDateString2Time(dateString string) (time.Time, error) {
	return time.Parse(xsdDateLayout, dateString)
}

func XSDDateTimeString2Time(dateTimeString string) (time.Time, error) {
	return time.Parse(xsdDateTimeLayout, dateTimeString)
}

func CreateAllFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)

}

func ReadFile(filePath string) ([]byte, error) {
	file, e := ioutil.ReadFile(filePath)

	if e != nil {
		zap.S().Fatal(e)
		return nil, e
	}
	return file, nil
}

func Write2File(filePath string, data []byte) error {
	var err error
	var f *os.File
	if f, err = os.Create(filePath); err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil

}

func MoveFiles(files []string, move2Folder string) error {
	var err error
	for i := 0; i < len(files); i++ {
		currentFileName := GetFileName(files[i])
		newOutput := move2Folder + string(os.PathSeparator) + currentFileName
		if err = os.Rename(files[i], newOutput); err != nil {
			return err

		}
	}
	return nil
}

func getTimeAsString() string {
	var currentTime time.Time

	currentTime = time.Now()

	formattedTime := currentTime.Format(datelayout_out)

	return formattedTime
}

func GetFileName(filePath string) string {
	return filepath.Base(filePath)
}

func GetFilesWithExtension(folder string, extension string) ([]string, error) {
	var files []string
	var err error
	fileSearch := folder + string(os.PathSeparator) + extension
	if files, err = filepath.Glob(fileSearch); err != nil {
		return files, err
	}
	return files, nil
}

func GetFolderPathForFile(filePath string) string {
	return filepath.Dir(filePath)
}

//returns filename without extension and extension
func GetFileNameAndExtension(filePath string) (string, string) {
	extension := filepath.Ext(filePath)
	fileName := GetFileName(filePath)
	fileNameWithoutExt := strings.Replace(fileName, extension, "", -1)
	return fileNameWithoutExt, extension

}

/*appends a timestamp to a filename with the given filepath*/
func AppendTimeAndDateToFile(filepath string) string {
	folderPath := GetFolderPathForFile(filepath)
	fileName, extension := GetFileNameAndExtension(filepath)
	//now create the timestamp
	time2Add := getTimeAsString()
	return folderPath + string(os.PathSeparator) + fileName + "_" + time2Add + extension
}

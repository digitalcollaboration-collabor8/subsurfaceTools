/*
* @Author: magsv
* @Date:   2017-08-11 09:32:44
* @Last Modified by:   magsv
* @Last Modified time: 2018-08-24 08:53:00
 */

package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

func CreateWorkbookFromDataSet(filepath string, datasets []DataSet, oneFilePerSheet bool, appendTimeInName bool) error {
	var err error
	var outputFolder, fileNameOnly, ext, dataOutFile string
	//if we have oneFilePerSheetSet we need to loop us through the dataset and create one by one
	zap.S().Infof("One per file for excel is set to:%s", oneFilePerSheet)
	if oneFilePerSheet {
		//create the base filename
		//get the outputfolder to store in
		outputFolder = GetFolderPathForFile(filepath)
		//get the file name without extension
		fileNameOnly, ext = GetFileNameAndExtension(filepath)
		for i := 0; i < len(datasets); i++ {

			//build the new filename using header name and removing
			dataOutFile = outputFolder + string(os.PathSeparator) + fileNameOnly + "_" + datasets[i].Name + ext
			//just create a new dataset to be processed based on going through each entry

			if err = createWorkbookFromDataSet(dataOutFile, []DataSet{datasets[i]}); err != nil {
				return err
			}
			zap.S().Infof("Wrote excel data for dataset name:%s to location:%s", datasets[i].Name, dataOutFile)
		}
	} else {
		outputFolder = GetFolderPathForFile(filepath)
		//get the file name without extension
		fileNameOnly, ext = GetFileNameAndExtension(filepath)
		dataOutFile = outputFolder + string(os.PathSeparator) + fileNameOnly + ext
		//just create a new dataset to be processed based on going through each entry

		if err = createWorkbookFromDataSet(dataOutFile, datasets); err != nil {
			return err
		}
		zap.S().Infof("Wrote excel data for dataset name to location:%s", dataOutFile)
	}
	return err
}

//Creates a new excel file where each dataset that is feed into this function is
//represented as a worksheet
func createWorkbookFromDataSet(filepath string, datasets []DataSet) error {
	var file *xlsx.File
	var err error
	var sheet *xlsx.Sheet
	var dTimOptions xlsx.DateTimeOptions
	loc := time.Now().Location()
	dTimOptions.Location = loc
	dTimOptions.ExcelTimeFormat = xlsx.DefaultDateTimeFormat
	file = xlsx.NewFile()
	for i := 0; i < len(datasets); i++ {
		//create the sheet
		if sheet, err = file.AddSheet(datasets[i].Name); err != nil {
			return err
		}
		//add the headers
		row := sheet.AddRow()
		for s := 0; s < len(datasets[i].HeadersName); s++ {
			addHeaderCellData(row, datasets[i].HeadersName[s])
		}
		//process the rowdata
		for x := 0; x < len(datasets[i].Rows); x++ {
			//add the new row
			row := sheet.AddRow()
			for y := 0; y < len(datasets[i].Rows[x].Columns); y++ {
				if datasets[i].Rows[x].Columns[y].IsStr {
					addCellData(row, datasets[i].Rows[x].Columns[y].StrVal)
				} else if datasets[i].Rows[x].Columns[y].IsInt {
					addCellDataInt(row, datasets[i].Rows[x].Columns[y].IntVal)
				} else if datasets[i].Rows[x].Columns[y].IsFloat {
					addCellDataFloat(row, datasets[i].Rows[x].Columns[y].FloatVal)
				} else if datasets[i].Rows[x].Columns[y].IsTime {
					//just add the cell data if it is not null otherwise add empty column
					timeValue := datasets[i].Rows[x].Columns[y].TimeValue
					if timeValue.IsZero() {
						//add the empty column
						addEmptyCell(row)
					} else {
						addCellDataTime(row, datasets[i].Rows[x].Columns[y].TimeValue, dTimOptions)
					}

				} else if datasets[i].Rows[x].Columns[y].IsEmptyColumn {
					addEmptyCell(row)
				} else {
					return errors.New(fmt.Sprintf("Undefined datatype found for data set with name: %s, row number:%d,column number:%d", datasets[i].Name, x, y))
				}
			}
		}
	}
	return file.Save(filepath)
}

//adds a cell formatted with center adjustement and bold txt
func addHeaderCellData(row *xlsx.Row, data string) {
	centerHalign := *xlsx.DefaultAlignment()
	centerHalign.Horizontal = "center"
	style := xlsx.NewStyle()
	font := *xlsx.NewFont(12, "Verdana")
	font.Bold = true
	style.Font = font
	style.Alignment = centerHalign
	cell := row.AddCell()
	cell.Value = data
	cell.SetStyle(style)
}

//builds the header row and columns for a worksheet
func buildObjSheetHeaders(row *xlsx.Row, headers []string) {
	for i := 0; i < len(headers); i++ {
		addHeaderCellData(row, headers[i])
	}

}

//adds a cell with the type of string to a row
func addCellData(row *xlsx.Row, data string) {

	cell := row.AddCell()
	cell.Value = data

}

//adds a cell with the type of int to a row
func addCellDataInt(row *xlsx.Row, data int) {

	cell := row.AddCell()
	cell.SetInt(data)

}

//adds an empty cell to a row
func addEmptyCell(row *xlsx.Row) {
	row.AddCell()
}

//adds a cell with the type of time to a row
func addCellDataTime(row *xlsx.Row, data time.Time, dTimOptions xlsx.DateTimeOptions) {
	cell := row.AddCell()
	//cell.SetDateTime(data)
	cell.SetDateWithOptions(data, dTimOptions)
	//cell.Value = TimeToString(data, xsdDateTimeLayout)
}

//adds a cell with the type of float to a row
func addCellDataFloat(row *xlsx.Row, data float64) {
	cell := row.AddCell()
	cell.SetFloat(data)
}

//adds a cell with the type of string and set the style of it
func addCellDataWithStyle(row *xlsx.Row, data string, style *xlsx.Style) {
	cell := row.AddCell()
	cell.Value = data
	cell.SetStyle(style)
}

# subsurfaceCollabor8DDR2Format


The following client can process a set of DDRML xml files and convert them to e.g. one spreadsheet with data


### Supported configuration parameters

- **XML_FOLDER** -> Path to the folder containing xml files to process
- **OUTPUT_FILE** -> Path and name of the output file
- **LOG_FILE** -> Path and name of the log file
- **MOVE_FOLDER** -> If specified after successful processing of files in the specified XML_FOLDER, the input files that has been processed will be moved to this folder
- **OUTPUT_FORMAT** -> Format to be used to the processing result either excel (default), json or csv, meaning that if excel is use the program will create an excel file with data as the path specified by the output_file option. Note that if csv is choosen as ouput format one file per datatype will be generated as a csv file cannot hold this information in one single file.

### Example processing a set of DDR xml files

#### To write data to an excel file:

subsurfaceCollabor8DDR2Format.exe -XML_FOLDER="C:\Temp\DDR" -OUTPUT_FILE="C:\Temp\DDR.xlsx" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="excel"

#### To write data to a json file

subsurfaceCollabor8DDR2Format.exe -XML_FOLDER="C:\Temp\DDR" -OUTPUT_FILE="C:\Temp\DDR.json" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="json"

#### To write data to a csv file

subsurfaceCollabor8DDR2Format.exe -XML_FOLDER="C:\Temp\DDR" -OUTPUT_FILE="C:\Temp\DDR.csv" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="csv"
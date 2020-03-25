# subsurfaceCollabor8DPR10Format


The following client can take a set of DPR 1.0 xml files and process them over to e.g. a spreadsheet

### Supported configuration parameters

- **XML_FOLDER** -> Path to the folder containing xml files to process
- **OUTPUT_FILE** -> Path and name of the output file
- **LOG_FILE** -> Path and name of the log file
- **MOVE_FOLDER** -> If specified after successful processing of files in the specified XML_FOLDER, the input files that has been processed will be moved to this folder
- **OUTPUT_FORMAT** -> Format to be used to the processing result either excel (default), csv or json, meaning that if excel is use the program will create an excel file with data as the path specified by the output_file option
- **APPEND_TIME2FILENAME** -> if set, a timestamp will be added to the output filename

### Example processing a set of DPR 1.0 xml files

#### To write data to an excel file:

subsurfaceCollabor8DPR10Format.exe -XML_FOLDER="C:\Temp\DPR1" -OUTPUT_FILE="C:\Temp\DPR10_RESULTS.xlsx" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="excel"

#### To write data to a json file

subsurfaceCollabor8DPR10Format.exe -XML_FOLDER="C:\Temp\DPR10" -OUTPUT_FILE="C:\Temp\DPR10_RESULTS.json" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="json"

#### To write data to a csv file

subsurfaceCollabor8DPR10Format.exe -XML_FOLDER="C:\Temp\DPR10" -OUTPUT_FILE="C:\Temp\DPR10_RESULTS.csv" -LOG_FILE="C:\Temp\Result_log.txt" -OUTPUT_FORMAT="csv"

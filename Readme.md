
## Subsurface Collabor8

The following project is a common library for working against the Collabor8 platform. It contains tools for downloading physical report files from the system and also tools to process raw xml file over to other formats like e.g. Excel or CSV.


## Building

### Dependencies

In order to build the following libraries the following is needed

1. Go language >= 1.13
2. Make sure that everything is checked out under <GOPATH>/src/github.com/EPIM-Association

### Building 

#### Using standard Go routines

1. Enter **go install ./..** in this folder or if you would like to have versions associated with each build use the ldflags option e.g. **go install -ldflags "-w -s -X main.Version=1. -X main.Build=2020-08-10" ./...**
2. Generated articacts will go into the default go bin folder

#### Using make

This package comes packaged with a simple makefile capable of building all of the clients. The build process using the makefile is relying on that ytou have both make and git installed locally.

1. Enter **make build** in this folder to build for your current platform
2. Enter **make linux** in this folder to build for a 64bit Linux platform
3. Enter **make windows** in this folder to build for a 64bit windows platform
4. Enter **make release** to build binaries for both linux and windows

## Configuration

The following environment variables need to be set to be able to authenticate against Azure and download physical files

1. AzureClientId - the client id from Azure
2. AzureClientSecret - the client secret from Azure
3. AzureTokenUrl - the token url from azure
4. AzureResourceId - the azure resource id to authenticate against
5. AzureFileDownloadUrl - the url from where to download files in azure
6. AzureSubscriptionKey - the service subscription key to use when calling the api's
7. AzureGraphUrl - the url for graph queries

The config folder contains an example of a download configuration file and also a sample on how to configure logging if something else than the default is wanted. The library utilises the zap logging library.

## Prebuild binaries - ready to use

Prebuild binaries are available as part of each release in Github.

## Running 

For running instructions check each cmd folder and its associated ReadMe file.





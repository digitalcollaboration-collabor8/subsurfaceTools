# subsurfaceCloudDownload


The following folder contains the source for the program to download physical files from the Collabor8 platform using the Azure cloud based services.


## Running

Just run the program from a command line or e.g. as a scheduled task. Below are the expected configuration parameters

- **configuration** -> path to the xml configuration file to use for downloading data 
- **logconfiguration** -> Optional path to a log configuration json file, if left out the default set-up will be used and logging will be directly to the stdout.
- **version** -> displays the build version and date for the client

Example of running, ./subsurfaceCloudDownload -configuration="./downloadConfig.xml" -logconfiguration="./logConfig.json"


## Configuration

The program accepts an xml file as its configuration containing information on what to download. Authentication is done using either environment variables of through the same configuration file.

### Authentication and Azure parameters 

Authentication can be configured using environment variables which is the default set-up or it can also be sent in as part of the xml download configuration file.


#### Configuration using environment variables

The following environment variables needs to be set if the program should function correctly (given that you use environment variables for configuration and not xml configuration).

1. AzureClientId - the client id from Azure
2. AzureClientSecret - the client secret from Azure
3. AzureTokenUrl - the token url from azure
4. AzureResourceId - the azure resource id to authenticate against
5. AzureFileDownloadUrl - the url from where to download files in azure
6. AzureSubscriptionKey - the service subscription key to use when calling the api's
7. AzureGraphUrl - the url for graph queries

#### Configuration using the xml configuration file

If the Azure needed client parameters haven't been set using the environment variables it is possible to configure the same using the xml configuration file.

Below is a sample of configuring the same parameters through the xml configuration file which is used to configure which assets to download data for.

```xml
<subsurface>
   <config>
	<clientId>XXXXX</clientId><!--the azure client id to use-->
	<clientSecret>YYYY</clientSecret><!-- the azure client secret to use-->
	<tokenUrl>SSSS</tokenUrl><!-- the Azure token url to use for authentication-->
	<resourceId>FFFFFF</resourceId> <!-- the resource id to authenticate against in Azure-->
	<fileDownloadUrl>GGGGG</fileDownloadUrl><!-- the file download url to use for downloading files-->
	<subscriptionKey>VVVVVV</subscriptionKey> <!--the subscription key to use for the API calls-->
	<graphUrl>BBBBBB</graphUrl><!--the url for running GraphQL queries>-->
   </config>
</subsurface>
```


### Download configuration

For examples of how to configure download of data, see the config/SampleCloudDownloadConfiguration.xml file
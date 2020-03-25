package cloud

import "testing"

var downloadDPR10XmlDataConfigDateFrom = `<subsurface>
<dpr>
	  <fieldName>ÅSGARD</fieldName>
	  <dateFrom>2020-02-24</dateFrom>
	  <dateTo>2020-02-29</dateTo>
	  <rollDays>10</rollDays>
	  <useUploadedFrom>false</useUploadedFrom>
	  <common>
		  <format>PDF</format>
		  <outputFolder>../../test/results</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
</subsurface>`

var downloadDPR10XmlDataConfigUploadedFrom = `<subsurface>
<dpr>
	  <fieldName>ÅSGARD</fieldName>
	  <rollDays>10</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>PDF</format>
		  <outputFolder>../../test/results</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
</subsurface>`

var downloadMPRMLGovXmlDataConfigDateFrom = `<subsurface>
<mprmlGov>
	  <fieldName>JOHAN SVERDRUP</fieldName>
	  <dateFrom>2019-11-01</dateFrom>
	  <dateTo>2019-12-31</dateTo>
	  <useUploadedFrom>false</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/MPRML_GOV_SVERDRUP</outputFolder>
		  <fileOutputPrefix>SVERDRUP_</fileOutputPrefix>
	  </common>
</mprmlGov>
</subsurface>`

var downloadMPRMLGovXmlDataConfigUploadedFrom = `<subsurface>
<mprmlGov>
	  <fieldName>JOHAN SVERDRUP</fieldName>
	  <dateFrom>2020-02-20</dateFrom>
	  <dateTo>2020-02-25</dateTo>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/MPRML_GOV_SVERDRUP</outputFolder>
		  <fileOutputPrefix>SVERDRUP_</fileOutputPrefix>
	  </common>
</mprmlGov>
</subsurface>`

var downloadDDRMLXmlDataConfigDateFrom = `<subsurface>
<ddrml>
			  <dateFrom>2019-01-01</dateFrom>
	          <dateTo>2019-01-20</dateTo>
	          
	          <common>
	              <format>XML</format>
	              <outputFolder>../../test/results/DDRML_DATE_FROM</outputFolder>
	              <fileOutputPrefix>DDRML</fileOutputPrefix>
	          </common>
		</ddrml>
</subsurface>`

var downloadDDRMLXmlDataConfigUploadedFrom = `<subsurface>
<ddrml>
			  <dateFrom>2019-01-01</dateFrom>
	          <dateTo>2019-01-20</dateTo>
	          <rollDays>30</rollDays>
	  		  <useUploadedFrom>true</useUploadedFrom>
	          <common>
	              <format>XML</format>
	              <outputFolder>../../test/results/DDRML_DATE_FROM</outputFolder>
	              <fileOutputPrefix>DDRML</fileOutputPrefix>
	          </common>
		</ddrml>
</subsurface>`

var testDownloadDDRML = `<subsurface>
<ddrml>
			  
	          <rollDays>10</rollDays>
	  		  <useUploadedFrom>true</useUploadedFrom>
	          <common>
	              <format>XML</format>
	              <outputFolder>../../test/results/DDRML_UPLOADED_FROM</outputFolder>
	              <fileOutputPrefix>DDRML_UPLOADED_FROM_</fileOutputPrefix>
	          </common>
		</ddrml>
		<ddrml>
			  
	          <rollDays>10</rollDays>
	  		  <useUploadedFrom>true</useUploadedFrom>
	          <common>
	              <format>PDF</format>
	              <outputFolder>../../test/results/DDRML_UPLOADED_FROM_PDF</outputFolder>
	              <fileOutputPrefix>DDRML_UPLOADED_FROM</fileOutputPrefix>
	          </common>
		</ddrml>
		<!--<ddrml>
		<dateFrom>2020-02-11</dateFrom>
		<dateTo>2020-03-13</dateTo>
		
		<common>
			<format>XML</format>
			<outputFolder>../../test/results/DDRML_PERIOD</outputFolder>
			<fileOutputPrefix>DDRML_PERIOD</fileOutputPrefix>
		</common>
  </ddrml>
  <ddrml>
  <dateFrom>2020-02-11</dateFrom>
  <dateTo>2020-03-13</dateTo>
		
		<common>
			<format>PDF</format>
			<outputFolder>../../test/results/DDRML_PERIOD_PDF</outputFolder>
			<fileOutputPrefix>DDRML_PERIOD</fileOutputPrefix>
		</common>
  </ddrml>-->
</subsurface>`

var testDownloadDDRMLPeriod = `<subsurface>

	<ddrml>
		<dateFrom>2020-02-17</dateFrom>
		<dateTo>2020-03-19</dateTo>
		<useUploadedFrom>false</useUploadedFrom>
		<common>
			<format>XML</format>
			<outputFolder>../../test/results/DDRML_PERIOD</outputFolder>
			<fileOutputPrefix>DDRML_PERIOD</fileOutputPrefix>
		</common>
  </ddrml>
  <ddrml>
	<dateFrom>2020-02-17</dateFrom>
	<dateTo>2020-03-19</dateTo>
	<useUploadedFrom>false</useUploadedFrom>
	<common>
		<format>XML</format>
		<outputFolder>../../test/results/DDRML_PERIOD_PDF</outputFolder>
		<fileOutputPrefix>DDRML_PERIOD</fileOutputPrefix>
	</common>
</ddrml>
</subsurface>`

var testDownloadMPRMLGov = `<subsurface>
<mprmlGov>
	  <fieldName>GINA KROG</fieldName>
	  <rollDays>10</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/GINAKROG_MPRMLGOV</outputFolder>
		  <fileOutputPrefix>GINAKROG_</fileOutputPrefix>
	  </common>
</mprmlGov>
<mprmlGov>
	  <fieldName>JOHAN SVERDRUP</fieldName>
	  <dateFrom>2019-02-28</dateFrom>
	  <dateTo>2019-03-28</dateTo>
	  <useUploadedFrom>false</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/SVERDRUP_MPRMLGOV</outputFolder>
		  <fileOutputPrefix>SVERDRUP_</fileOutputPrefix>
	  </common>
</mprmlGov>
</subsurface>`

var testDownloadDPR10 = `<subsurface>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <rollDays>1</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>PDF</format>
		  <outputFolder>../../test/results/GINA</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <rollDays>1</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/GINA</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <dateFrom>2012-04-30</dateFrom>
	  <dateTo>2012-05-01</dateTo>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/GINA</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <dateFrom>2012-04-30</dateFrom>
	  <dateTo>2012-05-01</dateTo>
	  <common>
		  <format>PDF</format>
		  <outputFolder>../../test/results/GINA</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
</subsurface>`

var testDownloadDPR10Created = `<subsurface>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <rollDays>10</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>PDF</format>
		  <outputFolder>../../test/results/GINA_PDF</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>
<dpr>
	  <fieldName>GINA KROG</fieldName>
	  <rollDays>10</rollDays>
	  <useUploadedFrom>true</useUploadedFrom>
	  <common>
		  <format>XML</format>
		  <outputFolder>../../test/results/GINA_XML</outputFolder>
		  <fileOutputPrefix>ASGARD_</fileOutputPrefix>
	  </common>
</dpr>

</subsurface>`

func TestBuildDPR10QueryNoUploadedFrom(t *testing.T) {
	var data, queryDPR10, queryDPR20 []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadDPR10XmlDataConfigDateFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of dpr 10 config:%s", err.Error())
	}
	fQueryDPR10 := createDPR10Query(cConfig.DPRS[0], "token")
	fQueryDPR20 := createDPR20Query(cConfig.DPRS[0], "token")
	if queryDPR10, err = BuildQueryForAssetUsingPeriod(fQueryDPR10); err != nil {
		t.Errorf("Failed in creating dpr10 query:%s", err.Error())
	}
	if queryDPR20, err = BuildQueryForAssetUsingPeriod(fQueryDPR20); err != nil {
		t.Errorf("Failed in creating dpr20 query:%s", queryDPR20)
	}
	t.Logf("QUERYDPR10:%s", string(queryDPR10))
	t.Logf("QUERYDPR20:%s", string(queryDPR20))
}

func TestBuildDPR10QueryUploadedFrom(t *testing.T) {
	var data, queryDPR10, queryDPR20 []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadDPR10XmlDataConfigUploadedFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of dpr 10 config:%s", err.Error())
	}
	fQueryDPR10 := createDPR10Query(cConfig.DPRS[0], "token")
	fQueryDPR20 := createDPR20Query(cConfig.DPRS[0], "token")
	if queryDPR10, err = BuildQueryForAssetUsingCreated(fQueryDPR10); err != nil {
		t.Errorf("Failed in creating dpr10 query:%s", err.Error())
	}
	if queryDPR20, err = BuildQueryForAssetUsingCreated(fQueryDPR20); err != nil {
		t.Errorf("Failed in creating dpr20 query:%s", queryDPR20)
	}
	t.Logf("QUERYDPR10:%s", string(queryDPR10))
	t.Logf("QUERYDPR20:%s", string(queryDPR20))
}

//TestDownloadDPR10FromConfiguration will test to download
//DPR 1.0 data from a xml configuration entry
func TestDownloadDPR10FromConfiguration(t *testing.T) {
	var data []byte
	var err error
	var errors []error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(testDownloadDPR10)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of dpr 10 config:%s", err.Error())
	}
	if errors = ProcessAndRunDownload(cConfig); len(errors) > 0 {
		for i := 0; i < len(errors); i++ {
			t.Errorf("Failed in download:%s", errors[i].Error())
		}
	} else {
		t.Logf("Finished download of dpr10")
	}
}

//TestDownloadDPR10FromConfigurationCreated will test to download
//DPR 1.0 data from a xml configuration entry and looking for data modified within the last 10 days
func TestDownloadDPR10FromConfigurationCreated(t *testing.T) {
	var data []byte
	var err error
	var errors []error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(testDownloadDPR10Created)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of dpr 10 config:%s", err.Error())
	}
	if errors = ProcessAndRunDownload(cConfig); len(errors) > 0 {
		for i := 0; i < len(errors); i++ {
			t.Errorf("Failed in download:%s", errors[i].Error())
		}
	} else {
		t.Logf("Finished download of dpr10")
	}
}

func TestBuildMPRMLGovQueryNoUploadedFrom(t *testing.T) {
	var data, queryMPRMLGov []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadMPRMLGovXmlDataConfigDateFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of mprmlgov config:%s", err.Error())
	}
	fQueryMPRMLGov := createMPRMLGovQuery(cConfig.MPRGovs[0], "token")
	if queryMPRMLGov, err = BuildQueryForAssetUsingPeriod(fQueryMPRMLGov); err != nil {
		t.Errorf("Failed in   creating mprmlgov query:%s", err.Error())
	}

	t.Logf("QUERYMPRMLGov:%s", string(queryMPRMLGov))
}

func TestBuildMPRMLGovQueryUploadedFrom(t *testing.T) {
	var data, queryMPRMLGov []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadMPRMLGovXmlDataConfigUploadedFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of mprmlgov config:%s", err.Error())
	}
	fQueryMPRMLGov := createMPRMLGovQuery(cConfig.MPRGovs[0], "token")
	if queryMPRMLGov, err = BuildQueryForAssetUsingCreated(fQueryMPRMLGov); err != nil {
		t.Errorf("Failed in creating mprmlgov query:%s", err.Error())
	}

	t.Logf("QUERYMPRMLGov:%s", string(queryMPRMLGov))
}

func TestDownloadMPRMLGovFromConfigurationUsingPeriod(t *testing.T) {
	var data []byte
	var err error
	var errors []error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadMPRMLGovXmlDataConfigDateFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of mprml gov config:%s", err.Error())
	}
	if errors = ProcessAndRunDownload(cConfig); len(errors) > 0 {
		for i := 0; i < len(errors); i++ {
			t.Errorf("Failed in download:%s", errors[i].Error())
		}
	} else {
		t.Logf("Finished download of mprmlgov")
	}
}

func TestDownloadMPRMLGovFromConfigurationUsingUploadedFrom(t *testing.T) {
	var data []byte
	var err error
	var errors []error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadMPRMLGovXmlDataConfigUploadedFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of mprml gov config:%s", err.Error())
	}
	if errors = ProcessAndRunDownload(cConfig); len(errors) > 0 {
		for i := 0; i < len(errors); i++ {
			t.Errorf("Failed in download:%s", errors[i].Error())
		}
	} else {
		t.Logf("Finished download of mprmlgov")
	}
}

func TestBuildDDRMLQueryNoUploadedFrom(t *testing.T) {
	var data, queryDDRML []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadDDRMLXmlDataConfigDateFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of ddrml config:%s", err.Error())
	}
	fQueryDDRML := createDDRMLQuery(cConfig.DDRMLS[0], "token")
	if queryDDRML, err = BuildQueryForAssetUsingPeriod(fQueryDDRML); err != nil {
		t.Errorf("Failed in creating ddrml query:%s", err.Error())
	}

	t.Logf("QUERYDDRML:%s", string(queryDDRML))
}

func TestBuildDDRMLQueryUploadedFrom(t *testing.T) {
	var data, queryDDRML []byte
	var err error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(downloadDDRMLXmlDataConfigUploadedFrom)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of ddrml config:%s", err.Error())
	}
	fQueryDDRML := createDDRMLQuery(cConfig.DDRMLS[0], "token")
	if queryDDRML, err = BuildQueryForAssetUsingPeriod(fQueryDDRML); err != nil {
		t.Errorf("Failed in creating ddrml query:%s", err.Error())
	}

	t.Logf("QUERYDDRML:%s", string(queryDDRML))
}

//TestDownloadDDRMLFromConfigurationForPeriod will test to download
//ddrml files for a given period of time
func TestDownloadDDRMLFromConfigurationForPeriod(t *testing.T) {
	var data []byte
	var err error
	var errors []error
	var cConfig CloudDownload
	//first test to unmarshal the xml config
	data = []byte(testDownloadDDRMLPeriod)
	if cConfig, err = CloudConfigArrayToStruct(data); err != nil {
		t.Errorf("Failed in uinmarshal of ddrml config:%s", err.Error())
	}
	if errors = ProcessAndRunDownload(cConfig); len(errors) > 0 {
		for i := 0; i < len(errors); i++ {
			t.Errorf("Failed in download:%s", errors[i].Error())
		}
	} else {
		t.Logf("Finished download of ddrml")
	}
}

/* ****************************************************************************
Copyright © 2015-2017 by opensoach. All rights reserved.
Filename    : config.go
File-type   : go-lang-1.6.2 source code file.
Compiler    : go version go1.6.2 linux/amd64

Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description :
- Loads configuration parameters from serverconfigparams.xml file.
- At present, serverconfigparams.xml file contains only following server configuration parameters:
ServerWebSocketPort (server websocket port), ServerWebServicePort (server webserver port)
and log configuration parameters.

<?xml version="1.0" encoding="UTF-8" ?>
<ConfigParams>
	<ProductRoot>/opt/sagacity/vrom</ProductRoot>
	<ServerConfigParams>
		<ServerWebServicePort>80</ServerWebServicePort>
	</ServerConfigParams>
	<LogConfigParams>
		<LogDir>/opt/sagacity/vrom/logs/server_logs/</LogDir>
		<LogFileNamePrefix>server.log</LogFileNamePrefix>
		<LogFile>/opt/sagacity/vrom/logs/server_logs/server.log</LogFile>
		<LogLevel>DEBUG</LogLevel>
		<LogFileSize>20971520</LogFileSize>
		<LogMaxFiles>10</LogMaxFiles>
	</LogConfigParams>
	<MysqlDBConfigParams>
		<DBName>vromdb</DBName>
		<DBServerName>localhost</DBServerName>
		<DBServerPort>3306</DBServerPort>
		<ConnectionString>root:Welcome@123@tcp(localhost:3306)/vromdb?parseTime=true</ConnectionString>
	</MysqlDBConfigParams>
</ConfigParams>

Version History
Version     : 1.0
Author      : sameer oak (sameer.oak@opensoach.com)
Description : Initial version
**************************************************************************** */
package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/* ****************************************************************************
Prototype   :
func init()

Arguments   : na

Description :
- Reads configuration parameters xml and loads a variable of type ServerConfig and LogConfig.
- Unmarshalling is done into variable of type ConfigParams.

Assumptions : na

TODO        :

Return Value: na. init() neither takes any argument nor returns any value.
**************************************************************************** */
func Init() *ConfigParams {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Current working directory: %s\n", dir)

	var fileSeperator string
	var configFileName string

	if runtime.GOOS == "windows" {
		fileSeperator = "\\"
		configFileName = "appWin.config"
	} else {
		fileSeperator = "/"
		configFileName = "appLinux.config"
	}

	filePathItems := []string{dir, configFileName}

	fileFullPath := strings.Join(filePathItems, fileSeperator)

	fmt.Printf("Reading configuration from file: %s\n", fileFullPath)

	configXML, err := os.Open(fileFullPath)

	if err != nil {
		fmt.Println("Error: Missing server configuration.")
		fmt.Println(err)
		os.Exit(1)
	}
	defer configXML.Close()

	configParameters := new(ConfigParams)
	xmlData, _ := ioutil.ReadAll(configXML)
	xml.Unmarshal(xmlData, configParameters)

	return configParameters
}

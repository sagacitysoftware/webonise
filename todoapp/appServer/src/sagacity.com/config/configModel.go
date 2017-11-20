/* ****************************************************************************
Copyright © 2015-2017 by Sagacity Software Pvt. Ltd.. All rights reserved.
Package     : sagacity.com/config
Filename    : sagacity.com/config/configModel.go
File-type   : golang-1.6.2 source code file

Compiler/Runtime: go, golang version go1.9 linux/amd64 (linux/x86_64)

Version History
Version     : 1.0
Author      : sameer oak (sameer.oak@sagacitysoftware.co.in)
Description : Initial version
- Models of config package.
**************************************************************************** */
package config

import (
	"encoding/xml"
)

/*
- ConfigParams contains all the configuration parameters.
- Parses /opt/opensoach/hkt/serverconfigparams.xml.
- At present, there exist following server configuration parameters.
1> network configuration parameters
2> log configuration parameters
3> MySQL DB configuration parameters
*/
type ConfigParams struct	{
	XMLName            xml.Name `xml:"ConfigParams"`
	ServerConfigParams ServerConfig
	LogConfigParams LogConfig
	MysqlDBConfigParams MysqlDBConfig
}

/*
- ServerConfig contains all server configuration parameters to be
fetched before all the server side subsystems start.
- At present, there exist only network server configuration parameters in
the server configuration.
- They're:
ServerWebSocketPort, ServerWebServicePort.
*/
type ServerConfig struct	{
	ServerWebSocketPort  int16 `xml:"ServerWebSocketPort"`
	ServerWebServicePort int16 `xml:"ServerWebServicePort"`
}

/* log configuration parameters */
type LogConfig struct	{
	LogDir      string `xml:"LogDir"`
	LogFileNamePrefix string `xml:"LogFileNamePrefix"`
	LogFile     string `xml:"LogFile"`
	LogLevel    string `xml:"LogLevel"`
	LogFileSize int64  `xml:"LogFileSize"`
	LogMaxFiles int8   `xml:"LogMaxFiles"`
}

/* MySQL DB configuration parameters */
type MysqlDBConfig struct	{
	DBName string				`xml:"DBName"`
	DBServerName string			`xml:"DBServerName"`
	DBServerPort uint16			`xml:"DBServerPort"`
	ConnectionString string		`xml:"ConnectionString"`
}

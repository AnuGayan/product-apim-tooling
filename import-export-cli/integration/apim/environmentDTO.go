/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package apim

// Environment : Environment DTO
type Environment struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"displayName"`
	Description string  `json:"description"`
	IsReadOnly  bool    `json:"isReadOnly"`
	VHosts      []VHost `json:"vhosts"`
	Provider    string  `json:"provider"`
}

type VHost struct {
	Host        string `json:"host"`
	HttpContext string `json:"httpContext"`
	HttpPort    int    `json:"httpPort"`
	HttpsPort   int    `json:"httpsPort"`
	WsPort      int    `json:"wsPort"`
	WssPort     int    `json:"wssPort"`
}

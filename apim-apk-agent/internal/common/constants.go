/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

// Package common contains the constants, utility methods shared across two or many packages.
package common

import "github.com/wso2/product-apim-tooling/apim-apk-agent/config"

const (
	// OrganizationID query parameter key.
	organizationID string = "organizationId"
	// OrganizationID query parameter value used when the global adapter is enabled and it is a shared gateway.
	commonOrganizationIDValue string = "ALL"
)

// PopulateQueryParamForOrganizationID add the query parameter "organizationId" with the value of "ALL"
func PopulateQueryParamForOrganizationID(queryParamMap map[string]string) map[string]string {

	if queryParamMap == nil {
		queryParamMap = make(map[string]string)
	}
	conf, _ := config.ReadConfigs()
	if conf.GlobalAdapter.Enabled {
		queryParamMap[organizationID] = commonOrganizationIDValue
	}
	return queryParamMap
}

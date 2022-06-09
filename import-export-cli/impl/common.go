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

package impl

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Jeffail/gabs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExecuteNewFileUploadRequest forms an HTTP request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func ExecuteNewFileUploadRequest(uri string, params map[string]string, paramName, path,
	accessToken string, isOAuthToken bool) (*resty.Response, error) {

	headers := make(map[string]string)
	if isOAuthToken {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	} else {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	}
	headers[utils.HeaderAccept] = "application/json"
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive
	return utils.InvokePOSTRequestWithFileAndQueryParams(params, uri, headers, paramName, path)
}

// From the template data (tmpl) writes the target file using the provided mainConfig
func WriteTargetFileFromTemplate(targetFile string, tmpl []byte, envs *utils.MainConfig) error {
	t, err := template.New("").Parse(string(tmpl))
	if err != nil {
		return err
	}

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, envs.Environments)
	if err != nil {
		return err
	}
	return nil
}

// Include x_meta.yaml (api_meta.yaml, application_meta.yaml,api_product_params.yaml ) into the sourceZipFile and create
// a new Zip file in the provided targetZipFile location. metaFile needs to be one of the supported x_meta.yaml.
//  api_meta.yaml, application_meta.yaml,api_product_params.yaml
func IncludeMetaFileToZip(sourceZipFile, targetZipFile, metaFile string, metaData utils.MetaData) error {
	//	Create a temp directory (tmpClonedLoc) by extracting the original zip file.
	tmpClonedLoc, err := utils.GetTempCloneFromDirOrZip(sourceZipFile)
	// Create the *_meta.yaml file inside the cloned directory.
	tmpLocationForAPIMetaFile := filepath.Join(tmpClonedLoc, metaFile)
	marshaledData, err := jsoniter.Marshal(metaData)
	if err != nil {
		return err
	}

	jsonMetaData, err := gabs.ParseJSON(marshaledData)
	metaContent, err := utils.JsonToYaml(jsonMetaData.Bytes())
	if err != nil {
		return err
	}

	//write the meta content into *_meta.yaml files
	err = ioutil.WriteFile(tmpLocationForAPIMetaFile, metaContent, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error creating api_meta.yaml inside the exported zip archive", err)
	}

	err = utils.Zip(tmpClonedLoc, targetZipFile)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}
	return nil
}

//Load the x_meta.yaml file in the provided path and return
func LoadMetaInfoFromFile(path string) (*utils.MetaData, error) {
	fileContent, err := GetFileContent(path)
	if err != nil {
		return nil, err
	}
	metaInfo := &utils.MetaData{}
	err = yaml.Unmarshal([]byte(fileContent), &metaInfo)
	if err != nil {
		return nil, err
	}
	return metaInfo, err
}

//Read the file content from the provided path
func GetFileContent(path string) (string, error) {
	r, err := os.Open(path)
	defer func() {
		_ = r.Close()
	}()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//Find the path for a file matching the provided pattern. If that file is not found in the root directory, an empty string
//will be returned
func GetFileLocationFromPattern(root, pattern string) (string, error) {
	var match string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			match, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			return io.EOF
		}
		return nil
	})
	return match, err
}

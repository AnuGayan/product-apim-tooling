package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// Match for $VAR or ${VAR} and capture VAR inside a group
var re = regexp.MustCompile(`\${?(\w+)}?`)

// Match for ${VAR} and capture VAR inside a group
var recb = regexp.MustCompile(`\${(\w+)}`)

// ErrRequiredEnvKeyMissing represents error used for indicate environment key missing
type ErrRequiredEnvKeyMissing struct {
	// Key is the missing entity
	Key string
}

func (e ErrRequiredEnvKeyMissing) Error() string {
	return fmt.Sprintf("%s is required, please set the environment variable", e.Key)
}

// EnvSubstituteForCurlyBraces substitutes variables from environment to the content.
// It uses regex to match in ${var} format for variables and look up them in the environment before processing.
// returns an error if anything happen
func EnvSubstituteForCurlyBraces(content string) (string, error) {
	var errorResults error
	missingEnvKeys := false
	matches := recb.FindAllStringSubmatch(content, -1) // matches is [][]string

	for _, match := range matches {
		Logln(LogPrefixInfo+"Looking for:", match[0])
		if os.Getenv(match[1]) == "" {
			missingEnvKeys = true
			errorResults = multierror.Append(errorResults, &ErrRequiredEnvKeyMissing{Key: match[0]})
		} else {
			content = strings.ReplaceAll(content, match[0], os.Getenv(match[1]))
		}
	}

	if missingEnvKeys {
		return "", errorResults
	}

	return content, nil
}

// Substitutes all the environment variables added in the file specified in the 'file' input and changes are
// updated in the file.
// If any required environment variable is not set will throw an error.
func EnvSubstituteInFile(file string, fileExtensions []string) error {
	var substitutedContent string
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// If the fileExtensions are nil, env variables will be substituted irrespective of the extension type
	if fileExtensions == nil {
		substitutedContent, err = EnvSubstituteForCurlyBraces(string(content))
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(file, []byte(substitutedContent), 0644)
		if err != nil {
			return err
		}
	} else {
		// If the fileExtensions is not nil, env variables will be substituted to the files that will only
		//  be ending with the set of extensions passed as fileExtensions
		for _, extension := range fileExtensions {
			if strings.HasSuffix(file, extension) {
				Logln(LogPrefixInfo+"Substituting env variables to restricted extensions: ", extension)
				substitutedContent, err = EnvSubstituteForCurlyBraces(string(content))
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(file, []byte(substitutedContent), 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Walks through all the files in the given folder and substitutes all the environment variables added in
// those files. The files will be updated with the substituted values.
// If any required environment variable is not set will throw an error.
func EnvSubstituteInFolder(folderPath string, fileExtensions []string) error {
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fi, err := os.Stat(path)
			if err != nil {
				return err
			}
			if fi.Mode().IsRegular() {
				Logln(LogPrefixInfo+"Substituting env variables in: ", path)
				err = EnvSubstituteInFile(path, fileExtensions)
				if err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

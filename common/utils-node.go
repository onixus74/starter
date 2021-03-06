package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	//"strconv"
	"encoding/json"
	"github.com/blang/semver"
	"io/ioutil"
	"regexp"
	"strconv"
)

// Looks for node version in the package.json. If found returns true, version if not false, ""
func GetNodeVersion(packageJsonFile string) (bool, []string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, []string{GetDefaultNodeVersion()}
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, []string{GetDefaultNodeVersion()}
	}

	if data["engines"] == nil {
		return false, []string{GetDefaultNodeVersion()}
	}

	if nodeVersionRanges, ok := data["engines"].(map[string]interface{})["node"].(string); ok {
		versions := []string{}
		version_ranges := strings.Split(nodeVersionRanges, "||")
		for _, version_range := range version_ranges {
			//remove spaces
			version_range = strings.Trim(version_range, " ")

			//remove v character
			version_range = strings.Replace(version_range, "v", "", 2)

			//change x to 0
			version_range = strings.Replace(version_range, "x", "0", 2)

			//pad version number
			version_range = PadVersionNumber(version_range)

			//tilda support
			if strings.Index(version_range, "~") == 0 {
				tilda_start_version, err := semver.Parse(strings.TrimLeft(version_range, "~"))
				if err == nil {
					version_range = ">=" + tilda_start_version.String() + " <" + strconv.FormatUint(tilda_start_version.Major, 10) + "." + strconv.FormatUint(tilda_start_version.Minor+1, 10) + ".0"
				}
			}

			//caret support
			if strings.Index(version_range, "^") == 0 {
				caret_start_version, err := semver.Parse(strings.TrimLeft(version_range, "^"))
				if err == nil {
					version_range = ">=" + caret_start_version.String() + " <" + strconv.FormatUint(caret_start_version.Major+1, 10) + ".0.0"
				}
			}

			allowed_versions := GetAllowedNodeVersions()

			for _, allowed_version := range allowed_versions {
				//pad version number
				allowed_version_padded := PadVersionNumber(allowed_version)

				version_semver, err := semver.Parse(allowed_version_padded)
				if err == nil {

					range_semver, err := semver.ParseRange(version_range)
					if err == nil {
						if range_semver(version_semver) {
							versions = append(versions, allowed_version)
						}
					}
				}
			}
		}
		if len(versions) == 0 {
			return false, []string{GetDefaultNodeVersion()}
		} else {
			return true, versions
		}
	}
	return false, []string{GetDefaultNodeVersion()}
}

func GetMeteorVersion(meteorReleaseFile string) (bool, string) {
	file, err := os.Open(meteorReleaseFile)
	if err != nil {
		return false, "version not detected"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var re = regexp.MustCompile(`METEOR@(.*)`)
	for scanner.Scan() {
		var version = scanner.Text()
		if re.MatchString(version) {
			return true, re.FindStringSubmatch(version)[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return false, "version not detected"
	}

	return false, "version not detected"
}

func GetClosedAllowedNodeVersion(major uint64, minor uint64, patch uint64) string {
	for _, version := range allowedNodeVersions {
		if strings.Index(version, fmt.Sprintf("%d.%d", major, minor)) == 0 {
			return version
		}
	}
	for _, version := range allowedNodeVersions {
		if strings.Index(version, fmt.Sprintf("%d", major)) == 0 {
			return version
		}
	}
	//last resort
	return "latest"
}

func GetNodeDatabase(packageJsonFile string, databaseNames ...string) (bool, string) {
	found, name := GetDependencyVersion(packageJsonFile, databaseNames...)
	return found, name
}

func GetDependencyVersion(packageJsonFile string, dependencyNames ...string) (bool, string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	if data["dependencies"] != nil {
		for dependency, version := range data["dependencies"].(map[string]interface{}) {
			for _, dependencyName := range dependencyNames {
				found := dependencyName == dependency
				if found {
					return true, version.(string)
				}
			}

		}
	}

	if data["optionalDependencies"] != nil {
		for dependency, version := range data["optionalDependencies"].(map[string]interface{}) {
			for _, dependencyName := range dependencyNames {
				found := dependencyName == dependency
				if found {
					return true, version.(string)
				}
			}

		}
	}

	return false, ""
}

func GetScriptsStart(packageJsonFile string) (bool, string) {
	buf, err := ioutil.ReadFile(packageJsonFile)
	if err != nil {
		return false, err.Error()
	}

	var data map[string](interface{})
	err = json.Unmarshal(buf, &data)

	if err != nil {
		return false, err.Error()
	}

	if data["scripts"] == nil {
		return false, ""
	}

	if _, ok := data["scripts"].(map[string]interface{})["start"].(string); ok {
		return true, "npm start"
	} else {
		return false, ""
	}
}

func PadVersionNumber(version string) string {
	if ok, _ := regexp.MatchString(`^\D{0,2}\d$`, version); ok {
		version = version + ".0.0"
	} else if ok, _ := regexp.MatchString(`^\D{0,2}\d+\.\d+$`, version); ok {
		version = version + ".0"
	}
	return version
}

func SetAllowedNodeVersions(versions []string) {
	allowedNodeVersions = versions
}

func GetAllowedNodeVersions() []string {
	return allowedNodeVersions
}

func GetDefaultNodeVersion() string {
	return defaultNodeVersion
}

func GetSupportedNodeFrameworks() []string {
	return []string{"meteor-node-stubs", "keystone", "express", "loopback", "restify", "actionhero", "hapi", "socket.io", "koa"}
}

var defaultNodeVersion = "4.6"
var allowedNodeVersions = []string{"4.6"}

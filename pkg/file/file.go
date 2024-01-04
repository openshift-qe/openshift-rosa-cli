package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/ini.v1"
)

// IniConnection builds the connection of the ini file
func IniConnection(filename string) *ini.File {
	iniCfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file %s: %v", filename, err)
		os.Exit(1)
	}

	return iniCfg
}

// IsRecordExist checks whether the section exists in the ini file
func IsRecordExist(clusterID string, iniCfg *ini.File) (*ini.Section, bool) {
	sec, err := iniCfg.GetSection(clusterID)
	return sec, err == nil
}

func WriteToFile(content string, fileName string, path ...string) (string, error) {
	KeyPath, _ := os.UserHomeDir()

	if len(path) != 0 {
		KeyPath = path[0]
	}

	filePath := fmt.Sprintf("%s/%s", KeyPath, fileName)
	if IfFileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Delete file err:%v", err)
			return "", err
		}
	}
	err := ioutil.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		fmt.Println("Write to file err:%v", err)
		return "", err
	}
	return filePath, nil
}

func IfFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

package utils

import (
	"encoding/json"
	"io"
	"os"
)

// WriteToJSONFile reads content stored in io.Reader, then write the content into json file with given fileName
func WriteToJSONFile(r io.Reader, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	io.Copy(file, r)
	return nil
}

// ReadFromJSONFile reads json file with given fileName, then save into object
func ReadFromJSONFile(fileName string, v interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

//ConvertJSONReaderToObject convert json reader to object
func ConvertJSONReaderToObject(data io.Reader, object interface{}) error {
	return json.NewDecoder(data).Decode(&object)
}

//ConvertObjectToJSON converts onject to json string
func ConvertObjectToJSON(object interface{}) string {
	b, _ := json.Marshal(object)
	return string(b)
}

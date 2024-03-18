package xml

import (
	"encoding/xml"
	"github.com/ydx1011/gobatis/logging"
	"io"
	"io/ioutil"
	"os"
)

const (
	MapperStart = "mapper"
)

func ParseFile(path string) (*Mapper, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}

	return Parse(data)
}

func Parse(data []byte) (*Mapper, error) {
	v := Mapper{}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}
	return &v, nil
}

func parseInner(r io.Reader) {
	decoder := xml.NewDecoder(r)
	var strName string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		name := getStartElementName(token)
		if name != nil {
			if name.Local == MapperStart {
				switch t := token.(type) {
				case xml.StartElement:
					stelm := xml.StartElement(t)
					logging.Debug("start: ", stelm.Name.Local)
					strName = stelm.Name.Local
				case xml.EndElement:
					endelm := xml.EndElement(t)
					logging.Debug("end: ", endelm.Name.Local)
				case xml.CharData:
					data := xml.CharData(t)
					str := string(data)
					switch strName {
					case "City":
						logging.Debug("city:", str)
					case "first":
						logging.Debug("first: ", str)
					}
				}
				break
			}
		}
	}
}

func getStartElementName(token xml.Token) *xml.Name {
	switch t := token.(type) {
	case xml.StartElement:
		stelm := xml.StartElement(t)
		logging.Debug("start: ", stelm.Name.Local)
		return &stelm.Name
	}
	return nil
}

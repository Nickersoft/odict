package odict

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type DictionaryUsage struct {
	POS              string                      `json:"pos" xml:"pos,attr,omitempty"`
	Definitions      []string                    `json:"definitions" xml:"definition"`
	DefinitionGroups []DictionaryDefinitionGroup `json:"groups" xml:"group"`
	XMLName          xml.Name                    `json:"-" xml:"usage"`
}

type DictionaryUsageMap struct {
	Iterable map[string]DictionaryUsage
}

func (m *DictionaryUsageMap) Set(key string, value DictionaryUsage) {
	m.Iterable[key] = value
}

func (m *DictionaryUsageMap) Get(key string) DictionaryUsage {
	return m.Iterable[key]
}

func (m DictionaryUsageMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Iterable)
}

func (m DictionaryUsageMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	for key := range m.Iterable {
		e.Encode(m.Get(key))
	}
	return nil
}

func (m *DictionaryUsageMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var usage DictionaryUsage

	d.DecodeElement(&usage, &start)

	if m.Iterable == nil {
		m.Iterable = make(map[string]DictionaryUsage)
	}

	m.Set(usage.POS, usage)

	for {
		_, err := d.Token()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

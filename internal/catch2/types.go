package catch2

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) parse(s string) error {
	n, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	if err != nil {
		return err
	}

	i, accuracy := n.Int64()
	if accuracy != big.Exact {
		return fmt.Errorf("value too wide: %s", s)
	}

	d.Duration = time.Duration(i)

	return nil
}

func (d *Duration) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string
	err := decoder.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	return d.parse(s)
}

func (d *Duration) UnmarshalXMLAttr(attr xml.Attr) error {
	return d.parse(attr.Value)
}

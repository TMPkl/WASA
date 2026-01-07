package attachments

import (
	"bytes"
	"encoding/gob"
)

// przygotowanie do zapisu i odczyu z DB

func (ap *AttachmentsPack) ConvertToGOB() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(ap)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeFromGOB(data []byte) (*AttachmentsPack, error) {
	buf := bytes.NewBuffer(data)
	var ap AttachmentsPack
	err := gob.NewDecoder(buf).Decode(&ap)
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

package decoder

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(data []byte) ([]byte, error) {
	return data, nil
}

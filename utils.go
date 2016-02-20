package bayes

import (
	"bytes"
	"encoding/gob"
)

func NumberPositionalBinning(data string) []string {
	// TODO: TEST
	// it does this: Useful for classification when you cannot extrapolate data
	// 1940 = 1*** 19** 194* 1945
	// 1945 = 1*** 19** 194* 1945
	// 1980 = 1*** 19** 198* 1980
	// 2000 = 2*** 20** 200* 2000
	// 2014 = 2*** 20** 201* 2014
	// 2013 = 2000 2010 2013
	datab := []byte(data)
	result := make([][]byte, 0)
	for _, fbyte := range datab {
		if len(result)-1 >= 0 {
			result = append(result, append(result[len(result)-1], fbyte))
		} else {
			last_one := make([]byte, 0)
			last_one = append(last_one, fbyte)
			result = append(result, last_one)
		}
	}
	for k, v := range result {
		d := make([]byte, len(v))
		copy(d, v)
		for i := 0; i <= (len(datab)-len(v))-1; i++ {
			d = append(d, []byte("*")...)
		}
		result[k] = d
	}
	return_data := make([]string, 0)
	for _, v := range result {
		return_data = append(return_data, string(v))
	}
	return return_data
}

func ByteToAny(bytes_data []byte, feedf func(*gob.Decoder) error) error {
	return feedf(gob.NewDecoder(bytes.NewReader(bytes_data)))
}

func AnyToByte(any interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(any)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type ClassifyResultSlice []ClassifyResult

func (slice ClassifyResultSlice) Len() int {
	return len(slice)
}

func (slice ClassifyResultSlice) Less(i, j int) bool {
	return slice[i].Log > slice[j].Log
}

func (slice ClassifyResultSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// tooler_test
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func test_Bolt(t *testing.T) {
	fileDb := "testBolt.db"
	bucket := "testBucket"
	// get1kv := func() ([]byte, []byte) {
	// 	return nil, nil
	// }
	getkv := func() ([]byte, []byte) {
		return nil, nil
	}
	// setkv := func(k, v []byte) {

	// }
	boltBatchWriteKeyValue(fileDb, bucket, getkv)

	assert.Equal(t, true, true, "test")
}

func test_Sscan(t *testing.T) {
	s := "430354,1555925,新疆,伊犁,0999,中国联通"
	var idx, number, province, area, code, isp string
	fmt.Sscanf(s, "%s,%s,%s,%s,%s,%s", &idx, &number, &province, &area, &code, &isp)
	assert.Equal(t, idx, "430354", "")
	assert.Equal(t, number, "1555925", "")
	assert.Equal(t, province, "新疆", "")
	assert.Equal(t, area, "伊犁", "")
	assert.Equal(t, code, "0999", "")
	assert.Equal(t, isp, "中国联通", "")
}
func test_Scan(t *testing.T) {
	cdr := "010TK6156567817600748685,2,3568,17092395243,13996370001,,,20190813143540,20190813143604,20190813143616,755d525a42012580,0101300044,1041-PSTN,0,x"
	scanner := bufio.NewScanner(strings.NewReader(cdr))
	//scanner.Split(bufio.ScanLines)
	// split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 	n := bytes.IndexByte(data, ',')
	// 	if n >= 0 {
	// 		return n + 1, data[:n+1], nil
	// 	} else if len(data) > 0 {
	// 		return len(data), data[:len(data)], nil
	// 	}
	// 	if atEOF {
	// 		return 0, nil, io.EOF
	// 	}
	// 	return 0, nil, nil
	// }

	log.Println(bytes.Count([]byte(cdr), []byte{','}))

	totel := 0
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if totel < 14 {
			totel = bytes.Count(data, []byte{','})
		}

		if totel >= 14 {
			c := 14
			l := 0
			p := data
			for c > 0 {
				n := bytes.IndexByte(p, ',') + 1
				p = p[n:]
				l += n
				c--
				totel--
			}
			return l, data[:l], nil
		} else if atEOF {
			return 0, nil, io.EOF
		}
		return 0, nil, io.EOF
	}

	scanner.Split(split)
	for scanner.Scan() {
		t.Log(scanner.Text())
	}
}

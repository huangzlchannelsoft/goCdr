// phone
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	//"github.com/plandem/xlsx"
)

const (
	PhonePropertyDB   = "phoneProperty.db"
	PhonePropertyBulk = "phone2property"
)

type PhoneProperty struct {
	productor string
	isp       string
	province  string
	area      string
}

var (
	phone2Property  map[string]*PhoneProperty
	phoneHttpClient = http.Client{Timeout: time.Second * 5}
	_phoneIspUri    = ""
	_phoneProUri    = ""
)

func init() {
	log.Println("init phone!")

	LoadPhoneProperty()
}

func SetPhonePropertyUri(ispUri string, proUri string) {
	_phoneIspUri = ispUri
	_phoneProUri = proUri
}

func LoadPhoneProperty() {
	phone2Property = make(map[string]*PhoneProperty)

	kv := func(k, v []byte) {
		// log.Printf("key=%s, value=%s\n", k, v)

		p := strings.Split(string(v), "_")
		phone2Property[string(k)] = &PhoneProperty{
			productor: p[0],
			isp:       p[1],
			province:  p[2],
			area:      p[3],
		}
	}
	boltEnumKeyValue(PhonePropertyDB, PhonePropertyBulk, kv)
}

// func LoadRemotePhoneProperty(phoneListFile string) {
// 	phoneList, err := ioutil.ReadFile(phoneListFile)
// 	if err != nil {
// 		log.Panic(err.Error())
// 	}

// 	kv := func() ([]byte, []byte) {
// 		n, number, err := bufio.ScanWords(phoneList, true)
// 		if n == 0 || err != nil {
// 			return nil, nil
// 		}

// 		//TODO: get phone property by number
// 		productor := ""
// 		isp := ""
// 		province := ""
// 		area := ""
// 		return number, []byte(PhoneProperty2Key(productor, isp, province, area))
// 	}

// 	boltBatchWriteKeyValue(PhonePropertyDB, PhonePropertyBulk, kv)

// }

/*
excel表结构 ： 话批 客户名称 省份 地市 电话号码 号码开通时间	号码类型 号码状态
*/
// func LoadExcelPhoneProperty(filepath string) {
// 	xl, err := xlsx.Open(filepath)
// 	if err != nil {
// 		log.Println("[Err]", err.Error())
// 		return
// 	}
// 	defer xl.Close()

// 	sht := xl.Sheet(0, 1)
// 	_, r := sht.Dimension()
// 	for j := 1; j < r; j++ {
// 		// log.Println(sht.Cell(0, j).Value(),
// 		// 	sht.Cell(2, j).Value(),
// 		// 	sht.Cell(3, j).Value(),
// 		// 	sht.Cell(4, j).Value())
// 		productor := sht.Cell(0, j).Value()
// 		isp := ""
// 		province := sht.Cell(1, j).Value()
// 		area := sht.Cell(2, j).Value()
// 		number := sht.Cell(3, j).Value()

// 		kv := func() ([]byte, []byte) {
// 			return []byte(number), []byte(PhoneProperty2Key(productor, isp, province, area))
// 		}
// 		boltWriteKeyValue(PhonePropertyDB, PhonePropertyBulk, kv)
// 	}
// }

/*
return productor_isp_province_area
*/
func GetPhoneProperty(na []string, nb []string) []string {
	requestPhoneProperty_(na, _phoneProUri)
	requestPhoneProperty_(nb, _phoneIspUri)

	if len(na) != len(nb) {
		log.Printf("[Err] GetPhoneProperty.%d!=%d\n", len(na), len(nb))
	} else {
		var keys []string

		productor := ""
		isp := ""
		province := ""
		area := ""
		l := len(na)
		for i := 0; i < l; i++ {
			pa := phone2Property[na[i]]
			pb := phone2Property[nb[i]]

			if pa == nil {
				productor = ""
			} else {
				productor = pa.productor
			}

			if pb == nil {
				isp = ""
				province = ""
				area = ""
			} else {
				isp = pb.isp
				province = pb.province
				area = pb.area
			}

			keys = append(keys, PhoneProperty2Key(productor, isp, province, area))
		}
		return keys
	}

	return nil
}

/*http request phone-property
request:
{"numbers":["17092395243","01053189237"]}
isp-json:
{"info": "success", "data": [{"serv_provider": "fixed", "phone": "01053189237", "flag": "fixed", "city_name": "北京", "city_code": "010", "province_name": "北京"}], "result": true}
{"info": "success", "data": [{"serv_provider": "中国联通", "phone": "17092395243", "flag": "mobile", "city_name": "重庆", "city_code": "023", "province_name": "重庆"}], "result": true}
pro-json:
{"info": "success", "data": [{"city_name": "", "phone": "01053189237", "province_name": "", "custom_Name": "西安初见数据网络科技有限公司 ", "voip_Name": "讯众"}], "result": true}
*/
type PhonePropertyJson struct {
	Phone     string `json:"phone"`
	Province  string `json:"province_name"`
	Area      string `json:"city_name"`
	Productor string `json:"voip_Name"`
	Isp       string `json:"serv_provider"`
}
type PhoneRequestJson struct {
	Info     string               `json:"info"`
	Property []*PhonePropertyJson `json:"data"`
}

//parse isp-json || pro-json
//write phoneProperty to phone2Property
//save phoneProperty boltdb
func PhonePropertyParser(data []byte) {
	log.Println(string(data))

	var pp PhoneRequestJson
	err := json.Unmarshal(data, &pp)
	if err != nil {
		log.Println("[Err] PhonePropertyParser", err.Error())
		return
	} else {

		if pp.Info == "success" {

			l := len(pp.Property)
			cursor := 0
			kv := func() ([]byte, []byte) {
				log.Println(cursor, l)
				if cursor == l {
					return nil, nil
				}
				number := pp.Property[cursor].Phone
				productor := pp.Property[cursor].Productor
				isp := pp.Property[cursor].Isp
				province := pp.Property[cursor].Province
				area := pp.Property[cursor].Area

				phone2Property[number] = &PhoneProperty{
					productor: productor,
					isp:       isp,
					province:  province,
					area:      area,
				}
				cursor++

				log.Println(number, productor, isp, province, area)

				return []byte(number), []byte(PhoneProperty2Key(productor, isp, province, area))
			}

			boltBatchWriteKeyValue("phoneProperty.db", "phone2property", kv)
		} else {
			log.Println("[Err] Request PhoneProperty", string(data))
		}
	}
}

func requestPhoneProperty_(na []string, uri string) {
	//{"numbers":["17092395243","01053189237"]}

	var buf bytes.Buffer
	count := 0
	buf.WriteString("{\"numbers\":[")
	for _, a := range na {
		pp := phone2Property[a]
		if pp == nil {
			if count > 0 {
				buf.WriteString(",")
			}
			buf.WriteString("\"")
			buf.WriteString(a)
			buf.WriteString("\"")
			count++
		}
	}
	buf.WriteString("]}")

	if count > 0 {
		HttpPost(&phoneHttpClient, uri, buf.String(), PhonePropertyParser)
	}
}

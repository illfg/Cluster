package conf

import (
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var conf map[string]string

const confLogFlag = "conf"

func ParseConf() {
	file := getConfFile()
	conf = make(map[string]string, 8)
	err := yaml.Unmarshal(file, conf)
	if err != nil {
		glog.Fatalf("[%s]:fail to parse conf file, err is %s\n", confLogFlag, err)
	}
}

func ParseYaml(in string) map[string]string {
	yMap := make(map[string]string, 8)
	err := yaml.Unmarshal([]byte(in), yMap)
	if err != nil {
		glog.Warningf("[%s]:fail to parse yaml from bytes, err is %s\n", confLogFlag, err)
		return nil
	}
	return yMap
}

func GetConf() *map[string]string {
	if conf == nil {
		ParseConf()
	}
	return &conf
}

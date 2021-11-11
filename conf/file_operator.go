package conf

import (
	"github.com/golang/glog"
	"io/ioutil"
)

const fileOperatorLogFlag = "file operator"

func getConfFile() []byte {
	data, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		return nil
	}
	if err != nil {
		glog.Fatalf("[%s]:fail to open conf file, err is %s\n", fileOperatorLogFlag, err)
	}
	return data
}

package utils

import (
	"github.com/golang/glog"
	"os"
)

type mappingFile struct {
	path    string
	file    *os.File
	currIdx int64 //当前读取文件下标
}

const mappingFileLogFlag = "mappingFile"

func newMappingFile(path string) *mappingFile {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		glog.Errorf("[%s]: can't open specified file, err is %s\n", mappingFileLogFlag, err)
		return nil
	}
	//fileInfo, err := file.Stat()
	//if err != nil {
	//	glog.Errorf("[%s]: can't read FileInfo of specified file, err is %s\n", mappingFileLogFlag, err)
	//	return nil
	//}
	//data := make([]byte, fileInfo.Size())
	//_, err = file.Read(data)
	//if err != nil {
	//	glog.Errorf("[%s]: can't read, err is %s\n", mappingFileLogFlag, err)
	//	return nil
	//}
	mf := &mappingFile{
		path:    path,
		file:    file,
		currIdx: 0,
	}
	return mf
}

func (receiver *mappingFile) getCheckPoint() int64 {
	return receiver.currIdx
}

func (receiver *mappingFile) readByOffset(b []byte, offset int64) (int, error) {
	return receiver.file.ReadAt(b, offset)
}

func (receiver mappingFile) next() {

}

// Code generated by go-bindata.
// sources:
// data/aws-simple/build/build-ruby.sh.tpl
// data/aws-simple/build/template.json.tpl
// data/aws-simple/deploy/main.tf.tpl
// data/aws-vpc-public-private/build/build-ruby.sh.tpl
// data/aws-vpc-public-private/build/template.json.tpl
// data/aws-vpc-public-private/deploy/main.tf.tpl
// data/common/dev/Vagrantfile.tpl
// DO NOT EDIT!

package rubyapp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dataAwsSimpleBuildBuildRubyShTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x56\x7f\x53\x1b\x39\x12\xfd\x7f\x3e\x45\xaf\xa1\xc2\x5d\x15\x9a\x01\xee\xf6\x6e\xd7\x5e\xb6\x0e\x58\x48\xa8\x4a\x99\x94\x49\xee\x47\x25\x39\x6a\x3c\xd3\xb6\x15\xc6\xd2\x44\xd2\x18\x0c\xf1\x77\xbf\xd7\x1a\xe3\x1f\x21\x97\xca\x3f\x30\x96\xd4\xdd\xaf\xbb\x5f\x3f\x69\xe7\xa7\x6c\xa8\x4d\x36\xcc\xfd\x24\x49\x3c\x07\x52\x96\x8c\x6d\xcc\xf2\x93\x9d\xe3\x7b\x1d\x3f\x6b\x5d\xf3\x28\xd7\xd5\x72\x39\xb8\xbc\xe0\x24\xc1\x97\x75\x7f\xfa\x33\x3d\x26\x44\x54\xd9\x22\xaf\xc8\xdb\xc6\x15\x3c\xd2\x15\x1f\xef\x1e\xae\x97\x2b\x6d\xd8\xd8\xe3\xdd\x23\x59\xe2\x62\x62\xa9\x73\x3e\x18\x5c\x0d\x28\x0f\xb4\xfb\xb8\x36\x5a\x74\x77\x1f\xdb\xb3\x8b\x1e\xbd\xce\x7d\x80\xfd\xd8\x77\x3b\x62\x36\x76\x5c\x93\x0d\xc1\x52\x36\xcb\x5d\x86\x8d\xcc\xcf\x3d\xfe\xd1\x17\x0a\x11\x9b\xa1\xa3\x83\x64\x91\x00\x5d\x4d\x7b\x11\x1c\x75\x76\x1f\x4f\x4f\xae\x5f\xdd\x5c\x5f\xbd\x1b\x9c\x9d\x2f\x3a\xb2\xf0\xfa\xb2\x7f\xde\xbf\x5a\x74\xf6\x08\x18\x92\xc4\xb2\xa4\x80\x8d\x7f\x74\xe8\xe8\xf7\x17\x87\x70\x07\xa7\x63\x76\xa4\x42\x1b\xef\x77\xca\x4a\x9e\x65\xa6\xa9\xaa\x1e\x2d\x12\x5b\x45\x83\x36\x8d\xf7\x72\xe2\x23\xc1\x58\xb6\x92\x1d\x2a\x2a\xdb\x94\xaa\xb0\x66\xa4\xc7\x54\xe4\x86\xb4\x09\xec\x46\xec\x98\xee\x74\x98\x50\x5e\x07\x2a\xec\x74\x9a\x9b\xd2\x93\x1e\x91\x0e\x7b\x9e\x7c\xd0\x55\x85\x93\x54\x3b\x8b\x3c\xbd\x47\x10\xea\xfc\x2b\xd7\x41\x9b\x31\x8d\x90\xc8\x96\x5b\x60\x82\x8b\xba\xe2\xc0\x69\x9a\x76\x92\xc6\xc0\x9e\xde\xbf\x27\x35\x5a\x16\x47\x0f\xb3\x68\x91\x69\xe3\x43\x6e\x0a\xce\x86\xd6\x06\x35\xd2\x46\xfb\x09\x97\xf4\xf1\x63\x8f\x4a\x8b\xb2\xfa\x8a\x51\xd6\x83\xf4\xe7\xa4\xb4\x06\x3d\x95\xb8\x27\x65\x29\x61\x05\x29\x6a\x6e\xbd\x0e\xd6\x69\xf6\x04\xc8\xd4\xd4\x65\x2e\xa0\x62\x5c\xcb\xe4\x9b\xd2\xca\x49\x35\x06\x69\xe2\x26\x3f\x5b\x8e\x18\x90\x9f\x9a\x53\x3d\x0f\x13\x6b\x94\xb7\xa3\x70\x97\x3b\x56\xc8\xb7\x66\x17\xc4\xfb\x37\xd6\x94\x14\xca\x9a\xe8\x08\x5d\x35\xbe\xb6\x2e\xa8\x49\x08\xb5\x5f\x07\x29\x4b\x25\xfb\x2b\xa4\xf3\x18\xa7\xce\xbb\xc5\xc4\x69\xaf\x2a\xce\x33\x63\x4b\x4e\x3f\xf9\x2d\x60\x62\xf7\xdc\x66\xe8\xf4\x78\x12\x86\xf6\x3e\x73\xcd\x70\xae\xcc\x78\xcb\xe6\x96\xe7\x88\x37\x23\x25\x5f\x9e\xdd\x0c\x24\x99\xdc\xd6\xdd\x2c\x5b\xfd\x4e\x9b\x21\xba\xd1\xa4\x40\xde\xfd\xe5\x00\x27\x1d\x17\xb3\x78\x9c\x7e\xfe\xdb\xe1\xc5\xaf\xa7\xbf\x9e\x9d\x9c\xfd\xf5\xe0\xf4\xe8\xe2\xef\x49\x64\xd0\x5e\xc9\x43\x8a\x29\xc1\x8d\xf5\x5e\x61\x24\x73\x29\x77\x5a\x4f\x1a\xaf\xad\xa9\x73\xef\xd9\x80\x8f\xe2\x33\x03\x8c\x6c\xb5\x42\xc1\x35\x3e\xcc\x69\x9a\x6b\xb3\x07\xde\x46\xa0\x81\x99\x32\x0e\x45\x3c\xda\x8e\x96\x4f\x2b\xed\x43\x5a\xae\x2d\xe3\xc2\x26\xb1\xff\x5f\x2f\x13\xbe\x97\xa2\xd3\xe0\xdd\xe9\x7f\x6e\xfe\x79\x3e\xb8\xbe\xbc\xea\x1f\x77\x1e\x1f\x49\xea\x73\x83\x84\x05\x22\x2d\x16\x9d\x96\x3a\x97\x6d\xab\x85\x3e\x03\x1c\xd8\xa7\x37\x4f\x11\xf7\xa9\x3f\xd6\xe6\x7e\x3f\xb2\xc8\x86\x09\xd0\xd7\x79\x71\x9b\x8f\x81\x4e\xb8\xb4\x8c\xf3\xc7\xf9\xe9\xe5\x49\xff\xe6\x62\x70\xd5\x7f\x7b\xde\xff\xe3\xd8\x58\x13\x07\x28\x2f\x82\x9e\x7d\x97\x5a\xc3\x07\x47\x63\xc8\xd5\x94\x5d\xd1\x38\x0d\xd5\x19\x36\xba\x2a\x15\x0b\x80\x20\xbf\x3f\x80\xef\x98\x8c\xfa\xb3\x42\xd6\xf4\x80\xcf\xc3\x71\xfc\xfc\x0e\xf5\xc4\x46\xe8\xf3\xc9\x3f\x99\xfb\xcf\x95\x0e\xfc\x97\x68\x28\x4b\x52\x88\xdd\xcd\xf2\x3c\x5f\x59\x9d\x35\x52\x02\xc5\xf7\x20\xb3\xa7\x55\x33\x9e\x95\xee\xb4\x31\x65\x85\x26\x6d\xce\xd8\x98\xa7\xab\x6c\x87\xed\x3e\xd8\x65\xac\x72\x7a\xf9\xbf\xb4\x45\xeb\xe9\x5c\xfc\x17\xa1\x9d\xe1\x3a\x7a\x89\x2e\xa6\xb7\xa5\x86\x51\x4d\x99\x77\xb3\x4c\x84\x0b\x93\x53\xb7\x7b\x21\x77\xf4\x70\x0f\xf9\x08\xd3\x7a\xb5\x95\x86\xf1\x03\xa9\xb3\xaf\xce\x6f\x6b\x44\x5d\xe9\x02\x8a\x80\x52\x35\xfe\x2b\xc8\x18\x31\x59\x03\xbc\x52\xfb\x7c\x58\x71\xa9\x24\xe7\x3b\xeb\x4a\xac\x8d\xb9\xb0\x9e\x3a\x1d\xda\x76\x7c\xcd\x21\x22\x47\x1f\xa6\xda\x0b\xbb\xfc\x96\x53\xcc\xcc\x9d\x21\x35\x58\x99\x75\xbf\x05\xef\x2c\x0a\x25\x68\x00\x4f\xb1\xe8\xd1\x07\xe4\xf9\xed\x44\x43\x76\x3d\x84\xed\x73\xa3\x1d\x94\x50\xc4\x75\x63\xa0\xa4\xd0\x81\x72\xec\xe7\xde\x1a\x01\x4d\x6c\x66\xda\x59\x33\x05\x8b\xe8\x6e\x22\x42\x0e\x96\x41\xd9\xe1\x0d\x7a\x5a\x12\xdf\x73\xd1\x04\x39\xea\xc1\x8f\x5b\x4c\x5f\xe3\x5d\xbc\x59\x61\xb9\xbf\xfe\x05\x56\x56\xfb\x84\xc9\x4c\xe9\x12\x21\x2a\x2f\x34\xae\x41\x3b\x13\xaa\x39\x9c\x19\x66\x5c\x09\x40\x60\x0b\x1c\xa5\x09\x94\x48\xae\x04\x8c\x0a\xb5\xba\x9f\xd2\x49\x5d\xb3\x89\x85\x07\x04\x49\xc4\xf8\x66\x34\xd2\x85\x86\x8f\x94\xba\xea\x4b\xdb\x4c\x8f\xbc\x94\xa6\xbd\x43\x9f\xfd\x57\x40\xd0\x9b\x93\xb7\xaf\x7a\x1f\x4c\xb6\xd7\x2a\x43\xac\x48\xfb\x37\x15\xd7\xdf\xb0\xda\xa1\x2b\x14\xb4\x4b\xf2\x16\x10\x6b\xcc\xc8\x46\x99\xe4\x5e\xf3\xd0\x99\x27\xad\xfa\x8e\x6b\x24\xd6\x47\x62\x92\x97\xe3\xa9\x9d\x31\x12\xd2\xa2\xf6\x6d\x5f\xa4\xd0\xc8\x1a\x32\x45\x50\x62\x6e\x91\xb8\xe9\xa6\x33\x59\xf7\x8a\x63\x33\x4a\x08\xd6\x28\x6f\xaa\xd0\xf6\x92\x3d\xc7\xb7\x05\xee\x26\xb4\xa5\xc6\xcd\x29\x4d\xc2\x6c\xc9\xf4\xe2\xd3\x3f\x15\x70\x05\x5d\x2d\x45\xa7\xa4\x35\xc6\x7d\x0c\x54\x80\xbf\x78\x15\xa3\xef\xba\x25\x42\x09\x31\x00\x13\x3c\xa3\x47\x90\x42\x7a\xba\x7c\x27\x48\x3e\xb4\xe5\xb2\x0d\xa2\x81\xf9\x66\x19\x2f\x4d\x30\x0b\xf4\xdb\x6f\xfd\x97\x97\xfd\x7f\x9f\x5d\xf5\x2f\x9e\x89\x72\x9b\x92\xb8\xda\x92\x63\x59\xd8\x92\xe3\x1d\x7a\xc9\x86\x25\x6e\x49\xc3\x79\x6c\x46\xb2\x3a\x7e\xe3\x70\x99\xb7\xc4\x92\x9b\x5e\xf4\x26\x9b\x81\x18\x16\x3b\xf2\xbd\xbc\x39\x6e\x56\x06\x99\x3c\xc3\x42\x9c\x25\xbc\x00\x7a\x9b\x9e\x70\x7e\x4d\xd1\xf5\xfa\xc8\x31\xc7\xcd\x5e\xb2\x4a\x26\xf9\xc1\xec\xb6\x1b\xb6\x12\x93\x1f\xca\x71\x79\xab\xc6\x27\x25\xc9\x25\xc5\x86\x7e\x39\xe8\xc5\x9f\x6d\xd6\x9b\xc3\x9e\xd5\xcd\x10\xfa\xd3\x6e\xaf\xc1\x2f\x43\x93\x35\x3d\xbc\x07\x37\xf0\x8b\x34\xb4\xe2\xfa\xa4\xa5\x32\x4b\x42\x90\x2d\xa1\x54\xcd\x4a\x5d\x64\x22\x56\xaf\x64\x52\x55\x41\x9d\xa2\xdc\x06\x41\x2f\x5e\x2c\x15\x79\x7d\x1d\x41\xf1\xeb\xca\xce\xa3\x66\x28\x25\x0f\x3f\xa1\x0a\x32\xe7\xca\xd6\x71\x15\x25\x0a\xcb\x4b\x13\x81\xe5\xf1\xf5\x53\x27\xf9\x5f\x00\x00\x00\xff\xff\xbe\x9a\xb3\xa9\x92\x0b\x00\x00"

func dataAwsSimpleBuildBuildRubyShTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsSimpleBuildBuildRubyShTpl,
		"data/aws-simple/build/build-ruby.sh.tpl",
	)
}

func dataAwsSimpleBuildBuildRubyShTpl() (*asset, error) {
	bytes, err := dataAwsSimpleBuildBuildRubyShTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-simple/build/build-ruby.sh.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataAwsSimpleBuildTemplateJsonTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x54\xcd\x8e\x9b\x30\x18\xbc\xf3\x14\x16\x52\xf6\x14\x20\xed\xae\xaa\xaa\xd7\x3e\xc6\x2a\x62\x0d\x78\x83\x15\xdb\x58\xfe\xec\x54\x59\xe4\x77\xef\x67\xfe\xa9\x36\x24\xdd\x5c\x8c\x3c\xc3\xcc\x78\xf8\x9c\x36\x22\xf8\x8b\x25\x57\xb9\xa6\xe5\x99\x99\xfc\xc2\x0c\xf0\x46\xc5\xbf\x48\x7c\x48\x7f\xa6\x87\x78\x1f\xf5\x9c\x0b\x35\x9c\x16\x82\x01\x42\xfd\x6b\xb8\x49\xff\x40\x4e\xcb\x92\x01\xe4\x67\x76\x45\x44\x39\x21\xf6\x4b\x14\x58\x69\x98\xbd\x85\x1a\x76\xea\xcd\x56\x08\x08\x77\xc2\x3c\xb6\x1e\x80\x6e\xdf\x8f\x41\xb4\x69\x2e\x3c\x64\xc4\xa4\x48\x78\x1d\xde\x6a\x77\xe4\xbd\x31\xa4\xe2\x86\x70\x85\x8f\x4e\x55\xd4\x22\x2b\xc7\x1d\x48\x0b\xc7\x45\x45\x76\x7e\x24\x0f\x2b\xca\xd9\xab\x66\xe1\xb4\x50\x33\x21\xe2\xfd\x0c\x70\x25\xb8\x0a\xd0\x6b\x2c\xcf\x41\x36\xd1\x24\xb3\x52\x67\x8d\xb5\x4d\x36\x1b\x24\x6d\x1b\x9c\x45\xd3\xe8\xf4\x37\xee\x5a\x66\x88\xf7\xf1\x71\x50\xf2\xfb\xdb\x9e\xef\x5c\xb0\xa5\x25\x34\xce\x94\x1d\x82\x9a\xc1\xd2\xfb\x6c\x89\x57\x0c\x2c\x57\x9d\x6b\x20\xfd\x47\x9a\x07\xc2\x6c\x15\x50\x56\x8f\x1e\xdd\x7b\xf2\xf4\x44\x0a\x0a\x35\x49\x33\x49\xb9\x4a\xa1\xfe\xa4\x8b\x1d\x61\xaa\x0a\xdf\x6b\xeb\x93\x6c\xd4\xb3\x23\x38\xa8\x05\x86\x90\xa8\x80\x29\x1c\xe0\x39\xdf\xa6\xc1\x79\xc3\x33\xf7\x1e\x0b\xda\x23\x4d\x26\x54\xeb\xd4\x9e\x3e\xbe\x54\x18\x94\x86\x6b\x1b\xa0\x6e\xdc\x12\xe3\x8a\x6b\x38\xfe\xa8\xd5\xad\xc7\x71\x8e\x3b\xce\x30\xc3\xd3\x85\x52\x54\x76\xda\x21\xcb\x24\x3d\x39\x52\x49\x3f\xb0\x75\x56\xc0\x8c\xad\xae\xdf\xad\x62\xd6\xf7\x74\xbb\x9d\x78\x75\x65\xb7\x14\x67\xe2\x1d\xc5\xe9\x9a\x6f\xa9\xf5\xa4\x7b\xd9\xba\x11\xc8\xa9\xe4\x7d\x1f\x3c\xf9\xfe\xed\xc7\xf3\xa1\x7a\x79\x99\x39\x5c\x81\xa5\x0a\x59\x63\x6d\xe5\x73\x2a\xa8\x39\xb1\x85\x0c\xd4\x79\x70\x1e\xeb\x76\x05\x0e\xaf\x5b\x94\x2a\x79\x3e\x62\x6d\x1b\x9e\x70\xae\xff\xcd\x8e\x2b\x4e\x11\x95\xfa\xb3\xc4\xfd\x7f\xd6\x31\x8a\x7c\xf4\x37\x00\x00\xff\xff\x16\x08\x5f\x5b\x65\x05\x00\x00"

func dataAwsSimpleBuildTemplateJsonTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsSimpleBuildTemplateJsonTpl,
		"data/aws-simple/build/template.json.tpl",
	)
}

func dataAwsSimpleBuildTemplateJsonTpl() (*asset, error) {
	bytes, err := dataAwsSimpleBuildTemplateJsonTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-simple/build/template.json.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataAwsSimpleDeployMainTfTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x52\xc1\x8e\xdb\x20\x10\xbd\xe7\x2b\x46\x74\x8f\x5d\x27\xed\xb1\x52\xcf\xbd\xb5\x1f\x50\xad\x10\xc6\x24\x45\x8b\x01\xc1\x90\xca\xb2\xfc\xef\x1d\xa0\x16\xc1\xbb\xbd\x36\x39\xf9\xcd\x63\xde\xcc\x9b\xf7\x01\xbe\x29\xab\x82\x40\x35\xc1\xb8\xc0\x0f\x44\xf7\x11\x26\x07\xd6\x21\xa8\x49\x23\xcc\xc2\x26\x61\xcc\x72\x3a\xdd\x45\xd0\x62\x34\x0a\x98\xb6\xd7\x20\xb8\x9e\x18\xac\xdb\x03\x2c\x7e\x47\x2e\xa4\x54\x31\xf2\x57\xb5\xbc\x53\x8c\x4a\x06\x85\xff\x28\x06\x75\xd3\xce\x1e\x0a\x44\xe5\x56\xcc\xaa\xc0\x8f\x0f\x66\x7d\x60\x6a\x1b\x51\x58\xa9\x38\x2e\x3e\xd3\x61\x52\x57\x91\x0c\xc2\x57\x60\xf8\x79\x98\xb5\x0c\x8e\xc1\xe3\x8b\x98\x46\x4b\xd3\xf8\x34\x1a\x2d\x0f\xdd\xee\x5e\x72\xa9\xa7\xf0\x0e\xfc\x77\xed\x93\x0f\xee\xae\x27\x15\xca\xf4\x04\x9d\x00\xda\xf2\x59\xf5\x69\xa5\x87\x43\x6f\xca\xc6\x88\xd6\x6c\xe8\x69\x0d\x2f\xb4\x6a\x08\xe4\x5f\x47\xab\x38\x51\x68\x88\xa0\xa2\x4b\x41\x36\x7f\x53\xd0\xb8\xf0\x5b\x70\xc9\x33\x02\xbd\xaf\x93\x65\x0f\x6b\x9f\x75\xad\x1f\xdb\xf6\x5c\x5b\xee\xc7\x2c\x9a\x75\xc1\xa6\x57\xbf\xa9\x44\x35\x6d\x6f\x24\x17\x4b\x3f\x00\x5a\x1f\x9d\x74\xa6\x8e\xf7\xfc\xa9\x80\xd7\xe0\x66\xee\x5d\xc0\x02\x5e\x0a\x86\x6e\x47\x1a\x96\xad\xe5\xa3\x71\xf2\x35\x12\xf6\x93\x5d\x86\xf2\x3f\x5f\xd8\x0b\xd5\xb7\xac\xa6\xfe\x9b\xd8\x1b\x1b\xf7\x28\x3d\x1a\x48\x81\x83\xf6\x6b\xf7\x98\x75\xf1\xad\x4b\x5f\x2b\x77\x70\xbd\x7d\x0d\x1d\x79\xdc\xf5\xe9\xb2\x58\x88\x7b\xf2\x0f\x82\x3b\x5c\x4f\x92\xcf\xd3\x1f\x9d\x3a\xd7\x2d\x9f\xd6\xb7\x89\x18\x68\x9d\x21\x9f\xf3\x25\x3f\x46\x71\x23\x7f\xe1\x7b\x16\xe9\x82\xc1\xaa\x29\x2e\xa1\x4f\x08\x2c\x05\x53\x3d\xb8\x0b\x93\x0a\xf5\x17\xa2\xff\x72\x3e\x57\x89\x7d\xc7\xd2\xbc\x2e\xc0\x27\x1b\xb7\x73\x0e\xe8\x9f\x00\x00\x00\xff\xff\x5f\x73\x79\x4b\x5f\x04\x00\x00"

func dataAwsSimpleDeployMainTfTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsSimpleDeployMainTfTpl,
		"data/aws-simple/deploy/main.tf.tpl",
	)
}

func dataAwsSimpleDeployMainTfTpl() (*asset, error) {
	bytes, err := dataAwsSimpleDeployMainTfTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-simple/deploy/main.tf.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataAwsVpcPublicPrivateBuildBuildRubyShTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x56\x7f\x53\x1b\x39\x12\xfd\x7f\x3e\x45\xaf\xa1\xc2\x5d\x15\x9a\x01\xee\xf6\x6e\xd7\x5e\xb6\x0e\x58\x48\xa8\x4a\x99\x94\x49\xee\x47\x25\x39\x6a\x3c\xd3\xb6\x15\xc6\xd2\x44\xd2\x18\x0c\xf1\x77\xbf\xd7\x1a\xe3\x1f\x21\x97\xca\x3f\x30\x96\xd4\xdd\xaf\xbb\x5f\x3f\x69\xe7\xa7\x6c\xa8\x4d\x36\xcc\xfd\x24\x49\x3c\x07\x52\x96\x8c\x6d\xcc\xf2\x93\x9d\xe3\x7b\x1d\x3f\x6b\x5d\xf3\x28\xd7\xd5\x72\x39\xb8\xbc\xe0\x24\xc1\x97\x75\x7f\xfa\x33\x3d\x26\x44\x54\xd9\x22\xaf\xc8\xdb\xc6\x15\x3c\xd2\x15\x1f\xef\x1e\xae\x97\x2b\x6d\xd8\xd8\xe3\xdd\x23\x59\xe2\x62\x62\xa9\x73\x3e\x18\x5c\x0d\x28\x0f\xb4\xfb\xb8\x36\x5a\x74\x77\x1f\xdb\xb3\x8b\x1e\xbd\xce\x7d\x80\xfd\xd8\x77\x3b\x62\x36\x76\x5c\x93\x0d\xc1\x52\x36\xcb\x5d\x86\x8d\xcc\xcf\x3d\xfe\xd1\x17\x0a\x11\x9b\xa1\xa3\x83\x64\x91\x00\x5d\x4d\x7b\x11\x1c\x75\x76\x1f\x4f\x4f\xae\x5f\xdd\x5c\x5f\xbd\x1b\x9c\x9d\x2f\x3a\xb2\xf0\xfa\xb2\x7f\xde\xbf\x5a\x74\xf6\x08\x18\x92\xc4\xb2\xa4\x80\x8d\x7f\x74\xe8\xe8\xf7\x17\x87\x70\x07\xa7\x63\x76\xa4\x42\x1b\xef\x77\xca\x4a\x9e\x65\xa6\xa9\xaa\x1e\x2d\x12\x5b\x45\x83\x36\x8d\xf7\x72\xe2\x23\xc1\x58\xb6\x92\x1d\x2a\x2a\xdb\x94\xaa\xb0\x66\xa4\xc7\x54\xe4\x86\xb4\x09\xec\x46\xec\x98\xee\x74\x98\x50\x5e\x07\x2a\xec\x74\x9a\x9b\xd2\x93\x1e\x91\x0e\x7b\x9e\x7c\xd0\x55\x85\x93\x54\x3b\x8b\x3c\xbd\x47\x10\xea\xfc\x2b\xd7\x41\x9b\x31\x8d\x90\xc8\x96\x5b\x60\x82\x8b\xba\xe2\xc0\x69\x9a\x76\x92\xc6\xc0\x9e\xde\xbf\x27\x35\x5a\x16\x47\x0f\xb3\x68\x91\x69\xe3\x43\x6e\x0a\xce\x86\xd6\x06\x35\xd2\x46\xfb\x09\x97\xf4\xf1\x63\x8f\x4a\x8b\xb2\xfa\x8a\x51\xd6\x83\xf4\xe7\xa4\xb4\x06\x3d\x95\xb8\x27\x65\x29\x61\x05\x29\x6a\x6e\xbd\x0e\xd6\x69\xf6\x04\xc8\xd4\xd4\x65\x2e\xa0\x62\x5c\xcb\xe4\x9b\xd2\xca\x49\x35\x06\x69\xe2\x26\x3f\x5b\x8e\x18\x90\x9f\x9a\x53\x3d\x0f\x13\x6b\x94\xb7\xa3\x70\x97\x3b\x56\xc8\xb7\x66\x17\xc4\xfb\x37\xd6\x94\x14\xca\x9a\xe8\x08\x5d\x35\xbe\xb6\x2e\xa8\x49\x08\xb5\x5f\x07\x29\x4b\x25\xfb\x2b\xa4\xf3\x18\xa7\xce\xbb\xc5\xc4\x69\xaf\x2a\xce\x33\x63\x4b\x4e\x3f\xf9\x2d\x60\x62\xf7\xdc\x66\xe8\xf4\x78\x12\x86\xf6\x3e\x73\xcd\x70\xae\xcc\x78\xcb\xe6\x96\xe7\x88\x37\x23\x25\x5f\x9e\xdd\x0c\x24\x99\xdc\xd6\xdd\x2c\x5b\xfd\x4e\x9b\x21\xba\xd1\xa4\x40\xde\xfd\xe5\x00\x27\x1d\x17\xb3\x78\x9c\x7e\xfe\xdb\xe1\xc5\xaf\xa7\xbf\x9e\x9d\x9c\xfd\xf5\xe0\xf4\xe8\xe2\xef\x49\x64\xd0\x5e\xc9\x43\x8a\x29\xc1\x8d\xf5\x5e\x61\x24\x73\x29\x77\x5a\x4f\x1a\xaf\xad\xa9\x73\xef\xd9\x80\x8f\xe2\x33\x03\x8c\x6c\xb5\x42\xc1\x35\x3e\xcc\x69\x9a\x6b\xb3\x07\xde\x46\xa0\x81\x99\x32\x0e\x45\x3c\xda\x8e\x96\x4f\x2b\xed\x43\x5a\xae\x2d\xe3\xc2\x26\xb1\xff\x5f\x2f\x13\xbe\x97\xa2\xd3\xe0\xdd\xe9\x7f\x6e\xfe\x79\x3e\xb8\xbe\xbc\xea\x1f\x77\x1e\x1f\x49\xea\x73\x83\x84\x05\x22\x2d\x16\x9d\x96\x3a\x97\x6d\xab\x85\x3e\x03\x1c\xd8\xa7\x37\x4f\x11\xf7\xa9\x3f\xd6\xe6\x7e\x3f\xb2\xc8\x86\x09\xd0\xd7\x79\x71\x9b\x8f\x81\x4e\xb8\xb4\x8c\xf3\xc7\xf9\xe9\xe5\x49\xff\xe6\x62\x70\xd5\x7f\x7b\xde\xff\xe3\xd8\x58\x13\x07\x28\x2f\x82\x9e\x7d\x97\x5a\xc3\x07\x47\x63\xc8\xd5\x94\x5d\xd1\x38\x0d\xd5\x19\x36\xba\x2a\x15\x0b\x80\x20\xbf\x3f\x80\xef\x98\x8c\xfa\xb3\x42\xd6\xf4\x80\xcf\xc3\x71\xfc\xfc\x0e\xf5\xc4\x46\xe8\xf3\xc9\x3f\x99\xfb\xcf\x95\x0e\xfc\x97\x68\x28\x4b\x52\x88\xdd\xcd\xf2\x3c\x5f\x59\x9d\x35\x52\x02\xc5\xf7\x20\xb3\xa7\x55\x33\x9e\x95\xee\xb4\x31\x65\x85\x26\x6d\xce\xd8\x98\xa7\xab\x6c\x87\xed\x3e\xd8\x65\xac\x72\x7a\xf9\xbf\xb4\x45\xeb\xe9\x5c\xfc\x17\xa1\x9d\xe1\x3a\x7a\x89\x2e\xa6\xb7\xa5\x86\x51\x4d\x99\x77\xb3\x4c\x84\x0b\x93\x53\xb7\x7b\x21\x77\xf4\x70\x0f\xf9\x08\xd3\x7a\xb5\x95\x86\xf1\x03\xa9\xb3\xaf\xce\x6f\x6b\x44\x5d\xe9\x02\x8a\x80\x52\x35\xfe\x2b\xc8\x18\x31\x59\x03\xbc\x52\xfb\x7c\x58\x71\xa9\x24\xe7\x3b\xeb\x4a\xac\x8d\xb9\xb0\x9e\x3a\x1d\xda\x76\x7c\xcd\x21\x22\x47\x1f\xa6\xda\x0b\xbb\xfc\x96\x53\xcc\xcc\x9d\x21\x35\x58\x99\x75\xbf\x05\xef\x2c\x0a\x25\x68\x00\x4f\xb1\xe8\xd1\x07\xe4\xf9\xed\x44\x43\x76\x3d\x84\xed\x73\xa3\x1d\x94\x50\xc4\x75\x63\xa0\xa4\xd0\x81\x72\xec\xe7\xde\x1a\x01\x4d\x6c\x66\xda\x59\x33\x05\x8b\xe8\x6e\x22\x42\x0e\x96\x41\xd9\xe1\x0d\x7a\x5a\x12\xdf\x73\xd1\x04\x39\xea\xc1\x8f\x5b\x4c\x5f\xe3\x5d\xbc\x59\x61\xb9\xbf\xfe\x05\x56\x56\xfb\x84\xc9\x4c\xe9\x12\x21\x2a\x2f\x34\xae\x41\x3b\x13\xaa\x39\x9c\x19\x66\x5c\x09\x40\x60\x0b\x1c\xa5\x09\x94\x48\xae\x04\x8c\x0a\xb5\xba\x9f\xd2\x49\x5d\xb3\x89\x85\x07\x04\x49\xc4\xf8\x66\x34\xd2\x85\x86\x8f\x94\xba\xea\x4b\xdb\x4c\x8f\xbc\x94\xa6\xbd\x43\x9f\xfd\x57\x40\xd0\x9b\x93\xb7\xaf\x7a\x1f\x4c\xb6\xd7\x2a\x43\xac\x48\xfb\x37\x15\xd7\xdf\xb0\xda\xa1\x2b\x14\xb4\x4b\xf2\x16\x10\x6b\xcc\xc8\x46\x99\xe4\x5e\xf3\xd0\x99\x27\xad\xfa\x8e\x6b\x24\xd6\x47\x62\x92\x97\xe3\xa9\x9d\x31\x12\xd2\xa2\xf6\x6d\x5f\xa4\xd0\xc8\x1a\x32\x45\x50\x62\x6e\x91\xb8\xe9\xa6\x33\x59\xf7\x8a\x63\x33\x4a\x08\xd6\x28\x6f\xaa\xd0\xf6\x92\x3d\xc7\xb7\x05\xee\x26\xb4\xa5\xc6\xcd\x29\x4d\xc2\x6c\xc9\xf4\xe2\xd3\x3f\x15\x70\x05\x5d\x2d\x45\xa7\xa4\x35\xc6\x7d\x0c\x54\x80\xbf\x78\x15\xa3\xef\xba\x25\x42\x09\x31\x00\x13\x3c\xa3\x47\x90\x42\x7a\xba\x7c\x27\x48\x3e\xb4\xe5\xb2\x0d\xa2\x81\xf9\x66\x19\x2f\x4d\x30\x0b\xf4\xdb\x6f\xfd\x97\x97\xfd\x7f\x9f\x5d\xf5\x2f\x9e\x89\x72\x9b\x92\xb8\xda\x92\x63\x59\xd8\x92\xe3\x1d\x7a\xc9\x86\x25\x6e\x49\xc3\x79\x6c\x46\xb2\x3a\x7e\xe3\x70\x99\xb7\xc4\x92\x9b\x5e\xf4\x26\x9b\x81\x18\x16\x3b\xf2\xbd\xbc\x39\x6e\x56\x06\x99\x3c\xc3\x42\x9c\x25\xbc\x00\x7a\x9b\x9e\x70\x7e\x4d\xd1\xf5\xfa\xc8\x31\xc7\xcd\x5e\xb2\x4a\x26\xf9\xc1\xec\xb6\x1b\xb6\x12\x93\x1f\xca\x71\x79\xab\xc6\x27\x25\xc9\x25\xc5\x86\x7e\x39\xe8\xc5\x9f\x6d\xd6\x9b\xc3\x9e\xd5\xcd\x10\xfa\xd3\x6e\xaf\xc1\x2f\x43\x93\x35\x3d\xbc\x07\x37\xf0\x8b\x34\xb4\xe2\xfa\xa4\xa5\x32\x4b\x42\x90\x2d\xa1\x54\xcd\x4a\x5d\x64\x22\x56\xaf\x64\x52\x55\x41\x9d\xa2\xdc\x06\x41\x2f\x5e\x2c\x15\x79\x7d\x1d\x41\xf1\xeb\xca\xce\xa3\x66\x28\x25\x0f\x3f\xa1\x0a\x32\xe7\xca\xd6\x71\x15\x25\x0a\xcb\x4b\x13\x81\xe5\xf1\xf5\x53\x27\xf9\x5f\x00\x00\x00\xff\xff\xbe\x9a\xb3\xa9\x92\x0b\x00\x00"

func dataAwsVpcPublicPrivateBuildBuildRubyShTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsVpcPublicPrivateBuildBuildRubyShTpl,
		"data/aws-vpc-public-private/build/build-ruby.sh.tpl",
	)
}

func dataAwsVpcPublicPrivateBuildBuildRubyShTpl() (*asset, error) {
	bytes, err := dataAwsVpcPublicPrivateBuildBuildRubyShTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-vpc-public-private/build/build-ruby.sh.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataAwsVpcPublicPrivateBuildTemplateJsonTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x54\xcd\x8e\x9b\x30\x18\xbc\xf3\x14\x16\x52\xf6\x14\x20\xed\xae\xaa\xaa\xd7\x3e\xc6\x2a\x62\x0d\x78\x83\x15\xdb\x58\xfe\xec\x54\x59\xe4\x77\xef\x67\xfe\xa9\x36\x24\xdd\x5c\x8c\x3c\xc3\xcc\x78\xf8\x9c\x36\x22\xf8\x8b\x25\x57\xb9\xa6\xe5\x99\x99\xfc\xc2\x0c\xf0\x46\xc5\xbf\x48\x7c\x48\x7f\xa6\x87\x78\x1f\xf5\x9c\x0b\x35\x9c\x16\x82\x01\x42\xfd\x6b\xb8\x49\xff\x40\x4e\xcb\x92\x01\xe4\x67\x76\x45\x44\x39\x21\xf6\x4b\x14\x58\x69\x98\xbd\x85\x1a\x76\xea\xcd\x56\x08\x08\x77\xc2\x3c\xb6\x1e\x80\x6e\xdf\x8f\x41\xb4\x69\x2e\x3c\x64\xc4\xa4\x48\x78\x1d\xde\x6a\x77\xe4\xbd\x31\xa4\xe2\x86\x70\x85\x8f\x4e\x55\xd4\x22\x2b\xc7\x1d\x48\x0b\xc7\x45\x45\x76\x7e\x24\x0f\x2b\xca\xd9\xab\x66\xe1\xb4\x50\x33\x21\xe2\xfd\x0c\x70\x25\xb8\x0a\xd0\x6b\x2c\xcf\x41\x36\xd1\x24\xb3\x52\x67\x8d\xb5\x4d\x36\x1b\x24\x6d\x1b\x9c\x45\xd3\xe8\xf4\x37\xee\x5a\x66\x88\xf7\xf1\x71\x50\xf2\xfb\xdb\x9e\xef\x5c\xb0\xa5\x25\x34\xce\x94\x1d\x82\x9a\xc1\xd2\xfb\x6c\x89\x57\x0c\x2c\x57\x9d\x6b\x20\xfd\x47\x9a\x07\xc2\x6c\x15\x50\x56\x8f\x1e\xdd\x7b\xf2\xf4\x44\x0a\x0a\x35\x49\x33\x49\xb9\x4a\xa1\xfe\xa4\x8b\x1d\x61\xaa\x0a\xdf\x6b\xeb\x93\x6c\xd4\xb3\x23\x38\xa8\x05\x86\x90\xa8\x80\x29\x1c\xe0\x39\xdf\xa6\xc1\x79\xc3\x33\xf7\x1e\x0b\xda\x23\x4d\x26\x54\xeb\xd4\x9e\x3e\xbe\x54\x18\x94\x86\x6b\x1b\xa0\x6e\xdc\x12\xe3\x8a\x6b\x38\xfe\xa8\xd5\xad\xc7\x71\x8e\x3b\xce\x30\xc3\xd3\x85\x52\x54\x76\xda\x21\xcb\x24\x3d\x39\x52\x49\x3f\xb0\x75\x56\xc0\x8c\xad\xae\xdf\xad\x62\xd6\xf7\x74\xbb\x9d\x78\x75\x65\xb7\x14\x67\xe2\x1d\xc5\xe9\x9a\x6f\xa9\xf5\xa4\x7b\xd9\xba\x11\xc8\xa9\xe4\x7d\x1f\x3c\xf9\xfe\xed\xc7\xf3\xa1\x7a\x79\x99\x39\x5c\x81\xa5\x0a\x59\x63\x6d\xe5\x73\x2a\xa8\x39\xb1\x85\x0c\xd4\x79\x70\x1e\xeb\x76\x05\x0e\xaf\x5b\x94\x2a\x79\x3e\x62\x6d\x1b\x9e\x70\xae\xff\xcd\x8e\x2b\x4e\x11\x95\xfa\xb3\xc4\xfd\x7f\xd6\x31\x8a\x7c\xf4\x37\x00\x00\xff\xff\x16\x08\x5f\x5b\x65\x05\x00\x00"

func dataAwsVpcPublicPrivateBuildTemplateJsonTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsVpcPublicPrivateBuildTemplateJsonTpl,
		"data/aws-vpc-public-private/build/template.json.tpl",
	)
}

func dataAwsVpcPublicPrivateBuildTemplateJsonTpl() (*asset, error) {
	bytes, err := dataAwsVpcPublicPrivateBuildTemplateJsonTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-vpc-public-private/build/template.json.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataAwsVpcPublicPrivateDeployMainTfTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x55\x4d\x6f\x13\x31\x10\xbd\xe7\x57\x8c\x5c\x4e\x88\x6e\x0a\xa7\xaa\x12\xb7\x4a\xdc\xe8\x85\x1b\x42\x2b\xaf\xd7\x09\x56\x1d\xdb\xf2\x47\xd0\x2a\xda\xff\xce\xd8\x8e\xe3\xfd\x48\x5b\x10\x12\xdb\x4b\xfb\xde\x8c\x67\xe6\x3d\x7b\x7a\x03\x5f\xb8\xe2\x96\x7a\xde\x43\x37\xc0\x93\xf7\xfa\x03\xf4\x1a\x94\xf6\xc0\x7b\xe1\xe1\x40\x55\xa0\x52\x0e\x9b\xcd\x91\x5a\x41\x3b\xc9\x81\x08\xb5\xb3\xb4\x15\x3d\x81\xd3\x38\x81\xe9\x2f\xd7\x52\xc6\xb8\x73\xed\x33\x1f\xae\x90\x8e\x33\xcb\xfd\x0b\xa4\xe5\x7b\xa1\xd5\x82\xc0\xd0\x56\xd1\x03\x4f\xf0\x34\xe1\x20\x16\x91\x42\x39\x4f\x15\xe3\xad\x1f\x4c\x0c\x87\x9e\xef\x68\x90\x1e\x3e\x03\xf1\x9f\x9a\x83\x60\x56\x13\x98\x66\x18\x2b\x8e\x38\x77\xeb\x42\xa7\xb0\xab\xd5\x38\x26\x74\x52\xb0\x17\xe9\xa3\x61\x2d\x13\xbd\xbd\x02\x9f\x63\x37\xc6\xea\xa3\xe8\xb9\x4d\x03\x22\xb4\x01\xa8\xfa\xc4\xc6\xde\x9d\x30\xb1\x99\xeb\x36\x12\x0c\xab\x4a\xcd\xc3\x2a\x9e\xc2\xb2\x66\x10\xbf\x59\x58\xc6\x31\x04\x9b\xb0\xdc\xe9\x60\x59\xb5\x20\x58\xe1\x87\x76\x6f\x75\x30\x04\x08\x97\x5d\xee\x2c\xca\x1c\x4f\x39\x9d\xf2\xaf\xe3\x78\x8b\xdc\x6d\x3e\xb4\x38\x9e\xaa\xe6\x11\x6b\xc5\xfc\x37\x52\xc8\xf1\x3d\xd6\x73\xe9\x40\x00\x9c\xdf\x6b\xa6\x65\xee\xef\xf6\x63\x02\x77\x56\x1f\x5a\xa3\xad\x4f\xe0\x5d\xc2\xbc\x2e\x48\xc5\xa2\xb6\x6d\x27\x35\x7b\x76\x88\x7d\x27\x77\x4d\xfa\xd9\xde\x91\x1f\xc8\x8f\xb1\x98\x50\x2f\x57\x23\x9e\x19\x72\xa5\xe0\xfd\xb5\x8a\xf7\x7f\x56\xf2\x6d\x35\xa9\x31\x13\x35\x61\xa1\xe7\x5f\x6a\xf9\xda\x78\xff\x28\x66\x2d\x16\x99\xf1\x3c\xdf\xff\xb4\x6f\xa5\x65\xba\x88\x2b\x01\x2f\xdf\x9b\x4a\xe6\x77\xea\x26\x09\x65\xcc\xe5\x43\xce\xe3\xce\xbd\x2b\xb2\xac\x5d\x6d\xb0\xb1\xa6\x24\x95\x2d\xe3\x66\x45\x62\x52\x61\x1a\x9c\xa0\x79\x7f\x4e\xc0\x8c\x1b\xf8\xf6\xf4\xf8\xf4\x80\x7b\xf4\x99\x83\x14\xce\x73\x85\xbe\x42\xd4\xcb\x01\xd3\x6a\x27\xf6\xc1\xc6\xd5\x81\xb1\x99\xc6\x7d\x91\xf5\x97\x5d\x95\x15\xe6\x37\x35\x52\x13\x77\x16\x37\xfe\xb2\x0b\xd7\x57\xbc\x52\x25\xbd\x26\x26\x53\x6e\xe0\x91\x1b\xa9\x07\xa0\xa8\x90\x07\xbd\xab\x33\x2f\x0c\x2b\xf8\xd4\x35\xdc\xcb\x73\xcf\xce\x3b\xe9\x20\x92\x47\xb3\x25\x5d\xe9\x19\x3c\x31\x33\xbe\x8c\xd9\x39\xab\x95\x9d\x82\xcb\x3f\x89\x45\xd1\x02\xe7\xc7\x14\xef\xfa\xdc\x58\x4c\x7f\xc5\xf5\x68\xe3\xc5\x44\x4f\xf7\xe5\x51\x7c\x5d\xad\xc9\x8b\x74\x3a\x78\x13\x3c\x90\x60\x65\x56\xe3\x48\x65\x48\xc1\x3f\xbd\x37\x0f\xdb\x6d\x2e\x14\xef\x53\x3c\xbd\x57\x2e\xf7\xb7\x8d\x7b\xfa\x77\x00\x00\x00\xff\xff\x11\x23\xa2\xbf\x89\x07\x00\x00"

func dataAwsVpcPublicPrivateDeployMainTfTplBytes() ([]byte, error) {
	return bindataRead(
		_dataAwsVpcPublicPrivateDeployMainTfTpl,
		"data/aws-vpc-public-private/deploy/main.tf.tpl",
	)
}

func dataAwsVpcPublicPrivateDeployMainTfTpl() (*asset, error) {
	bytes, err := dataAwsVpcPublicPrivateDeployMainTfTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aws-vpc-public-private/deploy/main.tf.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _dataCommonDevVagrantfileTpl = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x56\xeb\x6e\xdb\x3a\x12\xfe\xaf\xa7\x98\xa3\xf8\x34\x2d\x10\x49\x27\x69\xb1\x3f\xdc\x26\x68\x6e\x6d\x02\x74\xe3\xc2\x4e\x0b\x2c\x76\x17\x2e\x2d\x8d\x25\xb6\x32\xc9\x92\x94\x13\xd7\xf1\xbb\xef\x0c\x25\x5f\xb2\xe9\x76\x0b\xd4\x0d\x35\xe4\xdc\xbe\x99\xf9\xc8\x3d\x78\x8f\x0a\xad\xf0\x58\xc0\x64\x01\x03\xef\xf5\x01\x14\x1a\x94\xf6\x80\x85\xf4\x7f\x44\x7b\xd1\x1e\xdc\x56\xd2\x01\xfd\xf3\x15\xc2\x67\x51\x5a\xa1\xfc\x54\xd6\x08\xe5\x7f\xeb\xc2\x54\xdb\x70\xaa\xc0\x39\xd6\xda\xcc\x50\x79\xd0\x53\x32\xe1\xd9\x84\x30\xa6\x96\xb9\xf0\x52\xab\xcc\xa1\x9d\xcb\x1c\x53\xb8\xf6\xe0\x2a\xdd\xd4\x45\x70\x3a\x41\xa8\x84\x2a\x12\x76\x8e\x45\x0a\xb7\x1a\x66\xba\x90\xd3\x05\x9b\x25\x3b\x3b\xee\x0f\xa0\x71\x18\xbc\x9d\x1a\xc3\x82\x34\x8a\xba\xed\x34\xd7\x6a\x2a\xcb\xc6\xe2\xf3\xf8\x28\x7e\xc1\x19\x3d\xb4\xa2\x87\x08\xa0\x5d\xa5\xf3\x59\x3a\xd1\xf7\x70\x0c\x71\x25\x5c\x25\x73\x6d\x4d\x66\x2c\xe6\xd2\xe1\xdf\x5e\xc5\x11\x1d\xdc\x83\x2b\xed\x28\x01\x55\x2f\x40\xa1\xbf\xd3\xf6\xdb\x23\xf5\x4e\x06\xb1\xb1\x72\x4e\x38\x8c\x3b\x41\x7c\x00\xd2\xf4\x21\x5e\x2e\x19\x88\xb1\x34\x63\x51\x14\x16\x9d\x83\xd5\xaa\x33\x3c\x42\xdf\x18\x10\xe0\x16\x2a\x27\xfc\xa6\xba\x2e\xd0\xc2\xd4\xea\x19\xe8\xc6\x02\x5b\x91\xaa\x84\x42\x52\x40\x5e\x5b\x4a\x5f\x43\x36\x6f\xb3\x7b\x14\x43\x6b\x60\xdc\x19\xd8\x27\x97\x46\xf8\x2a\x5d\x1b\x58\xad\xf6\x0f\x20\x5e\x6b\xc6\x07\xa4\x0b\xa0\xef\xa8\x6e\x14\xdf\x46\x0a\xa5\xd5\x8d\xd9\x91\xb4\x41\x5e\x2a\x31\xa1\x32\x8f\x46\x57\x20\x4a\x2e\x25\x95\xf7\x4e\xd8\x82\x0d\x3b\x4d\xe5\xf7\x9e\x97\x5d\xf6\x94\xab\x41\x55\xa0\xca\x25\xba\x90\x81\xdb\x46\xea\x5c\x95\x76\xda\xe3\xd6\xd6\x31\x78\xdb\x60\xeb\xe8\x9d\x6e\x54\x11\xfa\x02\xd6\x95\x6b\xbf\x9e\xcb\x29\x08\xb5\x78\x41\xa7\x96\x7f\x86\xee\x22\x44\x40\x2a\x5a\xae\x35\xc6\x24\x71\x29\xe1\x0c\x7f\xae\xe8\x18\xef\x53\x49\x33\x4d\xed\x98\x6d\x4f\x25\x04\x0c\xa9\xd7\x5a\x9b\xf4\x9c\xa4\x9e\xc0\xe2\x62\xfc\x1a\x4a\x36\x16\x10\xa4\xc5\xa3\xa3\xc6\xea\xb9\x74\x1c\x61\xec\x2a\xac\x6b\xae\xb8\xaa\xa5\x42\xc2\x30\x2f\x60\x6f\x49\x0a\x2b\x78\xf6\x0c\x26\xd4\x5a\xdd\x67\x36\x13\x52\xa5\xae\x8a\xdb\x64\x08\x2a\xce\x87\x82\x0e\x10\x7c\xd0\xa2\x00\x51\xd7\xa1\xfc\x53\x2b\x4a\x9e\x1d\x07\x15\x5a\x0c\x79\x13\x0a\x8f\x00\x4e\xb7\x90\xac\x4f\x33\x2e\xdc\x6f\x5b\xed\x80\x08\x67\xde\x49\x1e\x2c\x92\x97\xd5\xea\xa7\x11\x5c\x2b\xe7\x39\x80\x61\x43\xd3\x3c\x69\x24\x4d\x24\xaa\xb9\xb4\x5a\xb1\xea\xef\xa6\xdf\x73\xb9\x95\xc6\x8f\x2d\x59\x89\x9e\x2a\x31\xb2\x7d\x23\x2c\x39\xc2\xda\x85\xc1\x34\x07\xa0\x1f\xda\xbe\xdc\x99\xc9\xcd\x99\xac\x99\x50\xb9\x9a\xe4\xf0\x28\xfd\xeb\x15\x63\x47\x61\x47\xfc\x8b\x76\x5d\x91\xca\x9b\x37\xa3\xf3\xe1\xf5\xc7\xdb\xc8\xa1\x87\x84\x39\xac\x51\xdd\x12\xad\xc5\x7b\x19\x96\x46\x1a\x9c\x0a\x59\x77\x62\x6f\x45\x4e\x4d\x48\x2b\x6d\x9f\xbf\x80\x25\xc7\x51\xeb\x5c\xd4\xd4\xe0\x8d\xcd\x91\x79\xe5\xb8\x77\xb8\x15\x73\x96\x4a\x1f\xf7\x8e\x58\x84\x79\xa5\x21\xbe\x1c\x0e\x07\x43\x10\x1e\x7a\xcb\xad\xd2\xaa\xdf\x5b\xb6\x67\x57\xaf\xe1\x83\x20\x1a\xa9\x75\xe9\xfa\x9c\x00\x4d\x1b\x1a\xe0\x06\xe5\x91\xb6\x19\x6d\x64\x6e\xe1\xe8\x0f\x3c\x80\x0f\xb1\x29\x38\xfa\x2b\x5a\x45\x14\x9d\x81\xfd\x10\x1c\xc4\xbd\xe5\xd9\xe9\xe8\x6a\x3c\x1a\x7c\x1a\x9e\x5f\xae\x62\x16\x7c\xb8\xbe\xb9\xbc\x19\xac\xe2\x7d\xa0\x18\x22\xa2\x47\x36\x9a\xe0\x3d\xe6\x7d\xe0\xff\x1b\x9a\xca\x5c\xcf\x66\xc4\xa8\x70\x27\x7d\x45\xcd\xe5\x4d\x13\x42\x29\x99\xb5\x69\xc9\xa4\x5b\x48\x67\x6a\xb1\xc0\x22\xd2\xc8\x20\x40\xef\x2d\x1c\x9d\x3c\x3b\xa4\x70\xc2\x49\x0b\x89\x6f\xe3\x3d\x81\x8c\x3a\x2c\x53\x4d\x5d\xbf\x86\xd5\xc6\x23\x9d\xea\xaf\x6d\x0b\xe2\x03\x42\xe0\x9e\xec\xcf\x88\xf2\x68\xd8\x23\x5d\x07\xab\x2d\x5a\xff\x64\x8d\x7f\x93\x8b\xb8\xb3\xf0\x77\xf1\x0d\x81\x8a\x43\x8c\xe2\x2b\x42\xf1\x4b\x47\x42\x40\x9c\xf1\x05\x4a\x4d\x64\xd2\xd2\x60\x1d\x58\x90\x09\x9f\xb8\x9a\x05\x61\x2c\x5b\xab\x34\x74\x1b\x92\x83\x13\x0a\xb3\xd2\x33\x5c\x4b\xb2\x94\xc7\xd0\xe6\xec\xed\xbc\xe3\x17\x26\x2e\x26\xb6\x30\x40\x54\x1e\x4a\x92\xb2\x90\x2a\x22\xc6\xf9\xa3\xad\x50\xfc\xc9\xe1\xc5\xcd\x88\x20\x8a\x21\x43\x9f\x67\x14\x10\xff\x8a\x71\xdb\xd4\x70\xb2\x03\x06\x85\xa5\xa2\x75\x47\xec\x28\x3e\x80\x6b\xa8\xcb\x3d\x22\x24\xe2\xff\x99\x21\x03\x1a\x5b\x85\xee\x7e\x64\x10\x80\xae\x0e\x2f\xac\x8f\xa6\x92\x3a\xf5\xde\x68\xeb\xe1\xe2\xf2\xec\xfa\xf4\x66\xfc\x6e\x38\xb8\xb9\xbd\xbc\xb9\x38\x56\x5a\x49\x26\x35\x91\x7b\x39\xa7\x86\xd6\x35\xc4\xa7\x45\x60\x6a\x61\x3c\x59\x30\xda\x49\xba\x48\x98\x9a\xb9\x1b\x1a\xc3\xbc\xa8\xca\x34\x4d\xe3\x68\xed\x93\x4e\x26\x44\xea\xed\x26\x3e\x11\xcb\x8e\x20\x92\x05\x98\x85\xaf\x88\x54\x9d\x9e\x7a\xa2\x74\x4c\x68\xb6\x0d\x5a\xcf\xd6\x7f\x22\x4b\xb8\x07\x89\x2c\xd8\x10\xb5\xb4\x72\x9c\x42\x52\x79\x6f\xdc\xd6\x49\x51\x24\xbc\xbf\x89\x74\x11\xfc\x18\xd1\xcf\x2b\x2b\x5d\x52\xa3\xc8\x94\x2e\x30\xfd\xea\x1e\x05\xc6\x7a\x4f\x75\x26\x56\x96\x95\x27\x2a\xc9\x98\x1c\x12\x55\xfe\xaf\x1c\xf9\x75\x33\xb8\x18\xf4\x81\xf9\x66\x86\x84\xa0\xfc\x81\x10\x18\x65\x8e\x36\x50\x9c\xa0\xb7\x8f\xa2\xce\x5e\x43\x3f\xfc\x74\xf6\x8f\xf1\xe7\xcb\xe1\xe8\x7a\x70\x73\xcc\x77\x3c\x9f\x1e\xaf\x4f\x87\x1b\x9e\xd1\xef\xe8\x94\x2b\x10\x18\xb5\xb7\xdc\x55\x5c\x85\x2a\xb8\xc6\xb0\xc9\x70\x89\x8a\xfc\x1b\xcd\x8a\x0b\x05\xf9\xbd\x22\xff\xa2\x3e\x93\x1f\x16\x4a\x9a\xa9\x19\xda\x9c\x7a\x9d\x78\x2b\x10\x7a\x42\x03\x49\x5c\xce\xdf\xff\xa2\x5e\xab\xe5\x84\x76\xeb\x97\x9b\x45\x52\xaa\xc6\x13\x2d\x77\xdf\xaf\x12\xaa\xa0\x72\xae\x4e\xf8\x7a\xed\x34\xcc\xf7\xf0\xf5\x83\x96\x87\x65\x58\xfe\xa2\xe2\x9d\x8e\xfb\x5e\xd3\x73\xee\xe5\xc6\x0c\x17\xf2\xab\x0b\x4b\x06\xaf\xb7\x8b\xcc\x53\x09\xab\x3d\xc1\xf4\x8c\x6e\xf6\x1a\xed\xba\x7f\x4b\x9c\x6d\xf2\x9f\xb4\x5b\x90\x24\x4a\x27\x56\x76\x7f\x0b\x9d\xb7\x46\x76\x19\xe0\x3d\x61\x44\x9c\xc2\x0f\x49\x26\x03\xb6\xc0\x37\xa4\x9e\xc2\xd5\xed\xed\x47\xa6\xa4\x3b\x62\x1b\xa1\xda\xf7\x4f\xd2\xbd\x60\x36\x2f\x1e\x6e\x3c\x10\x0d\xbd\xb7\xd6\x61\x90\xbd\x6e\xaa\x93\xa4\xac\xf5\x84\x90\x26\x20\xd3\x98\x36\xde\xd2\xaf\x6a\x26\xf4\x36\x9d\xf5\xe3\xb4\x73\x35\x98\xd2\xf3\x93\x07\xa1\x9f\x65\xdb\xfd\x2c\x8e\xba\xbb\xec\x3f\x01\x00\x00\xff\xff\x7f\xfb\x67\x82\xa1\x0b\x00\x00"

func dataCommonDevVagrantfileTplBytes() ([]byte, error) {
	return bindataRead(
		_dataCommonDevVagrantfileTpl,
		"data/common/dev/Vagrantfile.tpl",
	)
}

func dataCommonDevVagrantfileTpl() (*asset, error) {
	bytes, err := dataCommonDevVagrantfileTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/common/dev/Vagrantfile.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"data/aws-simple/build/build-ruby.sh.tpl": dataAwsSimpleBuildBuildRubyShTpl,
	"data/aws-simple/build/template.json.tpl": dataAwsSimpleBuildTemplateJsonTpl,
	"data/aws-simple/deploy/main.tf.tpl": dataAwsSimpleDeployMainTfTpl,
	"data/aws-vpc-public-private/build/build-ruby.sh.tpl": dataAwsVpcPublicPrivateBuildBuildRubyShTpl,
	"data/aws-vpc-public-private/build/template.json.tpl": dataAwsVpcPublicPrivateBuildTemplateJsonTpl,
	"data/aws-vpc-public-private/deploy/main.tf.tpl": dataAwsVpcPublicPrivateDeployMainTfTpl,
	"data/common/dev/Vagrantfile.tpl": dataCommonDevVagrantfileTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"data": &bintree{nil, map[string]*bintree{
		"aws-simple": &bintree{nil, map[string]*bintree{
			"build": &bintree{nil, map[string]*bintree{
				"build-ruby.sh.tpl": &bintree{dataAwsSimpleBuildBuildRubyShTpl, map[string]*bintree{
				}},
				"template.json.tpl": &bintree{dataAwsSimpleBuildTemplateJsonTpl, map[string]*bintree{
				}},
			}},
			"deploy": &bintree{nil, map[string]*bintree{
				"main.tf.tpl": &bintree{dataAwsSimpleDeployMainTfTpl, map[string]*bintree{
				}},
			}},
		}},
		"aws-vpc-public-private": &bintree{nil, map[string]*bintree{
			"build": &bintree{nil, map[string]*bintree{
				"build-ruby.sh.tpl": &bintree{dataAwsVpcPublicPrivateBuildBuildRubyShTpl, map[string]*bintree{
				}},
				"template.json.tpl": &bintree{dataAwsVpcPublicPrivateBuildTemplateJsonTpl, map[string]*bintree{
				}},
			}},
			"deploy": &bintree{nil, map[string]*bintree{
				"main.tf.tpl": &bintree{dataAwsVpcPublicPrivateDeployMainTfTpl, map[string]*bintree{
				}},
			}},
		}},
		"common": &bintree{nil, map[string]*bintree{
			"dev": &bintree{nil, map[string]*bintree{
				"Vagrantfile.tpl": &bintree{dataCommonDevVagrantfileTpl, map[string]*bintree{
				}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        // File
        if err != nil {
                return RestoreAsset(dir, name)
        }
        // Dir
        for _, child := range children {
                err = RestoreAssets(dir, filepath.Join(name, child))
                if err != nil {
                        return err
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


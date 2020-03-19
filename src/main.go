package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func parseArgs() string {
	var file string
	flag.StringVar(&file, "c", "", "config file")
	flag.Parse()

	var t []string
	if file == "" {
		t = append(t, "-c")
	}

	if len(t) > 0 {
		fmt.Fprintf(os.Stderr, "Missing required options: %s\n", strings.Join(t, ", "))
		os.Exit(1)
	}
	return file
}

func errExit(info string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, info)
		os.Exit(1)
	}
}

func loadConfigFile(path string) {
	f, err := os.Open(path)
	errExit(fmt.Sprintf("打开文件 %s 失败，%s", path, err), err)

	b, err := ioutil.ReadAll(f)
	errExit(fmt.Sprintf("读取配置文件 %s 内容失败，%s", path, err), err)

	if err = yaml.Unmarshal(b, conf); err != nil {
		errExit(fmt.Sprintf("解析配置文件 %s 失败，%s", path, err), err)
	}

}

func main() {
	loadConfigFile(parseArgs())
	ldapConf = conf.Ldap
	ifile = conf.File
	var f file

	switch ifile.Type {
	case "csv":
		f = &fromCSV{}
	case "excel":
		f = &fromExcel{}
	}

	importUser(f.getLdapUserAttr(ifile.Path))
}

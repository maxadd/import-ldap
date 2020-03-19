package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	ldapDNAttr     = "dn"
	ldapPasswdAttr = "userPassword"
)

type fileLine struct {
	passwordIndex int
	dnIndex       int
	// baseDN        string
	head []string
	// delimiter     string
	num int
}

type fromCSV struct{}

func (p *fromCSV) readFile(path string) *csv.Reader {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取 csv 文件 %s 失败, %s\n", path, err)
		os.Exit(1)
	}
	return csv.NewReader(bufio.NewReader(f))
}

// 判断是否在第一行中存在 userPassword 以及 dn 这两列
// 如果存在返回其索引也就是列数，否则它的值就是 -1
// 这两个字段都可以不指定
func (p *fileLine) columnParse(head []string) {
	if len(head) < 3 {
		log.Fatalln("文件列数少于 3")
	}
	p.passwordIndex = -1
	p.dnIndex = -1
	for i, v := range head {
		switch v {
		case ldapPasswdAttr:
			p.passwordIndex = i
		case ldapDNAttr:
			p.dnIndex = i
		case "":
			fmt.Fprintf(os.Stderr, "导入文件第一行的第 %d 列为空\n", i+1)
			os.Exit(1)
		}
	}
	if p.dnIndex == -1 {
		fmt.Fprintf(os.Stderr, "导入的文件中缺少 %s 字段\n", ldapDNAttr)
		os.Exit(1)
	}
	return
}

func (p *fileLine) lineParse(line []string) (attrs *ldapUserAttribute) {
	if len(p.head) != len(line) {
		fmt.Fprintf(os.Stderr, "第 %d 行的列数必须和第一行相同", p.num+2)
		os.Exit(1)
	}

	// 用来存放一个人的所有属性，每个属性都是一个 key 对应一个切片
	// 因为 ldap 中的属性可能对应多个值，比如 objectclass，因此 key 对应的值是切片
	var userAttr = make(map[string][]string)
	for i, v := range line {
		if v == "" {
			fmt.Fprintf(os.Stderr, "导入文件第 %d 行 %d 列为空\n", p.num+2, i+1)
			os.Exit(1)
		}
		if i == p.dnIndex {
			continue
		}
		// 用来存放一个单元格中使用分隔符分割的属性
		var values []string
		if ifile.Delimiter == "" {
			values = []string{v}
		} else {
			values = strings.Split(v, ifile.Delimiter)
		}
		// 根据索引从 head 里面取 key，也就是 ldap 中的字段名称
		userAttr[p.head[i]] = values
	}

	var password string
	if p.passwordIndex > 0 {
		password = userAttr[ldapPasswdAttr][0]
	} else {
		var pg passwdGen
		password = pg.get()
	}

	userAttr[ldapPasswdAttr] = []string{getLdapUserPassword(password)}
	attrs = &ldapUserAttribute{
		dn:       line[p.dnIndex] + "," + ldapConf.BaseDN,
		password: password,
		attrs:    userAttr,
	}
	p.num++
	return
}

// 将文件内容转换成 ldap 属性和值
func (p *fromCSV) parseFile(path string) (attrs []*ldapUserAttribute) {
	r := p.readFile(path)
	// head 是第一行，用来告知每一列对应 ldap 中的哪个属性，比如 cn、sn、mail、mobile 等
	// 要判断是否存在 userPassword 列和 dn 列，dn 列必须存在
	head, _ := r.Read()
	fl := &fileLine{head: head}
	fl.columnParse(head)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		attrs = append(attrs, fl.lineParse(line))
	}
	return
}

func (p *fromCSV) getLdapUserAttr(path string) []*ldapUserAttribute {
	return p.parseFile(path)
}

type fromExcel struct {
	delimiter string
	baseDN    string
}

func (p *fromExcel) readFile(path string) *excelize.File {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "打开 excel 文件 %s 失败, %s\n", path, err)
		os.Exit(1)
	}
	return f
}

func (p *fromExcel) parseFile(path string) (attrs []*ldapUserAttribute) {
	r := p.readFile(path)
	rows := r.GetRows("Sheet1")
	fl := &fileLine{head: rows[0]}
	fl.columnParse(rows[0])
	for _, row := range rows[1:] {
		attrs = append(attrs, fl.lineParse(row))
	}
	return
}

func (p *fromExcel) getLdapUserAttr(path string) []*ldapUserAttribute {
	return p.parseFile(path)
}

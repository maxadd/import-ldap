package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	ldap "gopkg.in/ldap.v3"
)

type ldapOps struct {
	conn *ldap.Conn
}

func (p *ldapOps) getLdapTLSConn() {
	tlsConfig := &tls.Config{InsecureSkipVerify: ldapConf.TLS.InsecureSkipVerify}
	address := ldapConf.IP + ":" + strconv.Itoa(ldapConf.Port)
	l, err := ldap.DialTLS("tcp", address, tlsConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接 ldap 服务器 %s 失败，%s\n", address, err)
		os.Exit(1)
	}

	err = l.Bind(ldapConf.BindDN, ldapConf.BindPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "使用用户 %s 登陆 ldap 失败，%s\n", ldapConf.BindDN, err)
		os.Exit(1)
	}
	p.conn = l
}

func (p *ldapOps) getLdapConn() {
	address := ldapConf.IP + ":" + strconv.Itoa(ldapConf.Port)
	l, err := ldap.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接 ldap 服务器 %s 失败，%s\n", address, err)
		os.Exit(1)
	}

	err = l.Bind(ldapConf.BindDN, ldapConf.BindPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "使用用户 %s 登陆 ldap 失败，%s\n", ldapConf.BindDN, err)
		os.Exit(1)
	}
	p.conn = l
}

func (p *ldapOps) addUser(attrs []*ldapUserAttribute) {
	for _, v := range attrs {
		addRequest := ldap.NewAddRequest(v.dn, nil)
		for attr, values := range v.attrs {
			addRequest.Attribute(attr, values)
		}
		if err := p.conn.Add(addRequest); err != nil {
			fmt.Fprintf(os.Stderr, "导入用户失败，%s\n", err)
			os.Exit(1)
		}
		printUserAndPassword(v.dn, v.password)
	}
}

func printUserAndPassword(dn, password string) {
	start := strings.Index(dn, "=")
	end := strings.Index(dn, ",")
	if start > 0 && end > 0 && end > start {
		fmt.Println(dn[start+1:end] + " " + password)
		return
	}
	fmt.Println(dn + " " + password)
}

func importUser(attrs []*ldapUserAttribute) {
	lo := &ldapOps{}
	if conf.Ldap.TLS == nil {
		lo.getLdapConn()
	} else {
		lo.getLdapTLSConn()
	}

	lo.addUser(attrs)
}

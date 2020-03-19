package main

type configItem struct {
	Ldap *ldapConfig   `yaml:"ldap"`
	File *importedFile `yaml:"file"`
}

var conf = &configItem{}
var ldapConf *ldapConfig
var ifile *importedFile

type ldapConfig struct {
	IP           string `yaml:"ip"`
	Port         int    `yaml:"port"`
	BindDN       string `yaml:"bind_dn"`
	BindPassword string `yaml:"bind_password"`
	// 这是配合用户 dn 属性用的，因为一个用户的 dn 会很长，每添加一个用户必须得提供 dn
	// 这样写起来会很麻烦，如果所有用户的 ou 都相同，那么可以将其作为 BaseDN
	// 如果要导入一个 uid=023123,ou=it,dc=example,dc=com 的用户
	// 可以将 BaseDN 指定为 ou=it,dc=example,dc=com，那么用户 dn 列只要为 uid=023123 即可
	BaseDN string         `yaml:"base_dn"`
	TLS    *ldapTLSConfig `yaml:"tls"`
}

type ldapTLSConfig struct {
	InsecureSkipVerify bool `yaml:"insecureSkipVerify"`
}

type importedFile struct {
	Path string `yaml:"path"`
	// csv or excel
	Type string `yaml:"type"`
	// 一个单元格的内容可能包含多个值，多个值通过这个分隔符进行分割
	Delimiter string `yaml:"delimiter"`
}

// 导入用户的属性
type ldapUserAttribute struct {
	// 要导入用户的完整 dn，比如 uid=023123,ou=it,dc=example,dc=com
	dn string
	// 该用户的密码，这是明文，以便导入后打印出来
	password string
	// 这是该用户的所有 ldap 属性，包含加密后的密码
	attrs map[string][]string
}

type file interface {
	getLdapUserAttr(string) []*ldapUserAttribute
}

通过 csv/excel 导入用户到 ldap。

文件第一行为 ldap 属性，必须存在以下属性：

- `dn`：用户的完整 dn。如果配置文件中定义了 base_dn，那么这个 dn 就是 base_dn 的相对值。当 base_dn 为 ou=it,dc=example,dc=com，dn 为 uid=023123 时，用户的完整 dn 为 uid=023123,ou=it,dc=example,dc=com；
- `objectClass`：这个必须存在，通常它的值会有多个，多个值之间使用配置文件中 `delimiter` 指定的分隔符分隔。

`userPassword` 作为用户密码可以不指定，如果不指定会自动生成，密码强度为 16 位大小写字母、数字和特殊字符混合；如果指定使用指定的密码。最终密码存入 ldap 后会对密码进行 `sha512_crypt` 加密（Linux 服务器的密码就是这种加密方式），目前不支持其他加密手段。

可以查看 file.csv 和 config.yml 作为参考。

每导入一个用户，会输出该用户的登录名和密码。
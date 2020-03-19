package main

import (
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
	"log"
	"math/rand"
	"time"
)

const (
	passwordLength         = 16
	passwordComplexity     = 4
	round                  = 65535
	saltLength             = 16
	ldapUserPasswordPrefix = "{CRYPT}"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type passwdGen struct {
	l            []int
	uniqSlice    []bool
	mappingArray [4]string
	passwords    []byte
}

func (p *passwdGen) init() {
	arr := [passwordLength]int{}
	p.l = arr[:]
	p.uniqSlice = make([]bool, 4)
	p.passwords = make([]byte, passwordLength)
	p.mappingArray = [4]string{"0123456789",
		"abcdefghigklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"!#$%&()*+,-./:;<=>?@[]^_`{|}~",
	}
}

func (p *passwdGen) validate() bool {
	for i := 0; i < len(p.uniqSlice); i++ {
		if !p.uniqSlice[i] {
			return true
		}
	}
	return false
}

func (p *passwdGen) get() string {
	p.init()
	for p.validate() {
		p.l = p.l[:0]
		p.uniqSlice = p.uniqSlice[0:]
		for i := 0; i < passwordLength; i++ {
			n := rand.Intn(passwordComplexity)
			p.l = append(p.l, n)
			p.uniqSlice[n] = true
		}
	}
	for i := 0; i < passwordLength; i++ {
		tmp := p.mappingArray[p.l[i]]
		p.passwords[i] = tmp[rand.Intn(len(tmp))]
	}
	return string(p.passwords)
}

func getLdapUserPassword(password string) string {
	c := sha512_crypt.New()
	salt := sha512_crypt.GetSalt()
	s, err := c.Generate([]byte(password), salt.GenerateWRounds(saltLength, round))
	if err != nil {
		log.Fatal(err)
	}
	return ldapUserPasswordPrefix + s
}

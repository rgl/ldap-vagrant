package main

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v2"
)

const (
	host          = "ldap.example.com"
	port          = 389
	peopleGroupDn = "ou=people,dc=example,dc=com"
)

func main() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer l.Close()

	search := ldap.NewSearchRequest(
		peopleGroupDn,
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=person)",
		[]string{
			"cn",
			"mail",
		},
		nil)

	r, err := l.Search(search)
	if err != nil {
		log.Fatalf("Failed people search: %v", err)
	}

	r.PrettyPrint(0)
}

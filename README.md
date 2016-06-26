This is a [Vagrant](https://www.vagrantup.com/) Environment for a [Directory/LDAP](https://en.wikipedia.org/wiki/Lightweight_Directory_Access_Protocol) Server.

This lets you easily test your application code against a real sandboxed server.

This uses the [slapd](http://www.openldap.org/software/man.cgi?query=slapd) daemon from [OpenLDAP](http://www.openldap.org/).

LDAP is described at [RFC 4510 (Technical Specification)](https://tools.ietf.org/html/rfc4510).

Also check the [OpenLDAP Server documentation at the Ubuntu Server Guide](https://help.ubuntu.com/lts/serverguide/openldap-server.html).

# Usage

Run `vagrant up` to configure the `ldap.example.com` LDAP server environment.

Configure your system `/etc/hosts` file with the `ldap.example.com` domain:

    192.168.33.253 ldap.example.com

The environment comes pre-configured with the following entries:

    uid=alice,ou=people,dc=example,dc=com
    uid=bob,ou=people,dc=example,dc=com
    uid=carol,ou=people,dc=example,dc=com
    uid=dave,ou=people,dc=example,dc=com
    uid=eve,ou=people,dc=example,dc=com
    uid=frank,ou=people,dc=example,dc=com
    uid=grace,ou=people,dc=example,dc=com
    uid=henry,ou=people,dc=example,dc=com

To see how these were added take a look at the end of the [provision.sh](provision.sh) file.

To troubleshoot, watch the logs with `vagrant ssh` and `sudo journalctl --follow`.

# Examples

There are examples available on how to connect to LDAP programmatically (e.g. from Go). Have a look at the [examples directory](examples).

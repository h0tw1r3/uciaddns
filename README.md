# UCI AD DNS

Very simple utility to output [uci] commands for creating the [minimum required AD records] in dnsmasq on [openwrt].

## Usage

    uciaddns SERVER REALM

## Example

    ./uciaddns 192.168.1.5:53 my.home
    uci add dhcp srvhost
    uci set dhcp.@srvhost[-1].srv="_ldap._tcp.gc._msdcs.my.home"
    uci set dhcp.@srvhost[-1].target="dc0.my.home."
    uci set dhcp.@srvhost[-1].port="3268"
    uci set dhcp.@srvhost[-1].class="100"
    uci set dhcp.@srvhost[-1].weight="0"
    uci commit dhcp
    uci add dhcp srvhost
    uci set dhcp.@srvhost[-1].srv="_ldap._tcp.my.home"
    uci set dhcp.@srvhost[-1].target="dc0.my.home."
    uci set dhcp.@srvhost[-1].port="389"
    uci set dhcp.@srvhost[-1].class="100"
    uci set dhcp.@srvhost[-1].weight="0"
    uci commit dhcp
    uci add dhcp srvhost
    uci set dhcp.@srvhost[-1].srv="_ldap._tcp.pdc._msdcs.my.home"
    uci set dhcp.@srvhost[-1].target="dc0.my.home."
    uci set dhcp.@srvhost[-1].port="389"
    uci set dhcp.@srvhost[-1].class="100"
    uci set dhcp.@srvhost[-1].weight="0"
    uci commit dhcp
    uci add dhcp srvhost
    uci set dhcp.@srvhost[-1].srv="_ldap._tcp.dc._msdcs.my.home"
    uci set dhcp.@srvhost[-1].target="dc0.my.home."
    uci set dhcp.@srvhost[-1].port="389"
    uci set dhcp.@srvhost[-1].class="100"
    uci set dhcp.@srvhost[-1].weight="0"
    uci commit dhcp
    uci add dhcp srvhost
    uci set dhcp.@srvhost[-1].srv="_kerberos._tcp.dc._msdcs.my.home"
    uci set dhcp.@srvhost[-1].target="dc0.my.home."
    uci set dhcp.@srvhost[-1].port="88"
    uci set dhcp.@srvhost[-1].class="100"
    uci set dhcp.@srvhost[-1].weight="0"
    uci commit dhcp
    uci add_list dhcp.@dnsmasq[0].address="/gc._msdcs.my.home/192.168.1.5"
    uci commit dhcp
    /etc/init.d/dnsmasq restart

This is *not* everything needed for a complete AD DNS setup through [openwrt].

For example, `dc0` forward and reverse resolution must be working properly.
If you're using an openwrt variant that uses kresd, configure reverse stub and forwarding.
Here's how I did it for Turris Omnia.

    customStubZones = policy.todnames({'1.168.192.in-addr.arpa', '1.168.192.in-addr.arpa.'})
    policy.add(policy.suffix(policy.FLAGS({'NO_CACHE'}), customStubZones))
    policy.add(policy.suffix(policy.STUB({'127.0.0.1@5353'}), customStubZones))
    forwardZones = policy.todnames({'my.home'})
    policy.add(policy.suffix(policy.FLAGS({'NO_CACHE'}), forwardZones))
    policy.add(policy.suffix(policy.FORWARD({'127.0.0.1@5353'}), forwardZones)

## Why?

While working on a [go] program that queried DNS records from a Samba DNS server, I discovered a [bug in the golang dns module].

[uci]: https://openwrt.org/docs/guide-user/base-system/uci
[openwrt]: https://openwrt.org/
[minimum required AD records]: https://docs.microsoft.com/en-us/archive/blogs/servergeeks/dns-records-that-are-required-for-proper-functionality-of-active-directory
[go]: https://golang.org
[bug in the golang dns module]: https://github.com/golang/go/issues/37362

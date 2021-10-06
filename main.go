package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

// miekg/dns is here because golang native net resolve does not work with Samba
// https://github.com/golang/go/issues/37362

type req struct {
	Name string
	Type uint16
}

const uci_srvhost string = "uci add dhcp srvhost\n" +
	"uci set dhcp.@srvhost[-1].srv=\"%s\"\n" +
	"uci set dhcp.@srvhost[-1].target=\"%s\"\n" +
	"uci set dhcp.@srvhost[-1].port=\"%d\"\n" +
	"uci set dhcp.@srvhost[-1].class=\"%d\"\n" +
	"uci set dhcp.@srvhost[-1].weight=\"%d\"\n"

const uci_address string = "uci add_list dhcp.@dnsmasq[0].address=\"/%s/%s\"\n"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s SERVER REALM\n", os.Args[0])
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid arguments\n")
		os.Exit(1)
	}

	var nameserver = flag.Arg(0)
	var domain = flag.Arg(1)

	var adrecs = []req{
		{"_ldap._tcp.gc._msdcs", dns.TypeSRV},
		{"_ldap._tcp", dns.TypeSRV},
		{"_ldap._tcp.pdc._msdcs", dns.TypeSRV},
		{"_ldap._tcp.dc._msdcs", dns.TypeSRV},
		{"_kerberos._tcp.dc._msdcs", dns.TypeSRV},
		{"gc._msdcs", dns.TypeA},
	}

	c := new(dns.Client)
	c.Net = "udp"
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  true,
			Opcode:            dns.OpcodeQuery,
		},
		Question: make([]dns.Question, 1),
	}
	if op, ok := dns.StringToOpcode["QUERY"]; ok {
		m.Opcode = op
	}
	qc := uint16(dns.ClassINET)

	for _, adrec := range adrecs {
		var name = adrec.Name + "." + domain
		m.Question[0] = dns.Question{Name: dns.Fqdn(name), Qtype: adrec.Type, Qclass: qc}
		m.Id = dns.Id()

		r, _, err := c.Exchange(m, nameserver)
		if err != nil {
			panic(err)
		}
		if r != nil && r.Rcode != dns.RcodeSuccess {
			fmt.Printf("%v", dns.RcodeToString[r.Rcode])
			return
		}

		if len(r.Answer) > 0 {
			for _, record := range r.Answer {
				switch addr := record.(type) {
				case *dns.SRV:
					fmt.Printf(uci_srvhost, name, addr.Target, addr.Port, addr.Weight, addr.Priority)
				case *dns.A:
					fmt.Printf(uci_address, name, addr.A)
				default:
					fmt.Fprintf(os.Stderr, "invalid SRV response for record %s", record)
					os.Exit(2)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "No answer for %s", name)
			os.Exit(2)
		}

		fmt.Print("uci commit dhcp\n")
	}
}

// Package netdb provides a Go interface for the protoent and servent
// structures as defined in netdb.h
//
// A pure Go implementation is used by parsing /etc/protocols and
// /etc/services
package netdb

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Protoent struct {
	Name    string
	Aliases []string
	Number  int
}

type Servent struct {
	Name     string
	Aliases  []string
	Port     int
	Protocol string
}


// These variables get populated from /etc/protocols and /etc/services
// respectively.
var (
	Protocols []Protoent
	Services []Servent
)

func init() {
	// Load protocols
	data, err := ioutil.ReadFile("/etc/protocols")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, "#", 2)
		fields := strings.Fields(split[0])
		if len(fields) < 2 {
			continue
		}

		num, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			panic(err)
		}

		Protocols = append(Protocols, Protoent{
			Name:    fields[0],
			Aliases: fields[2:],
			Number:  int(num),
		})
	}

	// Load services
	data, err = ioutil.ReadFile("/etc/services")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, "#", 2)
		fields := strings.Fields(split[0])
		if len(fields) < 2 {
			continue
		}

		name := fields[0]
		portproto := strings.SplitN(fields[1], "/", 2)
		port, err := strconv.ParseInt(portproto[0], 10, 32)
		if err != nil {
			panic(err)
		}

		proto := portproto[1]
		aliases := fields[2:]

		Services = append(Services, Servent{
			Name:     name,
			Aliases:  aliases,
			Port:     int(port),
			Protocol: proto,
		})
	}
}

// Equal checks if two Protoents are the same, which is the case if
// their protocol numbers are identical.
func (this Protoent) Equal(other Protoent) bool {
	return this.Number == other.Number
}

// GetProtoByNumber returns the Protoent for the correspondent
// protocol number.
func GetProtoByNumber(num int) (protoent Protoent, ok bool) {
	for _, protoent := range Protocols {
		if protoent.Number == num {
			return protoent, true
		}
	}
	return Protoent{}, false
}

// GetProtoByName returns the Protoent whose name or any of its
// aliases matches the argument.
func GetProtoByName(name string) (protoent Protoent, ok bool) {
	for _, protoent := range Protocols {
		if protoent.Name == name {
			return protoent, true
		}

		for _, alias := range protoent.Aliases {
			if alias == name {
				return protoent, true
			}
		}
	}

	return Protoent{}, false
}

// GetServByName returns the Servent for a given service name and
// protocol name. If the protocol name is empty, the first service
// matching the service name is returned.
func GetServByName(name, protocol string) (servent Servent, ok bool) {
	for _, servent := range Services {
		if servent.Protocol != protocol && protocol != "" {
			continue
		}

		if servent.Name == name {
			return servent, true
		}

		for _, alias := range servent.Aliases {
			if alias == name {
				return servent, true
			}
		}
	}

	return Servent{}, false
}

// GetServByPort returns the Servent for a given port number and
// protocol name. If the protocol name is empty, the first service
// matching the port number is returned.
func GetServByPort(port int, protocol string) (Servent, bool) {
	for _, servent := range Services {
		if servent.Port == port && (servent.Protocol == protocol || protocol == "") {
			return servent, true
		}
	}

	return Servent{}, false
}

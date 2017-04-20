package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0] + " <file>")
		fmt.Println("Where <file> is pd interfaces file, created after show-interfaces command")
		os.Exit(1)
	}

	//Open files for input and output

	filein, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer filein.Close()

	fileout, err := os.Create(os.Args[1] + "-" + "new" + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	defer fileout.Close()

	//Creating regexp to match invalid records in xml file

	notXMLStringRxp := regexp.MustCompile(`^[A-z].*$`)
	emptyStringRxp := regexp.MustCompile(`^$`)

	//Fixing xml

	//Before start we should append file with dumb root node
	fileout.WriteString("<DumbRootNode>\n")

	scanner := bufio.NewScanner(filein)
	for scanner.Scan() {
		if !notXMLStringRxp.MatchString(scanner.Text()) && !emptyStringRxp.MatchString(scanner.Text()) {
			fileout.WriteString(scanner.Text() + "\n")
		}
	}
	//Creting closing dumb root node
	fileout.WriteString("</DumbRootNode>\n")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fileout.Close()
	//

	//Parsing interfaces.xml

	type InterfaceAddress struct {
		AddressType string `xml:"addressType"`
		Address     string `xml:"address"`
	}

	type PhysicalInterface struct {
		Unit                    string           `xml:"unit"`
		PhysicalInterfaceName   string           `xml:"physicalInterfaceName"`
		OsInterfaceName         string           `xml:"osInterfaceName"`
		ConfiguredName          string           `xml:"configuredName"`
		MacAddress              string           `xml:"macAddress"`
		NetworkInterfaceAddress InterfaceAddress `xml:"interfaceAddress"`
		UseDHCP                 string           `xml:"useDHCP"`
		AllowAdmin              string           `xml:"allowAdmin"`
		AutoNegotiate           string           `xml:"autoNegotiate"`
		Speed                   string           `xml:"speed"`
		DuplexType              string           `xml:"duplexType"`
		Mtu                     string           `xml:"mtu"`
	}

	type VlanInterface struct {
		Unit                    string           `xml:"unit"`
		PhysicalInterface       string           `xml:"physicalInterface"`
		ConfiguredName          string           `xml:"configuredName"`
		VlanID                  string           `xml:"vlanId"`
		NetworkInterfaceAddress InterfaceAddress `xml:"interfaceAddress"`
		UseDHCP                 string           `xml:"useDHCP"`
		AllowAdmin              string           `xml:"allowAdmin"`
		Mtu                     string           `xml:"mtu"`
	}

	type EndPoint struct {
		AddressType string `xml:"addressType"`
		Address     string `xml:"address"`
	}

	type TunnelInterface struct {
		Unit              string   `xml:"unit"`
		PhysicalInterface string   `xml:"physicalInterface"`
		ConfiguredName    string   `xml:"configuredName"`
		LocalEndPoint     EndPoint `xml:"localEndPoint"`
		RemoteEndPoint    EndPoint `xml:"remoteEndPoint"`
		Mtu               string   `xml:"mtu"`
		TunnelNumber      string   `xml:"tunnelNumber"`
	}

	xmlFile, err := os.Open(os.Args[1] + "-" + "new" + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	xmlDecoder := xml.NewDecoder(xmlFile)

	for {
		tok, _ := xmlDecoder.Token()
		if tok == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := tok.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			// ...and its name is "page"
			if se.Name.Local == "PhysicalInterface" {
				var pi PhysicalInterface
				// decode a whole chunk of following XML into the
				// variable  which is a PhysicalInterface (see above)
				xmlDecoder.DecodeElement(&pi, &se)
				// Do some stuff with the page.
				//p.Title = CanonicalizeTitle(p.Title)
				fmt.Println("Phys" + ";" + "\"" + pi.ConfiguredName + "\"" + ";" + pi.NetworkInterfaceAddress.Address + ";" + pi.Mtu + ";" + pi.MacAddress)
			}

			if se.Name.Local == "VlanInterface" {
				var vi VlanInterface
				xmlDecoder.DecodeElement(&vi, &se)
				fmt.Println("Vlan" + ";" + "\"" + vi.ConfiguredName + "\"" + ";" + vi.NetworkInterfaceAddress.Address + ";" + vi.Mtu + ";" + vi.VlanID)
			}

			if se.Name.Local == "TunnelInterface" {
				var ti TunnelInterface
				xmlDecoder.DecodeElement(&ti, &se)
				fmt.Println("Tun" + ";" + "\"" + ti.ConfiguredName + "\"" + ";" + ti.LocalEndPoint.Address + ";" + "\"" + ti.RemoteEndPoint.Address + "\"" + ";" + ti.Mtu)
			}
		}
	}
}

package main

import "testing"

func TestValidateDestination(t *testing.T) {

	type destinationParts struct {
		hostName string
		port     string
		macAddr  string
	}
	type table struct {
		input  destinationParts
		output error
	}
	var inputs []table
	inputs = append(inputs, table{input: destinationParts{hostName: "atul", port: "9999", macAddr: "00"}, output: MacAddrLengthWrongErr})
	inputs = append(inputs, table{input: destinationParts{hostName: "atul", port: "9999", macAddr: "00000000ZZZZ"}, output: MacAddrFormatWrongErr})
	inputs = append(inputs, table{input: destinationParts{hostName: "atul", port: "999", macAddr: "00000000ffff"}, output: PortFormatWrongErr})
	inputs = append(inputs, table{input: destinationParts{hostName: "atul", port: "999a", macAddr: "00000000ffff"}, output: PortFormatWrongErr})
	inputs = append(inputs, table{input: destinationParts{hostName: "", port: "9999", macAddr: "00000000ffff"}, output: HostNameFormatWrongErr})

	for _, test := range inputs {
		err := validateDestination(test.input.hostName, test.input.port, test.input.macAddr)
		if err != test.output {
			t.Errorf("ValidateDestination(%v). Got - %v. Expected - %v", test.input, err, test.output)
		}

	}
}

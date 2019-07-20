package main

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

const quoter = `{
    "code": "1A",
    "sumInsured": 12000,
    "dateOfBirth" : "14/01/1977",
    "premium": 3000,
    "parties": [
        {
            "firstName": "Bhavesh",
            "lastName": "Yadav",
            "gender": "MALE",
            "dateOfBirth": "14/01/1977",
            "mobileNumber": "1234567890",
            "email": "primary@gmail.com",
            "panNumber": "AJBDD12345G",
            "aadhaar": 123456789012,
            "addressLine1": "ketaki",
            "addressLine2": "maneklal",
            "addressLine3": "Ghatkopar",
            "city": "mumbai",
            "pinCode": 400086,
            "latitude": 123223232,
            "longitude": 12345643,
            "relationship": "self",
            "isPrimary": true
        },
        {
            "firstName": "Usha",
            "lastName": "Patel",
            "gender": "FEMALE",
            "dateOfBirth": "14/01/1977",
            "mobileNumber": "1234567890",
            "email": "nominiee@gmail.com",
            "panNumber": "AJBDD12345G",
            "aadhaar": 123456789012,
            "addressLine1": "ketaki",
            "addressLine2": "maneklal",
            "addressLine3": "Ghatkopar",
            "city": "mumbai",
            "pincode": 400086,
            "latitude": 123223232,
            "Longitude": 12345643,
            "relationship": "self",
            "IsPrimary": false
        }
    ]
}`

func TestQuoteJsonMarshall(t *testing.T) {
	body := strings.NewReader(quoter)
	q := &quotereq{}
	err := json.NewDecoder(body).Decode(q)
	if err == nil {
		t.Errorf("error while unmarshalling %v", err)
	}
	log.Println(q)
}

func TestValidateReq(t *testing.T) {
	var qreq = quotereq{Code: "1A", SumInsured: 123, DateOfBirth: "14/01/1977"}
	var parties []party
	for i := range []int{1, 2} {
		var party = party{FirstName: "prashant", Email: "pras.p.in@gmail.com", PinCode: int32(i * 400086), MobileNumber: "asasaa"}
		parties = append(parties, party)
	}
	qreq.Parties = parties
	errors := validateReq(qreq)
	if errors != nil {
		t.Errorf("error in request %v", errors)
	}
}

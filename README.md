# quote

db.quotes.counter.insert({"_id" : "quoteid" , "value": 0 })
db.quotes.counter.find({}).pretty()

curl -i -X POST \
  http://172.17.0.12:8000/api/v1/healths/quotes \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: ef137bfc-bfa8-4c0a-bcc3-0dd6c5607075' \
  -H 'cache-control: no-cache' \
  -d '{
    "code": "1A",
    "SumInsured": 12000,
    "Premium": 3000,
    "parties": [
        {
            "firstName": "Bhavesh",
            "lastName": "Yadav",
            "gender": "MALE",
            "dataOfBirth": "14/01/1977",
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
            "dataOfBirth": "14/01/1977",
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
}'
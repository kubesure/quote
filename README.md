# quote

#### biz design

api creates quote for insured's risk. Creates parties, insured and nominee calling party grpc service. Parties are bound to quote and qoute number genereted.

#### components

Mongodb v4, GRPC, Golang 

#### Dev setup and test

1. create db quote and party 
    ```
       use quotes
       db.quotes.counter.insert({"_id" : "quoteid" , "value": 0 })
       db.quotes.counter.find({}).pretty()
    ```
2. Run party and quote 
   ```
       go run ../party/party.go
       go run quote.go 
   ```

3. Run curl to create quote

```
curl -i -X POST \
  http://localhost:8000/api/v1/healths/quotes   \
  -d '{                                         \
    "code": "1A",                               \
    "SumInsured": 12000,                        \
    "Premium": 3000,                            \
    "parties": [                                \
        {                                       \
            "firstName": "Bhavesh",             \
            "lastName": "Yadav",                \ 
            "gender": "MALE",                   \ 
            "dataOfBirth": "14/01/1977",        \ 
            "mobileNumber": "1234567890",       \
            "email": "primary@gmail.com",       \
            "panNumber": "AJBDD12345G",         \
            "aadhaar": 123456789012,            \
            "addressLine1": "ketaki",           \ 
            "addressLine2": "maneklal",         \
            "addressLine3": "Ghatkopar",        \ 
            "city": "mumbai",                   \
            "pinCode": 400086,                  \    
            "latitude": 123223232,              \
            "longitude": 12345643,              \
            "relationship": "self",             \
            "isPrimary": true                   \ 
        },                                      \
        {                                       \
            "firstName": "Usha",                \
            "lastName": "Patel",                \
            "gender": "FEMALE",                 \
            "dataOfBirth": "14/01/1977",        \
            "mobileNumber": "1234567890",       \
            "email": "nominiee@gmail.com",      \ 
            "panNumber": "AJBDD12345G",         \
            "aadhaar": 123456789012,            \
            "addressLine1": "ketaki",           \
            "addressLine2": "maneklal",         \
            "addressLine3": "Ghatkopar",        \
            "city": "mumbai",                   \
            "pincode": 400086,                  \
            "latitude": 123223232,              \
            "Longitude": 12345643,              \
            "relationship": "self",             \
            "IsPrimary": false                  \
        }                                       \
    ]                                           \
}                                               \'
```

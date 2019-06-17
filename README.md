# quote

#### biz design

api creates quote for insured's risk. Creates parties, insured and nominee calling party grpc service. Parties are bound to quote and qoute number generated.

#### components

Mongodb v4, GRPC, Golang

#### Dev setup and test

1. start mongo in replication mode
```
mongod --replSet=rs0 --bind_ip="0.0.0.0" --smallfiles --noprealloc --port="27017" --dbpath=$HOME/data/dbrs0
mongod --replSet=rs0 --bind_ip="0.0.0.0" --smallfiles --noprealloc --port="37017" --dbpath=$HOME/data/dbrs1
mongod --replSet=rs0 --bind_ip="0.0.0.0" --smallfiles --noprealloc --port="47017" --dbpath=$HOME/data/dbrs2

rs.initiate({ _id: "rs0", members:[ 
    { _id: 0, host: "localhost:27017" },
    { _id: 1, host: "localhost:37017" },
    { _id: 2, host: "localhost:47017" },
]});

rs.conf()
```

2. create db quote and party (follow instruction in party service) in mongodb primary node
    ```
       use quotes
       db.counter.insert({"_id" : "quoteid" , "value": 0 })
       db.counter.find({}).pretty()
       db.quote.find({}).pretty()
    ```
3. Run party and quote
   ```
       go run ../party/party.go
       go run quote.go
   ```

4. Run curl to create quote

```
curl -i -X POST http://localhost:8000/api/v1/healths/quotes  -H 'Content-Type: application/json' -d '{
    "code": "1A",
    "SumInsured": 12000,
    "dateOfBirth" : "14/01/1977",
    "Premium": 3000,
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
}'
```
## Pod 

1. Create and configure mongo in k8s
 
    ```
    alias k=kubectl
    complete -F __start_kubectl k

    k apply -f config/mongo.yaml
    k exec mongo-quote-0 -it mongo

    rs.initiate({ _id: "rs0", members:[ 
        { _id: 0, host: "mongo-quote-0.mongoquotesvc:27017" },
        { _id: 1, host: "mongo-quote-1.mongoquotesvc:27017" },
        { _id: 2, host: "mongo-quote-2.mongoquotesvc:27017" },
    ] });

    rs.conf()
    ```

    create document store follow step 2 in Dev setup 

2. Create quote docker image
    
    ```
    docker build . -t quote:v1
    k apply config/quote.yaml
    k get po -o wide
    ```

3.  Apply Quote to k8s

    ```
        k apply -f config/quote.yaml
        k get po -o wide 
    ```    

    curl test. Follow step 4.

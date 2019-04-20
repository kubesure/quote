package main

import (
	//"context"
	"encoding/json"
	"fmt"
	api "github.com/kubesure/party/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var mongoquotesvc = os.Getenv("mongoquotesvc")

const (
	address = "localhost:50051"
)

type quotereq struct {
	Code        string  `json:"code" bson:"code"`
	SumInsured  int32   `json:"sumInsured" bson:"sumInsured"`
	DateOfBirth string  `json:"dateOfBirth" bson:"dateOfBirth"`
	Premium     int32   `json:"premium" bson:"premium"`
	Parties     []party `json:"parties" bson:"parties"`
}

type party struct {
	FirstName    string  `json:"firstName" bson:"firstName"`
	LastName     string  `json:"lastName" bson:"lastName"`
	Gender       string  `json:"gender" bson:"gender"`
	DataOfBirth  string  `json:"dateOfBirth" bson:"dateOfBirth"`
	MobileNumber string  `json:"mobileNumber" bson:"mobileNumber"`
	Email        string  `json:"email" bson:"email"`
	PanNumber    string  `json:"panNumber" bson:"panNumber"`
	Aadhaar      int64   `json:"aadhaar" bson:"aadhaar"`
	AddressLine1 string  `json:"addressLine1" bson:"addressLine1"`
	AddressLine2 string  `json:"addressLine2" bson:"addressLine2"`
	AddressLine3 string  `json:"addressLine3" bson:"addressLine3"`
	City         string  `json:"city" bson:"city"`
	PinCode      int32   `bson:"pinCode" bson:"pinCode"`
	Latitude     float64 `json:"latitude" bson:"latitude"`
	Longitude    float64 `json:"longitude" bson:"latitude"`
	Relationship string  `json:"relationship" bson:"relationship"`
	IsPrimary    bool    `json:"isPrimary" bson:"isPrimary"`
}

type quoteres struct {
	QuoteNumber int64 `json:"quoteNumber"`
}

func main() {
	log.Println("quote api starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healths/quotes", quote)
	srv := http.Server{Addr: ":8080", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Print("shutting down receipt server...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
	//log.Fatal(http.ListenAndServe(":8000", mux))
}

func quote(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	q, merr := marshallReq(string(body))
	r, serr := save(q)

	if merr != nil {
		log.Println(merr)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else if serr != nil {
		log.Println(serr)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		data, _ := json.Marshal(r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", data)
	}
}

func save(q *quotereq) (*quoteres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoquotesvc+":27017"))
	errping := client.Ping(ctx, nil)

	if errping != nil {
		return nil, errping
	}

	collection := client.Database("quotes").Collection("quote")
	id, errSeq := nextcounter(client)
	if errSeq != nil {
		log.Println("err seq ", errSeq)
		return nil, errSeq
	}

	var parties []bson.D

	for _, p := range q.Parties {
		id, err := saveparty(&p)
		if err != nil {
			return nil, err
		}
		d := bson.D{
			{"partyId", id},
			{"relationship", p.Relationship},
			{"isPrimary", p.IsPrimary},
		}
		parties = append(parties, d)
	}

	quote := bson.M{
		"quoteNumber": id, "code": q.Code, "sumAssured": q.SumInsured, "premium": q.Premium,
		"parties": parties, "createdDate": time.Now().String(),
	}
	_, errcol := collection.InsertOne(context.Background(), quote)

	if errcol != nil {
		log.Println("errcol", errcol)
		return nil, errcol
	}

	res := quoteres{QuoteNumber: id}
	return &res, nil
}

func saveparty(qp *party) (int64, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {

		return 0, err
	}
	defer conn.Close()
	client := api.NewPartyServiceClient(conn)

	var p api.Party
	p.Aadhaar = qp.Aadhaar
	p.AddressLine1 = qp.AddressLine1
	p.AddressLine2 = qp.AddressLine2
	p.AddressLine3 = qp.AddressLine3
	p.City = qp.City
	p.DataOfBirth = qp.DataOfBirth
	p.Email = qp.Email
	p.FirstName = qp.FirstName
	p.LastName = qp.LastName
	p.Latitude = qp.Latitude
	p.Longitude = qp.Longitude
	if qp.Gender == "MALE" {
		p.Gender = api.Party_MALE
	}
	if qp.Gender == "FEMALE" {
		p.Gender = api.Party_FEMALE
	}

	var phones []*api.Party_PhoneNumber
	phone := api.Party_PhoneNumber{Number: qp.MobileNumber}
	phone.Type = api.Party_MOBILE
	phones = append(phones, &phone)
	p.Phones = phones

	req := api.PartyRequest{Party: &p}
	party, err := client.CreateParty(context.Background(), &req)
	if err != nil {
		return 0, err
	}
	return party.Id, nil
}

func nextcounter(c *mongo.Client) (int64, error) {
	collection := c.Database("quotes").Collection("counter")
	filter := bson.M{"_id": "quoteid"}
	update := bson.M{"$inc": bson.M{"value": 1}}
	aft := options.After
	opt := options.FindOneAndUpdateOptions{Upsert: new(bool), ReturnDocument: &aft}
	result := collection.FindOneAndUpdate(context.Background(), filter, update, &opt)
	type record struct {
		Quoteid string `bson:"quoteid"`
		Value   int64  `bson:"value"`
	}
	var data record
	err := result.Decode(&data)
	if err != nil {
		return 0, err
	}
	return data.Value, nil
}

func marshallReq(data string) (*quotereq, error) {
	var q quotereq
	err := json.Unmarshal([]byte(data), &q)
	if err != nil {

		return nil, err
	}
	return &q, nil
}

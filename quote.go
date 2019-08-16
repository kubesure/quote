package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	api "github.com/kubesure/party/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/go-playground/validator.v9/translations/en"
)

var mongoquotesvc = os.Getenv("mongoquotesvc")
var partysvc = os.Getenv("partysvc")

var validate *validator.Validate

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

}

type quotereq struct {
	Code        string  `json:"code" bson:"code" validate:"required"`
	SumInsured  int32   `json:"sumInsured" bson:"sumInsured" validate:"required"`
	DateOfBirth string  `json:"dateOfBirth" bson:"dateOfBirth" validate:"required"`
	Premium     int32   `json:"premium" bson:"premium" validate:"required"`
	Parties     []party `json:"parties" bson:"parties" validate:"required,len=2,dive"`
}

type party struct {
	FirstName    string  `json:"firstName" bson:"firstName" validate:"required,gte=3,lt=20"`
	LastName     string  `json:"lastName" bson:"lastName" validate:"required,gt=3,lt=20"`
	Gender       string  `json:"gender" bson:"gender" validate:"required"`
	DataOfBirth  string  `json:"dateOfBirth" bson:"dateOfBirth" validate:"required"`
	MobileNumber string  `json:"mobileNumber" bson:"mobileNumber" validate:"required,len=10,numeric"`
	Email        string  `json:"email" bson:"email" validate:"required,email"`
	PanNumber    string  `json:"panNumber" bson:"panNumber" validate:"required,min=10"`
	Aadhaar      int64   `json:"aadhaar" bson:"aadhaar" validate:"required,min=12"`
	AddressLine1 string  `json:"addressLine1" bson:"addressLine1" validate:"required"`
	AddressLine2 string  `json:"addressLine2" bson:"addressLine2" validate:"required"`
	AddressLine3 string  `json:"addressLine3" bson:"addressLine3" validate:"required"`
	City         string  `json:"city" bson:"city" validate:"required"`
	PinCode      int32   `bson:"pinCode" bson:"pinCode" validate:"required"`
	Latitude     float64 `json:"latitude" bson:"latitude" validate:"required,latitude"`
	Longitude    float64 `json:"longitude" bson:"latitude" validate:"required,longitude"`
	Relationship string  `json:"relationship" bson:"relationship" validate:"required"`
	IsPrimary    bool    `json:"isPrimary" bson:"isPrimary"`
}

type quoteres struct {
	QuoteNumber int `json:"quoteNumber"`
}

type errorresponse struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func main() {
	log.Debug("quote api starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healths/quotes", quote)
	mux.HandleFunc("/isready", isReady)
	srv := http.Server{Addr: ":8000", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Debug("shutting down quote server...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func isReady(w http.ResponseWriter, req *http.Request) {
	client, errping := conn()
	defer client.Disconnect(context.Background())
	if errping != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	coll := client.Database("quotes").Collection("quote")
	if coll == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func quote(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	q, merr := marshallReq(body)
	if merr != nil {
		log.Error(merr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	if err := validateHeader(w, req); err != nil {
		return
	}

	//fix pointer deref
	if err := validateReq(*q); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-type", "application/json")

	if merr != nil {
		log.Error(merr)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if r, serr := save(q); serr != nil {
			log.Error(serr)
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			data, _ := json.Marshal(r)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "%s", data)
		}
	}
}

func validateHeader(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		log.Debug("invalid method ", req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Invalid method %s", req.Method)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		log.Debug("invalid content type ", req.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Invalid content-type require %s", "application/json")
	}
	return nil
}

func validateReq(q quotereq) map[string][]string {
	validate := validator.New()
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	_ = en.RegisterDefaultTranslations(validate, trans)

	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	errv := validate.Struct(q)
	errors := make(map[string][]string)

	_, ok := errv.(*validator.InvalidValidationError)

	if !ok {
		return nil
	}

	for _, e := range errv.(validator.ValidationErrors) {
		if val, ok := errors[e.StructField()]; !ok {
			var err = make([]string, 0)
			err = append(err, e.Translate(trans))
			errors[e.StructField()] = err
		} else {
			val = append(val, e.Translate(trans))
			errors[e.StructField()] = val
		}
	}
	log.Println(errors)
	return errors
}

func save(q *quotereq) (*quoteres, error) {

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

	client, errping := conn()

	if errping != nil {
		return nil, errping
	}
	defer client.Disconnect(context.Background())

	id, errSeq := nextcounter(client)
	if errSeq != nil {
		return nil, errSeq
	}

	quote := bson.M{
		"quoteNumber": id, "code": q.Code, "sumAssured": q.SumInsured, "dateOfBirth": q.DateOfBirth,
		"premium": q.Premium, "parties": parties, "createdDate": time.Now().String(),
	}
	collection := client.Database("quotes").Collection("quote")
	_, errcol := collection.InsertOne(context.Background(), quote)

	if errcol != nil {
		return nil, errcol
	}

	res := quoteres{QuoteNumber: id}
	return &res, nil
}

func saveparty(qp *party) (int64, error) {
	conn, err := grpc.Dial(partysvc+":50051", grpc.WithInsecure())
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
	p.PinCode = qp.PinCode
	p.PanNumber = qp.PanNumber
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

func nextcounter(c *mongo.Client) (int, error) {
	collection := c.Database("quotes").Collection("counter")
	filter := bson.M{"_id": "quoteid"}
	update := bson.M{"$inc": bson.M{"value": 1}}
	aft := options.After
	opt := options.FindOneAndUpdateOptions{Upsert: new(bool), ReturnDocument: &aft}
	result := collection.FindOneAndUpdate(context.Background(), filter, update, &opt)
	type record struct {
		Quoteid string `bson:"quoteid"`
		Value   int    `bson:"value"`
	}
	var data record
	err := result.Decode(&data)
	if err != nil {
		return 0, err
	}
	return data.Value, nil
}

func marshallReq(data []byte) (*quotereq, error) {
	var q quotereq
	err := json.Unmarshal(data, &q)
	if err != nil {

		return nil, err
	}
	return &q, nil
}

func conn() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uri := "mongodb://" + mongoquotesvc + "/?replicaSet=rs0"
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	errping := client.Ping(ctx, nil)
	return client, errping
}

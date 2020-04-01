package main

import (
	"log"
	"fmt"
	"time"
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

type loc struct {
	Name   string    `json:"type,omitempty"`
}

type relation struct {
	Type   string    `json:"type,omitempty"`
}

// If omitempty is not set, then edges with empty values (0 for int/float, "" for string, false
// for bool) would be created for values not specified explicitly.

type Person struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Relationship  relation   `json:"relation,omitempty"`
}

type Dream struct {
	Date 		*time.Time `json:"dob,omitempty"`
	Actors 		[]Person   `json:"friend,omitempty"`
	Location 	loc        `json:"loc,omitempty"`
	Intimity	int        `json:"intimity,omitempty"`
	Text		string	   `json:"text,omitempty"`
}

type CancelFunc func()

func getDgraphClient() (*dgo.Dgraph, CancelFunc) {
	//Listens on 9080 by default
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

func addSchema(dg *dgo.Dgraph) error {
	op := &api.Operation{}
	op.Schema = `
	name: string @index(exact,term) .
	location: [uid] .
	date: datetime .
	actors: [uid] .
	location: [uid] .
	relationship: [uid] .
	intimity: int  .
	text: string .
	type: string .

	type Dream {
	date
	actors
	location
	intimity
	text
	}

	type Person {
	name
	relationship
	}

	type Loc {
	name
	}

	type Relation {
	type
	}

	`

	err := dg.Alter(context.Background(), op)
	if err != nil {
		return err
	}
	return nil
}

func insertData(ctx context.Context, dg *dgo.Dgraph) (*api.Response, error) {
	//year int, month Month, day, hour, min, sec, nsec int, loc *Location
	dreamDate := time.Date(2020, 04, 01, 07, 0, 0, 0, time.UTC)
	d := Dream{
		Date: &dreamDate,
		Actors: []Person{{
			Uid:     "_:bob",
			Name: "Bob",
			Relationship:  relation{
				Type: "friend",
			},
		}, {
			Uid:     "_:alice",
			Name: "Alice",
			Relationship:  relation{
				Type: "girlfriend",
			},
		}},
		Location: loc{
			Name: "Home",
		},
		Intimity: 1,
		Text: "There was a witch with a dick",
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	db, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	mu.SetJson = db
	res, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return res, nil
}
 
func main() {
	dg, cancel := getDgraphClient()
	defer cancel()
	/*
	if err := addSchema(dg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Schema added successfully")

	ctx := context.Background()
	res, err := insertData(ctx, dg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)
	*/
	resp, err := dg.NewTxn().Query(context.Background(), `{
		bobActor(func: eq(name, "Bob")) {
		  uid
		  name
		  relationship
		}
	  }`)
		  
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", resp.Json)
}
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/romana/rlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getOrderByOrderID(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getOrderByOrderID() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	var jx []byte
	var order ORDER

	redisC := REDISCLIENT.Get(r.URL.Path)

	if redisC.Err() != redis.Nil {

		jx = []byte(redisC.Val())
		mapBytes(w, r, &order, jx)

	} else {

		pth := strings.Split(r.URL.Path, "/")
		oid := pth[len(pth)-1]

		dbcol := getAccessToken(r) + OrdersExtension

		var opts options.FindOptions

		results := findMongoDocument(ExternalDB, dbcol, bson.M{"OrderID": oid}, &opts)

		if len(results) != 1 {
			respondWith(w, r, nil, OrderNotFoundMessage, nil, http.StatusNotFound, false)
			return
		}

		mapDocument(w, r, &order, results[0])

		jx = mapToBytes(w, r, results[0])

		REDISCLIENT.Set(r.URL.Path, jx, 0)

	}

	respondWith(w, r, nil, OrderFoundMessage, order, http.StatusOK, true)

}

func getOrderByCustomerID(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("getOrderByCustomerID() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	var jx []byte
	var order ORDER

	redisC := REDISCLIENT.Get(r.URL.Path)

	if redisC.Err() != redis.Nil {

		jx = []byte(redisC.Val())
		mapBytes(w, r, &order, jx)

	} else {

		pth := strings.Split(r.URL.Path, "/")
		cid := pth[len(pth)-1]

		dbcol := getAccessToken(r) + OrdersExtension

		var opts options.FindOptions

		results := findMongoDocument(ExternalDB, dbcol, bson.M{"CustomerID": cid}, &opts)

		if len(results) != 1 {
			respondWith(w, r, nil, OrderNotFoundMessage, nil, http.StatusNotFound, false)
			return
		}

		mapDocument(w, r, &order, results[0])

		jx = mapToBytes(w, r, results[0])

		REDISCLIENT.Set(r.URL.Path, jx, 0)

	}

	respondWith(w, r, nil, OrderFoundMessage, order, http.StatusOK, true)

}

func postOrder(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("postOrder() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	var order ORDER

	if !mapInput(w, r, &order) {
		return
	}

	order.OrderCreationDate = time.Now().UnixNano()

	if order.OrderID == "" {
		order.OrderID = uuid.New().String()
	}

	dbcol := getAccessToken(r) + OrdersExtension

	insertMongoDocument(ExternalDB, dbcol, order)

	respondWith(w, r, nil, OrderCreatedMessage, order, http.StatusCreated, true)

}

func putOrder(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("putOrder() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	pth := strings.Split(r.URL.Path, "/")
	oid := pth[len(pth)-1]

	var order ORDER

	if !mapInput(w, r, &order) {
		return
	}

	order.OrderCreationDate = time.Now().UnixNano()
	order.OrderID = oid

	dbcol := getAccessToken(r) + OrdersExtension

	result := updateMongoDocument(ExternalDB, dbcol, bson.M{"OrderID": order.OrderID}, bson.M{"$set": order})

	if result[0] == 1 && result[1] == 1 {

		REDISCLIENT.Del(r.URL.Path)
		respondWith(w, r, nil, OrderUpdatedMessage, order, http.StatusAccepted, true)

	} else if result[0] == 1 && result[1] == 0 {

		respondWith(w, r, nil, OrderNotUpdatedMessage, nil, http.StatusNotModified, false)

	} else if result[0] == 0 && result[1] == 0 {

		respondWith(w, r, nil, OrderNotFoundMessage, nil, http.StatusNotModified, false)

	}

}

func deleteOrder(w http.ResponseWriter, r *http.Request) {

	rlog.Debug("deleteOrder() handle function invoked ...")

	if !pre(w, r) {
		return
	}

	dbcol := getAccessToken(r) + OrdersExtension

	pth := strings.Split(r.URL.Path, "/")
	oid := pth[len(pth)-1]

	var opts options.FindOptions

	results := findMongoDocument(ExternalDB, dbcol, bson.M{"OrderID": oid}, &opts)

	if len(results) != 1 {
		respondWith(w, r, nil, OrderNotFoundMessage, nil, http.StatusNotFound, false)
		return
	}

	var order ORDER

	mapDocument(w, r, &order, results[0])

	if deleteMongoDocument(ExternalDB, dbcol, bson.M{"OrderID": oid}) == 1 {

		REDISCLIENT.Del(r.URL.Path)
		respondWith(w, r, nil, OrderDeletedMessage, nil, http.StatusOK, true)

	} else {

		respondWith(w, r, nil, OrderNotFoundMessage, nil, http.StatusNotModified, false)

	}

}

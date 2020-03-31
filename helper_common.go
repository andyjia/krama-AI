package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/romana/rlog"
	"go.mongodb.org/mongo-driver/bson"
)

func isValidJSON(s string) bool {

	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func isValidURL(toTest string) bool {

	_, err := url.ParseRequestURI(toTest)

	return !(err != nil)

}

func isValidEmail(toTest string) bool {

	emailRegularExpression := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return emailRegularExpression.MatchString(toTest)

}

func isValidAttributeKey(potentialAttributeKey string) bool {

	attributeKey := regexp.MustCompile("^[A-Za-z0-9]+([-_ ]{1}[A-Za-z0-9]+)*$")

	return attributeKey.MatchString(potentialAttributeKey)

}

func hashString(Txt string) string {

	h := sha256.New()
	h.Write([]byte(Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x", bs))
	return sh

}

func fileExists(fileName string) bool {

	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true

}

func typeof(value interface{}) string {

	return reflect.TypeOf(value).String()

}

func areCoreServicesUp() bool {

	rlog.Debug("areCoreServicesUp() function invoked ...")

	return pingMongoDB(true) && pingES(true) && pingRedis(true)

}

func resetProductCacheKeys(p *PRODUCT, pg *PRODUCTGROUP) {

	if p != nil {
		REDISCLIENT.Del(CatalogPath + "/products/" + p.Sku)
		REDISCLIENT.Del(CatalogPath + "/productgroups/" + p.GroupID)
	} else if pg != nil {
		REDISCLIENT.Del(CatalogPath + "/productgroups/" + pg.GroupID)
		for _, sku := range pg.Skus {
			REDISCLIENT.Del(CatalogPath + "/products/" + sku)
		}
	}

}

func resetCustomerCacheKeys(customer *CUSTOMER) {

	if customer != nil {
		REDISCLIENT.Del(CustomersPath + "/customers/" + customer.CustomerID)
	}

}

func mapInput(w http.ResponseWriter, r *http.Request, object interface{}) bool {

	err := json.NewDecoder(r.Body).Decode(&object)

	if err != nil {

		respondWith(w, r, err, HTTPBadRequestMessage, nil, http.StatusBadRequest, false)
		return false

	}

	return true

}

func mapDocument(w http.ResponseWriter, r *http.Request, object interface{}, document interface{}) {

	j, err0 := bson.MarshalExtJSON(document, false, false)

	if err0 != nil {
		respondWith(w, r, err0, HTTPInternalServerErrorMessage, nil, http.StatusInternalServerError, false)
		return
	}

	err1 := json.Unmarshal([]byte(j), &object)

	if err1 != nil {
		respondWith(w, r, err1, HTTPInternalServerErrorMessage, nil, http.StatusInternalServerError, false)
		return
	}

}

func mapBytes(w http.ResponseWriter, r *http.Request, object interface{}, bytes []byte) {

	err1 := json.Unmarshal(bytes, &object)

	if err1 != nil {
		respondWith(w, r, err1, HTTPInternalServerErrorMessage, nil, http.StatusInternalServerError, false)
		return
	}

}

func mapToBytes(w http.ResponseWriter, r *http.Request, document interface{}) []byte {

	j, err0 := bson.MarshalExtJSON(document, false, false)

	if err0 != nil {
		respondWith(w, r, err0, HTTPInternalServerErrorMessage, nil, http.StatusInternalServerError, false)
		return nil
	}

	return j

}

func getAccessToken(r *http.Request) string {

	authToken := r.Header.Get("Authorization")
	splitToken := strings.Split(authToken, "Bearer ")
	authToken = splitToken[1]

	return REDISCLIENT.Get(authToken).Val()

}

func logDebugMessage(message string) {

	rlog.Debug(message)

}

func logErrorMessage(message string) {

	rlog.Error(message)

}

func logInfoMessage(message string) {

	rlog.Info(message)

}

func containsInArray(array []string, value string) bool {

	for _, k := range array {
		if k == value {
			return true
		}
	}

	return false
}

func containsInCategoryNodes(array []CATEGORYTREENODE, value string) bool {

	for _, k := range array {
		if k.Name == value {
			return true
		}
	}

	return false
}

func removeElementsFromArray(array []string, remove []string) []string {

loop:
	for i := 0; i < len(array); i++ {
		element := array[i]
		for _, rem := range remove {
			if element == rem {
				array = append(array[:i], array[i+1:]...)
				i--
				continue loop
			}
		}
	}

	return array

}

func cleanCategoryPath(path string) string {

	path = strings.TrimSpace(path)
	path = strings.Join(strings.Fields(strings.TrimSpace(path)), " ")
	path = strings.Trim(path, ">")

	if strings.Contains(path, ">") {
		splits := strings.Split(path, ">")
		var cleanerCategoryPath string
		for _, val := range splits {

			if val == "" {
				continue
			}

			val = strings.TrimSpace(val)
			cleanerCategoryPath += val + ">"
		}
		path = cleanerCategoryPath[:len(cleanerCategoryPath)-1]
	}

	return strings.TrimSpace(path)

}

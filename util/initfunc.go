package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandString return a random string in give length
func RandomString(length int) string {
	buf := make([]byte, length/2+1)
	rand.Read(buf)
	ret := hex.EncodeToString(buf)
	return ret[:length]
}

// StringArrayContains test if string is contains in a array
func StringArrayContains(array []string, name string) bool {
	for _, i := range array {
		if i == name {
			return true
		}
	}
	return false
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	payload, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	return err
}

func Sha256HMAC(payload, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(payload))
	return h.Sum(nil)
}

type initCallback struct {
	Callback func()
	Name     string
}

var initFuncs = map[int][]initCallback{}

func RegisterInitFunction(name string, init func(), priority int) {
	_, ok := initFuncs[priority]
	i := initCallback{
		Name:     name,
		Callback: init,
	}
	if ok {
		initFuncs[priority] = append(initFuncs[priority], i)
	} else {
		initFuncs[priority] = []initCallback{i}
	}
}

func DoInit() {
	var priority = make([]int, 0, len(initFuncs))
	for i := range initFuncs {
		priority = append(priority, i)
	}
	sort.Ints(priority)
	for _, i := range priority {
		for _, f := range initFuncs[i] {
			log.Printf("run init hook %d: %s", i, f.Name)
			f.Callback()
		}
	}
}

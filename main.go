package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// json.UnmarshalとMarshalとエンコード

type T struct{}

type Person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age,omitempty"`
	Nicknames []string `json:"nicknames,omitempty"`
	T         *T       `json:"T,omitempty"`
}

// func (p *Person) UnmarshalJSON(b []byte) error {
// 	type Person2 struct {
// 		Name string
// 	}
// 	var p2 Person2
// 	err := json.Unmarshal(b, &p2)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	p.Name = p2.Name + "!"
// 	return err
// }

// func (p Person) MarshalJSON() ([]byte, error) {
// 	v, err := json.Marshal(&struct{
// 		Name string
// 	}{
// 		Name: "Mr." + p.Name,
// 	})
// 	return v, err
// }

// hmacでAPI認証

var DB = map[string]string{
	"User1Key": "User1Secret",
	"User2Key": "User2Secret",
}

func Server(apiKey, sign string, data []byte) {
	apiSecret := DB[apiKey]
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)
	expectedHMAC := hex.EncodeToString(h.Sum(nil))
	fmt.Println(sign == expectedHMAC)
}

func main() {
	// http

	// resp, err := http.Get("http://example.com")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	base, _ := url.Parse("https://example.com")
	reference, _ := url.Parse("/test?a=1&b=2")
	endpoint := base.ResolveReference(reference).String()
	fmt.Println(endpoint)

	req, _ := http.NewRequest("GET", endpoint, nil)
	// req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte("password")))
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	q := req.URL.Query()
	q.Add("c", "3&%")
	fmt.Println(q)
	fmt.Println(q.Encode())
	req.URL.RawQuery = q.Encode()

	var client *http.Client = &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	// json.UnmarshalとMarshalとエンコード
	b := []byte(`{"name":"Mike", "age":20, "nicknames":[]}`)
	var p Person
	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Name, p.Age, p.Nicknames)

	v, _ := json.Marshal(p)
	fmt.Println(v)
	fmt.Println(string(v))

	// hmacでAPI認証

	const apiKey = "User2Key"
	const apiSecret = "User2Secret"

	data := []byte("data")
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)
	sign := hex.EncodeToString(h.Sum(nil))

	fmt.Println(sign)

	Server(apiKey, sign, data)
}

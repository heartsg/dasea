// Package gorequest inspired by Nodejs Requests provides easy-way to write http client
// Package requests based on gorequest, and also add in some functions provided by goreq
// that is not provided by gorequest

// Note that goreq's drawback is that, when CookieJar is supplied, it will always create
// a new client (however, it is recommended to reuse the client when possible).
// Goreq also does not support a way to clear request state.

// Therefore, our implementation is mainly based on gorequest.

// What goreq provided but not in goreqest (so we will add in here):
//  - default redirect policy
//  - goreq supports compression and request timeout, but 
//    no use currently, so we don't implement it

// some functions we add for goreqeust
//  - optimization to not insert body back to resp after it's read
//  - auto close resp body if in setting

package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
	"log"

	"golang.org/x/net/publicsuffix"
)


// HTTP methods we support
const (
	POST   = "POST"
	GET    = "GET"
	HEAD   = "HEAD"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

// A Requests is a object storing all request data for client.
type Requests struct {
	// gorequest original except logger changed
	Url               string
	Method            string
	Header            map[string]string
	TargetType        string
	ForceType         string
	Data              map[string]interface{}
	SliceData         []interface{}
	FormData          url.Values
	QueryData         url.Values
	BounceToRawString bool
	RawString         string
	Client            *http.Client
	Transport         *http.Transport
	Cookies           []*http.Cookie
	Errors            []error
	BasicAuth         struct{ Username, Password string }
	Debug             bool
	logger            *log.Logger
	
	// Keep resp.Body still readable after resp is already read out
	responseBodyValid  bool
	
}

// Used to create a new Requests object.
func New() *Requests {
	// We create a default cookiejar. If a different cookiejar is needed,
	// uses shall provide their own client.
	cookiejarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookiejarOptions)
	r := &Requests{
		TargetType:        "json",
		Data:              make(map[string]interface{}),
		Header:            make(map[string]string),
		RawString:         "",
		SliceData:         []interface{}{},
		FormData:          url.Values{},
		QueryData:         url.Values{},
		BounceToRawString: false,
		Client:            &http.Client{Jar: jar},
		Transport:         &http.Transport{},
		Cookies:           make([]*http.Cookie, 0),
		Errors:            nil,
		BasicAuth:         struct{ Username, Password string }{},
		Debug:             false,
		logger:            log.New(os.Stderr, "[gorequest]", log.LstdFlags),
		responseBodyValid: false,
	}

	return r
}

func (s *Requests) SetResponseBodyValid(valid bool) *Requests {
	s.responseBodyValid = valid
	return s
}

// Enable the debug mode which logs request/response detail
func (s *Requests) SetDebug(enable bool) *Requests {
	s.Debug = enable
	return s
}

func (s *Requests) SetLogger(logger *log.Logger) *Requests {
	s.logger = logger
	return s
}

// Clear Requests data for another new request.
func (s *Requests) ClearRequests() {
	s.Url = ""
	s.Method = ""
	s.Header = make(map[string]string)
	s.Data = make(map[string]interface{})
	s.SliceData = []interface{}{}
	s.FormData = url.Values{}
	s.QueryData = url.Values{}
	s.BounceToRawString = false
	s.RawString = ""
	s.ForceType = ""
	s.TargetType = "json"
	s.Cookies = make([]*http.Cookie, 0)
	s.Errors = nil
	s.responseBodyValid = false
}

func (s *Requests) Request(targetUrl string, method string) *Requests {
	switch method {
	case GET:
		return s.Get(targetUrl)
	case POST:
		return s.Post(targetUrl)
	case HEAD:
		return s.Head(targetUrl)
	case PUT:
		return s.Put(targetUrl)
	case DELETE:
		return s.Delete(targetUrl)
	case PATCH:
		return s.Patch(targetUrl)
	default:
		return s.Get(targetUrl)
	}
}

func (s *Requests) Get(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = GET
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *Requests) Post(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = POST
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *Requests) Head(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = HEAD
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *Requests) Put(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = PUT
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *Requests) Delete(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = DELETE
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *Requests) Patch(targetUrl string) *Requests {
	s.ClearRequests()
	s.Method = PATCH
	s.Url = targetUrl
	s.Errors = nil
	return s
}

// Set is used for setting header fields.
// Example. To set `Accept` as `application/json`
//
//    gorequest.New().
//      Post("/gamelist").
//      Set("Accept", "application/json").
//      End()
func (s *Requests) Set(param string, value string) *Requests {
	s.Header[param] = value
	return s
}

// SetBasicAuth sets the basic authentication header
// Example. To set the header for username "myuser" and password "mypass"
//
//    gorequest.New()
//      Post("/gamelist").
//      SetBasicAuth("myuser", "mypass").
//      End()
func (s *Requests) SetBasicAuth(username string, password string) *Requests {
	s.BasicAuth = struct{ Username, Password string }{username, password}
	return s
}

// AddCookie adds a cookie to the request. The behavior is the same as AddCookie on Request from net/http
func (s *Requests) AddCookie(c *http.Cookie) *Requests {
	s.Cookies = append(s.Cookies, c)
	return s
}

// AddCookies is a convenient method to add multiple cookies
func (s *Requests) AddCookies(cookies []*http.Cookie) *Requests {
	s.Cookies = append(s.Cookies, cookies...)
	return s
}

var Types = map[string]string{
	"html":       "text/html",
	"json":       "application/json",
	"xml":        "application/xml",
	"text":       "text/plain",
	"urlencoded": "application/x-www-form-urlencoded",
	"form":       "application/x-www-form-urlencoded",
	"form-data":  "application/x-www-form-urlencoded",
}

// Type is a convenience function to specify the data type to send.
// For example, to send data as `application/x-www-form-urlencoded` :
//
//    gorequest.New().
//      Post("/recipe").
//      Type("form").
//      Send(`{ name: "egg benedict", category: "brunch" }`).
//      End()
//
// This will POST the body "name=egg benedict&category=brunch" to url /recipe
//
// GoRequest supports
//
//    "text/html" uses "html"
//    "application/json" uses "json"
//    "application/xml" uses "xml"
//    "text/plain" uses "text"
//    "application/x-www-form-urlencoded" uses "urlencoded", "form" or "form-data"
//
func (s *Requests) Type(typeStr string) *Requests {
	if _, ok := Types[typeStr]; ok {
		s.ForceType = typeStr
	} else {
		s.Errors = append(s.Errors, errors.New("Type func: incorrect type \""+typeStr+"\""))
	}
	return s
}

// Query function accepts either json string or strings which will form a query-string in url of GET method or body of POST method.
// For example, making "/search?query=bicycle&size=50x50&weight=20kg" using GET method:
//
//      gorequest.New().
//        Get("/search").
//        Query(`{ query: 'bicycle' }`).
//        Query(`{ size: '50x50' }`).
//        Query(`{ weight: '20kg' }`).
//        End()
//
// Or you can put multiple json values:
//
//      gorequest.New().
//        Get("/search").
//        Query(`{ query: 'bicycle', size: '50x50', weight: '20kg' }`).
//        End()
//
// Strings are also acceptable:
//
//      gorequest.New().
//        Get("/search").
//        Query("query=bicycle&size=50x50").
//        Query("weight=20kg").
//        End()
//
// Or even Mixed! :)
//
//      gorequest.New().
//        Get("/search").
//        Query("query=bicycle").
//        Query(`{ size: '50x50', weight:'20kg' }`).
//        End()
//
func (s *Requests) Query(content interface{}) *Requests {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.queryString(v.String())
	case reflect.Struct:
		s.queryStruct(v.Interface())
	default:
	}
	return s
}

func (s *Requests) queryStruct(content interface{}) *Requests {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				s.QueryData.Add(k, v.(string))
			}
		}
	}
	return s
}

func (s *Requests) queryString(content string) *Requests {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		if queryVal, err := url.ParseQuery(content); err == nil {
			for k, _ := range queryVal {
				s.QueryData.Add(k, queryVal.Get(k))
			}
		} else {
			s.Errors = append(s.Errors, err)
		}
		// TODO: need to check correct format of 'field=val&field=val&...'
	}
	return s
}

// As Go conventions accepts ; as a synonym for &. (https://github.com/golang/go/issues/2210)
// Thus, Query won't accept ; in a querystring if we provide something like fields=f1;f2;f3
// This Param is then created as an alternative method to solve this.
func (s *Requests) Param(key string, value string) *Requests {
	s.QueryData.Add(key, value)
	return s
}

func (s *Requests) Timeout(timeout time.Duration) *Requests {
	s.Transport.Dial = func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			s.Errors = append(s.Errors, err)
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(timeout))
		return conn, nil
	}
	return s
}

// Set TLSClientConfig for underling Transport.
// One example is you can use it to disable security check (https):
//
//      gorequest.New().TLSClientConfig(&tls.Config{ InsecureSkipVerify: true}).
//        Get("https://disable-security-check.com").
//        End()
//
func (s *Requests) TLSClientConfig(config *tls.Config) *Requests {
	s.Transport.TLSClientConfig = config
	return s
}

// Proxy function accepts a proxy url string to setup proxy url for any request.
// It provides a convenience way to setup proxy which have advantages over usual old ways.
// One example is you might try to set `http_proxy` environment. This means you are setting proxy up for all the requests.
// You will not be able to send different request with different proxy unless you change your `http_proxy` environment again.
// Another example is using Golang proxy setting. This is normal prefer way to do but too verbase compared to GoRequest's Proxy:
//
//      gorequest.New().Proxy("http://myproxy:9999").
//        Post("http://www.google.com").
//        End()
//
// To set no_proxy, just put empty string to Proxy func:
//
//      gorequest.New().Proxy("").
//        Post("http://www.google.com").
//        End()
//
func (s *Requests) Proxy(proxyUrl string) *Requests {
	parsedProxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		s.Errors = append(s.Errors, err)
	} else if proxyUrl == "" {
		s.Transport.Proxy = nil
	} else {
		s.Transport.Proxy = http.ProxyURL(parsedProxyUrl)
	}
	return s
}

func (s *Requests) RedirectPolicy(policy func(req *http.Request, via []*http.Request) error) *Requests {
	s.Client.CheckRedirect = policy
	return s
}

func (s *Requests) Redirect(maxRedirects int, redirectHeaders bool) *Requests {
	return s.RedirectPolicy(func(req *http.Request, via []*http.Request) error {
		if len(via) > maxRedirects {
			return errors.New("Error redirecting. MaxRedirects reached")
		}
	
		//By default Golang will not redirect request headers
		// https://code.google.com/p/go/issues/detail?id=4800&q=request%20header
		if redirectHeaders {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}
		}
		return nil		
	})
}

// Send function try to unmarshal everything into map first, which is not efficient.
// Provide raw string functionality (have to call each request)
func (s *Requests) SendRawString(content string)  *Requests {
	s.BounceToRawString = true
	s.RawString = content
	return s
}
func (s *Requests) SendRawStruct(content interface{}) *Requests {
	s.BounceToRawString = true
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		s.RawString = string(marshalContent)
	}
	return s
}

// Send function accepts either json string or query strings which is usually used to assign data to POST or PUT method.
// Without specifying any type, if you give Send with json data, you are doing requesting in json format:
//
//      gorequest.New().
//        Post("/search").
//        Send(`{ query: 'sushi' }`).
//        End()
//
// While if you use at least one of querystring, GoRequest understands and automatically set the Content-Type to `application/x-www-form-urlencoded`
//
//      gorequest.New().
//        Post("/search").
//        Send("query=tonkatsu").
//        End()
//
// So, if you want to strictly send json format, you need to use Type func to set it as `json` (Please see more details in Type function).
// You can also do multiple chain of Send:
//
//      gorequest.New().
//        Post("/search").
//        Send("query=bicycle&size=50x50").
//        Send(`{ wheel: '4'}`).
//        End()
//
// From v0.2.0, Send function provide another convenience way to work with Struct type. You can mix and match it with json and query string:
//
//      type BrowserVersionSupport struct {
//        Chrome string
//        Firefox string
//      }
//      ver := BrowserVersionSupport{ Chrome: "37.0.2041.6", Firefox: "30.0" }
//      gorequest.New().
//        Post("/update_version").
//        Send(ver).
//        Send(`{"Safari":"5.1.10"}`).
//        End()
//
// If you have set Type to text or Content-Type to text/plain, content will be sent as raw string in body instead of form
//
//      gorequest.New().
//        Post("/greet").
//        Type("text").
//        Send("hello world").
//        End()
//
func (s *Requests) Send(content interface{}) *Requests {
	// TODO: add normal text mode or other mode to Send func
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.SendString(v.String())
	case reflect.Struct:
		s.SendStruct(v.Interface())
	case reflect.Slice:
		// change to slice
		slice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			slice[i] = v.Index(i).Interface()
		}
		s.SendSlice(slice)
	case reflect.Ptr:
		switch v.Type().Elem().Kind() {
		case reflect.Struct:
			s.SendStruct(v.Interface())
		default:
		}
	default:
		// TODO: leave default for handling other types in the future such as number, byte, etc...
		// TODO: Add support for slice and array
	}
	return s
}

// SendSlice (similar to SendString) returns Requests's itself for any next chain and takes content []interface{} as a parameter.
// Its duty is to append slice of interface{} into s.SliceData ([]interface{}) which later changes into json array in the End() func.
func (s *Requests) SendSlice(content []interface{}) *Requests {
	s.SliceData = append(s.SliceData, content...)
	return s
}

// SendStruct (similar to SendString) returns Requests's itself for any next chain and takes content interface{} as a parameter.
// Its duty is to transfrom interface{} (implicitly always a struct) into s.Data (map[string]interface{}) which later changes into appropriate format such as json, form, text, etc. in the End() func.
func (s *Requests) SendStruct(content interface{}) *Requests {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				s.Data[k] = v
			}
		}
	}
	return s
}

// SendString returns Requests's itself for any next chain and takes content string as a parameter.
// Its duty is to transform String into s.Data (map[string]interface{}) which later changes into appropriate format such as json, form, text, etc. in the End func.
// Send implicitly uses SendString and you should use Send instead of this.
func (s *Requests) SendString(content string) *Requests {
	if !s.BounceToRawString {
		var val interface{}
		d := json.NewDecoder(strings.NewReader(content))
		d.UseNumber()
		if err := d.Decode(&val); err == nil {
			switch v := reflect.ValueOf(val); v.Kind() {
			case reflect.Map:
				for k, v := range val.(map[string]interface{}) {
					s.Data[k] = v
				}
			// add to SliceData
			case reflect.Slice:
				s.SendSlice(val.([]interface{}))
			// bounce to rawstring if it is arrayjson, or others
			default:
				s.BounceToRawString = true
			}
		} else if formVal, err := url.ParseQuery(content); err == nil {
			for k, _ := range formVal {
				// make it array if already have key
				if val, ok := s.Data[k]; ok {
					var strArray []string
					strArray = append(strArray, formVal.Get(k))
					// check if previous data is one string or array
					switch oldValue := val.(type) {
					case []string:
						strArray = append(strArray, oldValue...)
					case string:
						strArray = append(strArray, oldValue)
					}
					s.Data[k] = strArray
				} else {
					// make it just string if does not already have same key
					s.Data[k] = formVal.Get(k)
				}
			}
			s.TargetType = "form"
		} else {
			s.BounceToRawString = true
		}
	}
	// Dump all contents to RawString in case in the end user doesn't want json or form.
	s.RawString += content
	return s
}

func changeMapToURLValues(data map[string]interface{}) url.Values {
	var newUrlValues = url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			newUrlValues.Add(k, val)
		case []string:
			for _, element := range val {
				newUrlValues.Add(k, element)
			}
		// if a number, change to string
		// json.Number used to protect against a wrong (for GoRequest) default conversion
		// which always converts number to float64.
		// This type is caused by using Decoder.UseNumber()
		case json.Number:
			newUrlValues.Add(k, string(val))
		}
	}
	return newUrlValues
}

// End is the most important function that you need to call when ending the chain. The request won't proceed without calling it.
// End function returns Response which matchs the structure of Response type in Golang's http package (but without Body data). The body data itself returns as a string in a 2nd return value.
// Lastly but worth noticing, error array (NOTE: not just single error value) is returned as a 3rd value and nil otherwise.
//
// For example:
//
//    resp, body, errs := gorequest.New().Get("http://www.google.com").End()
//    if (errs != nil) {
//      fmt.Println(errs)
//    }
//    fmt.Println(resp, body)
//
// Moreover, End function also supports callback which you can put as a parameter.
// This extends the flexibility and makes GoRequest fun and clean! You can use GoRequest in whatever style you love!
//
// For example:
//
//    func printBody(resp gorequest.Response, body string, errs []error){
//      fmt.Println(resp.Status)
//    }
//    gorequest.New().Get("http://www..google.com").End(printBody)
//
func (s *Requests) End(callback ...func(response *http.Response, body string, errs []error)) (*http.Response, string, []error) {
	var bytesCallback []func(response *http.Response, body []byte, errs []error)
	if len(callback) > 0 {
		bytesCallback = []func(response *http.Response, body []byte, errs []error){
			func(response *http.Response, body []byte, errs []error) {
				callback[0](response, string(body), errs)
			},
		}
	}
	resp, body, errs := s.EndBytes(bytesCallback...)
	bodyString := string(body)
	return resp, bodyString, errs
}

// EndBytes should be used when you want the body as bytes. The callbacks work the same way as with `End`, except that a byte array is used instead of a string.
func (s *Requests) EndBytes(callback ...func(response *http.Response, body []byte, errs []error)) (*http.Response, []byte, []error) {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	// check whether there is an error. if yes, return all errors
	if len(s.Errors) != 0 {
		return nil, nil, s.Errors
	}
	// check if there is forced type
	switch s.ForceType {
	case "json", "form", "xml", "text":
		s.TargetType = s.ForceType
		// If forcetype is not set, check whether user set Content-Type header.
		// If yes, also bounce to the correct supported TargetType automatically.
	default:
		for k, v := range Types {
			if s.Header["Content-Type"] == v {
				s.TargetType = k
			}
		}
	}

	// if slice and map get mixed, let's bounce to rawstring
	if len(s.Data) != 0 && len(s.SliceData) != 0 {
		s.BounceToRawString = true
	}

	switch s.Method {
	case POST, PUT, PATCH:
		if s.TargetType == "json" {
			// If-case to give support to json array. we check if
			// 1) Map only: send it as json map from s.Data
			// 2) Array or Mix of map & array or others: send it as rawstring from s.RawString
			var contentJson []byte
			if s.BounceToRawString {
				contentJson = []byte(s.RawString)
			} else if len(s.Data) != 0 {
				contentJson, _ = json.Marshal(s.Data)
			} else if len(s.SliceData) != 0 {
				contentJson, _ = json.Marshal(s.SliceData)
			}
			contentReader := bytes.NewReader(contentJson)
			req, err = http.NewRequest(s.Method, s.Url, contentReader)
			if err != nil {
				s.Errors = append(s.Errors, err)
				return nil, nil, s.Errors
			}
			req.Header.Set("Content-Type", "application/json")
		} else if s.TargetType == "form" {
			var contentForm []byte
			if s.BounceToRawString || len(s.SliceData) != 0 {
				contentForm = []byte(s.RawString)
			} else {
				formData := changeMapToURLValues(s.Data)
				contentForm = []byte(formData.Encode())
			}
			contentReader := bytes.NewReader(contentForm)
			req, err = http.NewRequest(s.Method, s.Url, contentReader)
			if err != nil {
				s.Errors = append(s.Errors, err)
				return nil, nil, s.Errors
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else if s.TargetType == "text" {
			req, err = http.NewRequest(s.Method, s.Url, strings.NewReader(s.RawString))
			req.Header.Set("Content-Type", "text/plain")
		} else if s.TargetType == "xml" {
			req, err = http.NewRequest(s.Method, s.Url, strings.NewReader(s.RawString))
			req.Header.Set("Content-Type", "application/xml")
		} else {
			// TODO: if nothing match, let's return warning here
		}
	case GET, HEAD, DELETE:
		req, err = http.NewRequest(s.Method, s.Url, nil)
		if err != nil {
			s.Errors = append(s.Errors, err)
			return nil, nil, s.Errors
		}
	}
	for k, v := range s.Header {
		req.Header.Set(k, v)
	}
	// Add all querystring from Query func
	q := req.URL.Query()
	for k, v := range s.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()

	// Add basic auth
	if s.BasicAuth != struct{ Username, Password string }{} {
		req.SetBasicAuth(s.BasicAuth.Username, s.BasicAuth.Password)
	}

	// Add cookies
	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}

	// Set Transport
	s.Client.Transport = s.Transport

	// Log details of this request
	if s.Debug {
		dump, err := httputil.DumpRequest(req, true)
		s.logger.SetPrefix("[http] ")
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Request: %s", string(dump))
		}
	}
	// Send request
	resp, err = s.Client.Do(req)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, s.Errors
	}
	defer resp.Body.Close()

	// Log details of this response
	if s.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	// Reset resp.Body so it can be use again
	if s.responseBodyValid {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	// deep copy response to give it to both return and callback func
	respCallback := *resp
	if len(callback) != 0 {
		callback[0](&respCallback, body, s.Errors)
	}
	return resp, body, nil
}

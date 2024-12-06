package main

import (
	"fmt"
	"net/url"
	"strings"
)

// type of http/rest request method
type RestMethod string

func (m RestMethod) String() string {
	return string(m)
}

const (
	GET    RestMethod = "GET"
	POST   RestMethod = "POST"
	PUT    RestMethod = "PUT"
	PATCH  RestMethod = "PATCH"
	DELETE RestMethod = "DELETE"
)

func ParseRestMethod(s string) (RestMethod, error) {
	switch s {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "PATCH":
		return PATCH, nil
	case "DELETE":
		return DELETE, nil
	case "get":
		return GET, nil
	case "post":
		return POST, nil
	case "put":
		return PUT, nil
	case "patch":
		return PATCH, nil
	case "delete":
		return DELETE, nil
	default:
		return "", fmt.Errorf("Unknown method: %s. Only GET, POST, PUT, PATCH and DELETE are supported", s)
	}
}

// type of argument
type ConfigArgType string

const (
	STRING ConfigArgType = "string"
	INT    ConfigArgType = "int"
	FLOAT  ConfigArgType = "float"
	BOOL   ConfigArgType = "bool"
)

func ParseConfigArgType(s string) (ConfigArgType, error) {
	switch s {
	case "string":
		return STRING, nil
	case "int":
		return INT, nil
	case "float":
		return FLOAT, nil
	case "bool":
		return BOOL, nil
	default:
		return "", fmt.Errorf("Unknown argument type: %s. Only string, int, float and bool are supported", s)
	}
}

// type of body type
type BodyType string

const (
	TEXT BodyType = "text"
	JSON BodyType = "json"
)

func ParseBodyType(s string) (BodyType, error) {
	switch s {
	case "text":
		return TEXT, nil
	case "json":
		return JSON, nil
	default:
		return "", fmt.Errorf("Unknown body type: %s. Only text and json are supported", s)
	}
}

// Selector to show result
type Selector interface {
	Run(headers map[string]string, body string, bodyType BodyType) string
	Kind() string
}

func CreateSelector(showString string) (Selector, error) {
	parts := strings.Split(showString, ":")
	switch parts[0] {
	case "header":
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid header selector: %s. Format should be header.<key>", showString)
		}
		return HeaderSelector{Key: parts[1]}, nil
	case "cookie":
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid cookie selector: %s. Format should be cookie.<key>", showString)
		}
		return CookieSelector{Key: parts[1]}, nil
	case "body":
		if len(parts) == 1 {
			return BodySelector{}, nil
		}
		depth := strings.Split(parts[1], ".")
		return BodySelector{Depth: depth}, nil
	}
	return nil, fmt.Errorf("Unknown selector: %s", showString)
}

// to select header
type HeaderSelector struct {
	Key string
}

func (hs HeaderSelector) Run(headers map[string]string, body string, bodyType BodyType) string {
	return headers[hs.Key]
}

func (hs HeaderSelector) Kind() string {
	return "header"
}

// to select body
type BodySelector struct {
	Depth []string
}

func (bs BodySelector) Run(headers map[string]string, body string, bodyType BodyType) string {
	return body
}

func (bs BodySelector) Kind() string {
	return "body"
}

// cookie selector
type CookieSelector struct {
	Key string
}

func (cs CookieSelector) Run(headers map[string]string, body string, bodyType BodyType) string {
	return headers[cs.Key]
}

func (cs CookieSelector) Kind() string {
	return "cookie"
}

// Main config
type Config struct {
	Settings Optional[GlobalSettings]
	Requests map[string]RequestConfig
}

type GlobalSettings struct {
	BaseUrl url.URL
	Output  string
	EnvVars map[string]string
	Include []string
}

type RequestConfig struct {
	Args     []Argument
	Method   RestMethod
	Url      url.URL
	Select   []Selector
	Headers  []Header
	Body     Optional[[]byte]
	BodyType BodyType
}

type Argument struct {
	Key     string
	ArgType ConfigArgType
}

type Header struct {
	Key   string
	Value string
}

func NewConfigFromRaw(raw RawConfig) (Config, error) {
	// parses global settings
	opt_gs := NewOptionalEmpty[GlobalSettings]()
	if raw_gs := raw.Settings; raw_gs != nil {
		gs := GlobalSettings{}
		if base_url := raw_gs.BaseUrl; base_url != nil {
			bu, err := url.Parse(*base_url)
			if err != nil {
				return Config{}, err
			}
			gs.BaseUrl = *bu
		}
		if output := raw_gs.Output; output != nil {
			gs.Output = *output
		}
		if load_env := raw_gs.LoadEnv; load_env != nil {
			gs.EnvVars = make(map[string]string)
		}
		if include := raw_gs.Include; include != nil {
			gs.Include = *include
		}

		opt_gs.SetValue(gs)
	}
	// parses config requests
	requests := make(map[string]RequestConfig)
	for name, raw_req := range raw.Requests {
		raw_req, err := FillExtend(raw_req, raw.Requests)
		if err != nil {
			return Config{}, err
		}
		rc := RequestConfig{}

		// parses arguments with their types
		for _, arg := range raw_req.Args {
			parts := strings.Split(arg, ":")
			// if no type is specified, default to string
			if len(parts) == 1 {
				rc.Args = append(rc.Args, Argument{Key: parts[0], ArgType: "string"})
			} else if len(parts) == 2 {
				// else if the type is specified
				arg_type, err := ParseConfigArgType(parts[1])
				if err != nil {
					return Config{}, err
				}
				rc.Args = append(rc.Args, Argument{Key: parts[0], ArgType: arg_type})
			} else {
				return Config{}, fmt.Errorf("Invalid argument: %s. Format should be key:type or key", arg)
			}
		}

		// parses rest method
		if method := raw_req.Method; method != nil {
			m, err := ParseRestMethod(*method)
			if err != nil {
				return Config{}, err
			}
			rc.Method = m
		} else {
			rc.Method = GET
		}

		// parses url and path
		if path := raw_req.Path; path != nil {

			// if global base url is not defined, returns error
			if gs, ok := opt_gs.GetValue(); ok {
				rc.Url = gs.BaseUrl
				rc.Url.Path = *path
			} else {
				return Config{}, fmt.Errorf("Url is not defined in global settings, so path can't be used")
			}

			// if url is defined, path is ignored
			if raw_url := raw_req.Url; raw_url != nil {
				u, err := url.Parse(*raw_url)
				if err != nil {
					return Config{}, err
				}
				rc.Url = *u
			}
		}

		// parses headers
		for _, header := range raw_req.Headers {
			parts := strings.Split(header, ":")
			if len(parts) != 2 {
				return Config{}, fmt.Errorf("Header must be of format key:value")
			}
			rc.Headers = append(rc.Headers, Header{Key: parts[0], Value: parts[1]})
		}

		// parses body type
		if body_type := raw_req.BodyType; body_type != nil {
			bt, err := ParseBodyType(*body_type)
			if err != nil {
				return Config{}, err
			}
			rc.BodyType = bt
		} else {
			rc.BodyType = JSON
		}

		// parses body
		if body := raw_req.Body; body != nil {
			rc.Body = NewOptional([]byte(*body))
		} else {
			rc.Body = NewOptionalEmpty[[]byte]()
		}

		// parses selector
		for _, raw_sel := range raw_req.Select {
			selector, err := CreateSelector(raw_sel)
			if err != nil {
				return Config{}, err
			}
			rc.Select = append(rc.Select, selector)
		}

		requests[name] = rc
	}
	return Config{Settings: opt_gs, Requests: requests}, nil
}

// Fills extended request
func FillExtend(raw_req RawRequestConfig, all_reqs map[string]RawRequestConfig) (RawRequestConfig, error) {
	if extend := raw_req.Extend; extend != nil {
		target, found := all_reqs[*extend]
		if !found {
			return RawRequestConfig{}, fmt.Errorf("Can't find extended request: %s", *extend)
		}
		if raw_req.Url == nil {
			raw_req.Url = target.Url
		}
		if raw_req.Path == nil {
			raw_req.Path = target.Path
		}
		if raw_req.Method == nil {
			raw_req.Method = target.Method
		}
		if raw_req.Body == nil {
			raw_req.Body = target.Body
		}
		raw_req.Select = append(raw_req.Select, target.Select...)
		raw_req.Headers = append(raw_req.Headers, target.Headers...)
		raw_req.Args = append(raw_req.Args, target.Args...)
	}
	return raw_req, nil
}

func (c Config) String() string {
	baseString := ""
	if gs, ok := c.Settings.GetValue(); ok {
		baseString += "Settings:\n"
		baseString += fmt.Sprintf("  BaseUrl: %s\n", gs.BaseUrl.String())
		baseString += fmt.Sprintf("  Output: %s\n", gs.Output)
		baseString += fmt.Sprintf("  EnvVars: %s\n", gs.EnvVars)
		baseString += fmt.Sprintf("  Include: %s\n", gs.Include)
		baseString += "\n"
	}
	for name, req := range c.Requests {
		baseString += fmt.Sprintf("Request: %s\n", name)
		baseString += fmt.Sprintf("  Args: %s\n", req.Args)
		baseString += fmt.Sprintf("  Method: %s\n", req.Method)
		baseString += fmt.Sprintf("  Url: %s\n", req.Url.String())
		baseString += fmt.Sprintf("  Select: %s\n", req.Select)
		baseString += fmt.Sprintf("  Headers: %s\n", req.Headers)
		if body, ok := req.Body.GetValue(); ok {
			baseString += fmt.Sprintf("  Body: %s\n", body)
		}
		baseString += "\n"
	}
	return baseString
}

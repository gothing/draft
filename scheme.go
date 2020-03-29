package draft

import (
	"sync"

	"github.com/gothing/draft/reflect"
)

// Scheme — описательная часть api и сбосов его использования
type Scheme struct {
	mu         sync.Mutex
	url        string
	name       string
	descr      string
	cases      []*SchemeCase
	defAccess  AccessType
	defMethod  string
	defParams  interface{}
	defBody    interface{}
	activeCase *SchemeCase
}

// SchemeCase — описание и пример использование
type SchemeCase struct {
	name   string
	descr  string
	access AccessType
	status StatusType
	method string
	params interface{}
	body   interface{}
}

// Name —
func (s *Scheme) Name(v string) {
	s.name = v
}

// Access — выставить права доступа к апишке или `case`
func (s *Scheme) Access(v AccessType) {
	if s.activeCase != nil {
		s.activeCase.access = v
	} else {
		s.defAccess = v
	}
}

// Method — выставить метод к апишке или `case`
func (s *Scheme) Method(v string) {
	if s.activeCase != nil {
		s.activeCase.method = v
	} else {
		s.defMethod = v
	}
}

// Description — выставить описание к апишке или `case`
func (s *Scheme) Description(v string) {
	if s.activeCase != nil {
		s.activeCase.descr = v
	} else {
		s.descr = v
	}
}

// Params — выставить параметры к апишке или `case`
func (s *Scheme) Params(v interface{}) {
	if s.activeCase != nil {
		s.activeCase.params = v
	} else {
		s.defParams = v
	}
}

// Body — выставить ответ к апишке или `case`
func (s *Scheme) Body(v interface{}) {
	if s.activeCase != nil {
		s.activeCase.body = v
	} else {
		s.defBody = v
	}
}

// Case — определить описание и пример использования
func (s *Scheme) Case(status StatusType, name string, fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeCase = &SchemeCase{
		status: status,
		name:   name,
		access: s.defAccess,
		params: s.defParams,
		body:   s.defBody,
	}
	s.cases = append(s.cases, s.activeCase)
	fn()
	s.activeCase = nil
}

// JSONScheme —
type JSONScheme struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Detail      map[StatusType]*JSONSchemeDetail `json:"detail"`
	Cases       []*SchemeCase
}

// JSONSchemeDetail —
type JSONSchemeDetail struct {
	Access   AccessType          `json:"access"`
	Request  *JSONSchemeRequest  `json:"request"`
	Response *JSONSchemeResponse `json:"response"`
}

// JSONSchemeRequest -
type JSONSchemeRequest struct {
	Params map[string]reflect.Item `json:"params"`
}

// JSONSchemeResponse -
type JSONSchemeResponse struct {
	Body map[string]reflect.Item `json:"body"`
}

// GetCaseByStatus — определить описание и пример использования
func (s *Scheme) GetCaseByStatus(v StatusType) *SchemeCase {
	for _, c := range s.cases {
		if c.status == v {
			return c
		}
	}
	return nil
}

// ToJSON — определить описание и пример использования
func (s *Scheme) ToJSON() JSONScheme {
	json := JSONScheme{
		Name:        s.name,
		Description: s.descr,
		Detail:      make(map[StatusType]*JSONSchemeDetail),
		Cases:       s.cases,
	}

	for _, c := range s.cases {
		d, exists := json.Detail[c.status]
		if !exists {
			d = &JSONSchemeDetail{
				Request: &JSONSchemeRequest{
					Params: make(map[string]reflect.Item),
				},

				Response: &JSONSchemeResponse{
					Body: make(map[string]reflect.Item),
				},
			}
			json.Detail[c.status] = d
		}

		d.Access = c.access
		setReflectObjectMap(d.Request.Params, c.params)
		setReflectObjectMap(d.Response.Body, c.body)
	}

	return json
}

func setReflectObjectMap(obj map[string]reflect.Item, v interface{}) {
	if v != nil {
		ref := reflect.Get(v, reflect.Options{
			SnakeCase: true,
		})
		for _, item := range ref.Nested {
			obj[item.Name] = item
		}
	}
}

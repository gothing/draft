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
	Name        string
	Description string
	Access      AccessType
	Status      StatusType
	Method      string
	Params      interface{}
	Body        interface{}
}

// Name —
func (s *Scheme) Name(v string) {
	s.name = v
}

// Access — выставить права доступа к апишке или `case`
func (s *Scheme) Access(v AccessType) {
	if s.activeCase != nil {
		s.activeCase.Access = v
	} else {
		s.defAccess = v
	}
}

// Method — выставить метод к апишке или `case`
func (s *Scheme) Method(v string) {
	if s.activeCase != nil {
		s.activeCase.Method = v
	} else {
		s.defMethod = v
	}
}

// Description — выставить описание к апишке или `case`
func (s *Scheme) Description(v string) {
	if s.activeCase != nil {
		s.activeCase.Description = v
	} else {
		s.descr = v
	}
}

// Params — выставить параметры к апишке или `case`
func (s *Scheme) Params(v interface{}) {
	if s.activeCase != nil {
		s.activeCase.Params = v
	} else {
		s.defParams = v
	}
}

// Body — выставить ответ к апишке или `case`
func (s *Scheme) Body(v interface{}) {
	if s.activeCase != nil {
		s.activeCase.Body = v
	} else {
		s.defBody = v
	}
}

// Case — определить описание и пример использования
func (s *Scheme) Case(status StatusType, name string, fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeCase = &SchemeCase{
		Status: status,
		Name:   name,
		Access: s.defAccess,
		Params: s.defParams,
		Body:   s.defBody,
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
		if c.Status == v {
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
		d, exists := json.Detail[c.Status]
		if !exists {
			d = &JSONSchemeDetail{
				Request: &JSONSchemeRequest{
					Params: make(map[string]reflect.Item),
				},

				Response: &JSONSchemeResponse{
					Body: make(map[string]reflect.Item),
				},
			}
			json.Detail[c.Status] = d
		}

		d.Access = c.Access
		setReflectObjectMap(d.Request.Params, c.Params)
		setReflectObjectMap(d.Response.Body, c.Body)
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

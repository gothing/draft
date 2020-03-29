


```go



func (api) Handle(ctx, params) {
	mock, err := api.InitEndpoint(req)
	if err != nil {
		return
	}

	if mock != nil {
		return mock.Body.(*APIResponse), nil
	}


	return apidoc.Mock(api, req)
}
```

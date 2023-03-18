package fs

import (
	"context"
	"go.keploy.io/server/pkg/models"
	"gopkg.in/yaml.v3"
	"sync"
	"testing"
)

func Test_mockExport_ReadByID(t *testing.T) {
	tests := sync.Map{}
	data := `version: api.keploy.io/v1beta2
kind: Http
name: test-1
spec:
    metadata: {}
    req:
        method: POST
        proto_major: 1
        proto_minor: 1
        url: /api/regression/testcase
        header:
            Accept-Encoding: gzip
            Content-Length: "1667"
            Content-Type: application/json
            User-Agent: Go-http-client/1.1
        body: '{"captured":1674553625,"app_id":"grpc-nested-app","uri":"","http_req":{"method":"","proto_major":0,"proto_minor":0,"url":"","url_params":null,"header":null,"body":"","binary":"","form":null},"http_resp":{"status_code":0,"header":null,"body":"","status_message":"","proto_major":0,"proto_minor":0,"binary":""},"grpc_req":{"body":"{\"x\":1,\"y\":23}","method":"api.Adder.Add"},"grpc_resp":{"body":"{\"result\":81,\"data\":{\"name\":\"Fabio Di Gentanio\",\"team\":{\"name\":\"Ducati\",\"championships\":\"0\",\"points\":\"1001\"}}}","error":""},"deps":[{"name":"mongodb","type":"NO_SQL_DB","meta":{"InsertOneOptions":"[]","document":"x:1  y:23","name":"mongodb","operation":"InsertOne","type":"NO_SQL_DB"},"data":["LP+BAwEBD0luc2VydE9uZVJlc3VsdAH/ggABAQEKSW5zZXJ0ZWRJRAEQAAAAT/+CATNnby5tb25nb2RiLm9yZy9tb25nby1kcml2ZXIvYnNvbi9wcmltaXRpdmUuT2JqZWN0SUT/gwEBAQhPYmplY3RJRAH/hAABBgEYAAAY/4QUAAxj/8//qRkRZP/hI//V//IO/9sA","Cv+FBQEC/4gAAAAF/4YAAQE="]}],"test_case_path":"/Users/ritikjain/Desktop/go-practice/skp-workspace/go/grpc-example-app/keploy/tests","mock_path":"/Users/ritikjain/Desktop/go-practice/skp-workspace/go/grpc-example-app/keploy/mocks","mocks":[{"Version":"api.keploy.io/v1beta2","Kind":"Generic","Spec":{"Metadata":{"InsertOneOptions":"[]","document":"x:1  y:23","name":"mongodb","operation":"InsertOne","type":"NO_SQL_DB"},"Objects":[{"Type":"*mongo.InsertOneResult","Data":"LP+BAwEBD0luc2VydE9uZVJlc3VsdAH/ggABAQEKSW5zZXJ0ZWRJRAEQAAAAT/+CATNnby5tb25nb2RiLm9yZy9tb25nby1kcml2ZXIvYnNvbi9wcmltaXRpdmUuT2JqZWN0SUT/gwEBAQhPYmplY3RJRAH/hAABBgEYAAAY/4QUAAxj/8//qRkRZP/hI//V//IO/9sA"},{"Type":"*keploy.KError","Data":"Cv+FBQEC/4gAAAAF/4YAAQE="}]}}],"type":"gRPC"}'
        body_type: utf-8
    resp:
        status_code: 200
        header:
            Content-Type: application/json; charset=utf-8
            Vary: Origin
        body: |
            {"id":"415e534e-10f5-488b-a781-19a4bede11b5"}
        body_type: utf-8
        status_message: ""
        proto_major: 1
        proto_minor: 1
    objects:
        - type: error
          data: H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=
    mocks:
        - mock-1-0
    assertions:
        noise:
            - header.Content-Type
            - body.id
    created: 1674553625
`

	mock := models.Mock{}
	yaml.Unmarshal([]byte(data), &mock)
	var mocks []models.Mock
	mocks = append(mocks, mock)

	tests.Store(KeyForDummyDataForRead, mocks)

	mockExport := NewMockExportFS(true, tests)
	tc, err := mockExport.ReadByID(context.TODO(), "1", "dummy-tc-path", "dummy-mock-path")
	if err != nil {
		t.Fatal(err)
	}

	if tc.ID != "test-1" {
		t.Fatalf("ReadByID: expected test-1, got %s", tc.ID)
	}
}

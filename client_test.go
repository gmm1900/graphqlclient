package graphqlclient

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	testHeaderKey = "cust-header-key"
	testHeaderVal = "cust-header-val"
	testGraphqlQuery = `employees { id name }`
	testMockedRespBody = []byte(`
	{
		"data": {
			"employees": [
				{
					"id": 15,
					"name": "john"
				},
				{
					"id": 18,
					"name": "karl"
				}
			]
		},
		"errors": [
			{
				"message": "some error"
			}
		]
	}`)
)

func TestDo(t *testing.T) {
	type employee struct {
		ID int `json:"id"`
		Name string `json:"name"`
	}
	type employees struct {
		Employees []employee `json:employees`
	}
	type employeesResp struct {
		Data employees `json:"data"`
		Errors []map[string]interface{} `json:"errors"`
	}

	// GIVEN:
	graphqlSvr := MockGraphqlServer{
		MockedRespBody: [][]byte{
			testMockedRespBody,
		},
	}
	graphqlSvr.Start(t)
	defer graphqlSvr.Close()
	ctx := context.Background()

	// WHEN:
	graphqlClient := New(graphqlSvr.URL, nil, http.Header{testHeaderKey: []string{testHeaderVal}})
	var actualResp employeesResp
	err := graphqlClient.Do(ctx, Request{ Query: testGraphqlQuery }, &actualResp)

	// THEN:
	assert.NoError(t, err)

	// assert request header
	assert.Equal(t, testHeaderVal, graphqlSvr.CapturedReqHeaders[0].Get(testHeaderKey))

	// assert request body
	expectedReqBody := map[string]interface{}{
		"query": testGraphqlQuery,
	}
	assert.Equal(t, expectedReqBody, graphqlSvr.CapturedReqBody[0])

	// assert response received (mocked) is parsed to the desired object
	expectedResp := employeesResp {
		Data: employees{
			Employees: []employee{
				{ ID: 15, Name: "john" },
				{ ID: 18, Name: "karl" },
			},
		},
		Errors: []map[string]interface{}{
			{
				"message": "some error",
			},
		},
	}
	assert.Equal(t, expectedResp, actualResp)
}
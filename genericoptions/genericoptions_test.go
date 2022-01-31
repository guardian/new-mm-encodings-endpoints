package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestGenericOptionsAlways200(t *testing.T) {
	result, err := HandleEvent(context.Background(), &events.APIGatewayProxyRequest{})

	if err != nil {
		t.Error("HandleEvent returned an unexpected error: ", err)
	} else {
		if result.StatusCode != 200 {
			t.Errorf("HandleEvent returned a status of %d, should have been 200", result.StatusCode)
		}
		if allowOrigin, haveAllowOrigin := result.Headers["Access-Control-Allow-Origin"]; haveAllowOrigin {
			if allowOrigin != "*" {
				t.Errorf("HandleEvent returned unexpected Access-Control-Allow-Origin: %s", allowOrigin)
			}
		} else {
			t.Error("HandleEvent returned no Access-Control-Allow-Origin header")
		}
	}
}

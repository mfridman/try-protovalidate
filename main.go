package main

import (
	apiv1 "github.com/mfridman/try-protovalidate/gen/go/api/v1"
)

func main() {
	_ = apiv1.GetUserInfoRequest{}
	_ = apiv1.APIService_GetUserInfo_FullMethodName
}

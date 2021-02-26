package server

import (
	"context"
	"encoding/json"

	"felix.bs.com/felix/BeStrongerInGO/gRPC-Service/pkg/bapi"
	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
)

type TagServer struct{}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, err
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}

	return &tagList, nil
}

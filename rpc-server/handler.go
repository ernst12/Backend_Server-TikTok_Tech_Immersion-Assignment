package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/database"
	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

type Chat database.Chat

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	chat := &database.Chat{
		Chat:      req.Message.Chat,
		Text:      req.Message.Text,
		Sender:    req.Message.Sender,
		Send_time: time.Now().Unix(),
	}
	(*DB).Set(req.Message.Chat, chat)

	resp := rpc.NewSendResponse()
	resp.Code = 0
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = areYouLucky()
	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	message := &rpc.Message{
		Chat:     req.Message.Chat,
		Text:     req.Message.Text,
		Sender:   req.Message.Sender,
		SendTime: time.Now().Unix(),
	}
	err := (*DB).Append(req.Message.Chat, message)

	resp := rpc.NewSendResponse()
	if err != nil {
		resp.Code = 500
		log.Default().Print(err.Error())
	} else {
		resp.Code = 0
	}

	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	messages, err := (*DB).Get(req.Chat)
	if err != nil {
		resp.Code = 500
		log.Default().Print(err.Error())
		return resp, nil
	}
	lenMessages := len(messages)

	cursor := req.Cursor
	limit := req.Limit
	var reverse bool

	if limit == 0 {
		limit = 10 // default
	}

	if req.Reverse == nil {
		reverse = false // default
	} else {
		reverse = *req.Reverse
	}

	if reverse {
		for i, j := 0, lenMessages-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
	}

	// chat array slice. Eg. chats[messagesFrom : chatTo]
	messagesFrom := cursor * int64(limit)
	messagesTo := messagesFrom + int64(limit)
	if messagesTo > int64(lenMessages) {
		messagesTo = int64(lenMessages) // max
	}

	hasMore := messagesTo < int64(lenMessages)

	resp.Code = 0
	resp.Messages = messages[messagesFrom:messagesTo]
	resp.HasMore = &hasMore
	if hasMore {
		var nextCursor int64 = cursor + 1
		resp.NextCursor = &nextCursor
	}

	return resp, nil
}

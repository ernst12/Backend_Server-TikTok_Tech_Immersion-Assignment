package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/database"
	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	/*message := &rpc.Message{
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

	return resp, nil*/

	if err := validateSendRequest(req); err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	message := &database.Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}

	roomID, roomIdErr := getRoomID(req.Message.GetChat())
	if roomIdErr != nil {
		return nil, roomIdErr
	}

	err := rdb.SaveMessage(roomID, message)
	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	/*resp := rpc.NewPullResponse()

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

	return resp, nil*/

	roomID, err := getRoomID(req.GetChat())
    if err != nil {
       return nil, err
    }

    start := req.GetCursor()
    end := start + int64(req.GetLimit()) // did not minus 1 on purpose for hasMore check later on

    messages, err := rdb.GetMessagesByRoomID(roomID, start, end, req.GetReverse())
    if err != nil {
       return nil, err
    }

    respMessages := make([]*rpc.Message, 0)
    var counter int32 = 0
    var nextCursor int64 = 0
    hasMore := false
    for _, msg := range messages {
       if counter+1 > req.GetLimit() {
          // having extra value here means it has more data
          hasMore = true
          nextCursor = end
          break // do not return the last message
       }
       temp := &rpc.Message{
          Chat:     req.GetChat(),
          Text:     msg.Message,
          Sender:   msg.Sender,
          SendTime: msg.Timestamp,
       }
       respMessages = append(respMessages, temp)
       counter += 1
    }

    resp := rpc.NewPullResponse()
    resp.Messages = respMessages
    resp.Code = 0
    resp.Msg = "success"
    resp.HasMore = &hasMore
    resp.NextCursor = &nextCursor

    return resp, nil
}

func getRoomID(chat string) (string, error) {
	var roomID string

	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("invalid Chat ID '%s', should be in the format of user1:user2", chat)
		return "", err
	}

	sender1, sender2 := senders[0], senders[1]
	// Compare the sender and receiver alphabetically, and sort them asc to form the room ID
	if comp := strings.Compare(sender1, sender2); comp == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}

	return roomID, nil
}

func validateSendRequest(req *rpc.SendRequest) error {
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("invalid Chat ID '%s', should be in the format of user1:user2", req.Message.GetChat())
		return err
	}
	sender1, sender2 := senders[0], senders[1]

	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("sender '%s' not in the chat room", req.Message.GetSender())
		return err
	}

	return nil
}

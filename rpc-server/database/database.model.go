package database

//import ("github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc")

// will store the chats array in ascending order of the sendTime
/*type Database interface {
	Append(key string, value *rpc.Message) (error)
	Get(key string) ([]*rpc.Message, error)
	Delete(key string) (error)
}*/

type Message struct {
    Sender    string  `json:"sender"`
    Message   string  `json:"message"`
    Timestamp int64   `json:"timestamp"`
}

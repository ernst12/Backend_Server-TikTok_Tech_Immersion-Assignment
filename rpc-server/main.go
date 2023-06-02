package main

import (
	//"context"
	"fmt"
	"log"

	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/database"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	rpc "github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc/imservice"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
    rdb = &database.RedisClient{} // make the RedisClient with global visibility in the 'main' scope
)

func main() {
	//ctx := context.Background() // https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
	
	err := rdb.InitClient("redis:6379", "")
    if err != nil {
       errMsg := fmt.Sprintf("failed to init Redis client, err: %v", err)
       log.Fatal(errMsg)
    }

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

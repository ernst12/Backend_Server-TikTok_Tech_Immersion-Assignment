package main

import (
	"log"
	"github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/database"

	rpc "github.com/ernst12/Backend_Server-TikTok_Tech_Immersion-Assignment/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	_, redisErr := database.Factory("redis")
	//redisDb, redisErr := database.Factory("redis")
	if redisErr != nil {
		panic(redisErr)
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

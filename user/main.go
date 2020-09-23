package main


import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/dao"
	"github.com/longjoy/micro-go-course/endpoint"
	"github.com/longjoy/micro-go-course/redis"
	"github.com/longjoy/micro-go-course/service"
	"github.com/longjoy/micro-go-course/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)


func main(){
	fmt.Println("fuck")
	var (
		// 服务监听端口
		servicePort = flag.Int("service.port", 7999, "service port")
	)

	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)

	// 初始化 MySQL
	err := dao.InitMysql("localhost", "3306", "root", "ghost2111", "flaskr")
	if err != nil{
		log.Fatal(err)
	}
	// 初始化 redis
	err = redis.InitRedis("localhost", "6379", "")
	if err != nil{
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	// 定义注册与登录的处理函数
	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	// 加载路由处理
	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func(){
		errChan <- http.ListenAndServe(":" + strconv.Itoa(*servicePort), r)
	}()

	go func(){
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <- errChan
	log.Println(error)
}
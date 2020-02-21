package main

import (
	"context"
	pb "github.com/aau-network-security/haaukins/daemon/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

const (

	timeLayout = "2006-01-02 15:04:05"
)
type MD map[string][]string

type UserCredentials struct {
	// make variable name uppercase to make it accessible
	GrpcAddress string `yaml:"grpcendpoint"`
	GrpcPort 	string `yaml:"grpcport"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}


func main() {
	//setEnvVariables() // setting environment variables for test purposes
	f, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Debug().Msgf("Error while reading user credentials file")
	}
	var credentials UserCredentials
	err = yaml.Unmarshal(f, &credentials)
	if err != nil {
		log.Error().Msgf("Error while reading user credentials")
	}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(credentials.GrpcAddress+":"+credentials.GrpcPort, opts...)
	if err != nil {
		log.Debug().Msgf("Error : %s", err)
	}
	c := pb.NewDaemonClient(conn)

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	loginResp, err := c.LoginUser(ctx, &pb.LoginUserRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	})
	if err != nil {
		log.Debug().Msgf("Login user error ! %s", err)
	}
	md := metadata.New(map[string]string{"token": loginResp.GetToken()})
	contextWithMetaData := metadata.NewOutgoingContext(context.Background(), md)

	events, err := c.ListEvents(contextWithMetaData, &pb.ListEventsRequest{})
	if err != nil {
		log.Error().Msgf("Error on listing events %s", err)
	}
	for _, e := range events.Events {
		log.Debug().Msgf("Checking event %s", e.Name)
		t, err := time.Parse(timeLayout, e.FinishTime)
		if err != nil {
			log.Error().Msgf("Finish time parsing error! %s ", err)
		}

		if t.Before(time.Now()) {
			_, err := c.StopEvent(contextWithMetaData, &pb.StopEventRequest{
				Tag: e.Tag,
			})
			if err != nil {
				log.Fatal().Msgf("Error while stopping event %s", e.Name)
			}
			t, err := time.Parse(timeLayout,time.Now().String())
			if err !=nil {
				log.Warn().Msgf("time.Now() parsing error")
			}
			log.Debug().Msgf("Stopping event %s Expected finis time is  %s and stopping event at %s ", e.Name, e.FinishTime,t.String())
		}
	}
}

// this is used when testing the application
// grpc is not enabled !
//func setEnvVariables() {
//	os.Setenv("HKN_HOST","")
//	os.Setenv("HKN_PORT","")
//	os.Setenv("HKN_SSL_OFF", "true")
//}
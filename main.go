package main

import (
	"context"
	"flag"
	pb "github.com/aau-network-security/haaukins/daemon/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (

	timeLayout = "2006-01-02 15:04:05"
)
type MD map[string][]string

type UserCredentials struct {
	// make variable name uppercase to make it accessible
	User struct{
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"user,omitempty"`
	Grpc struct {
		Address string `yaml:"endpoint"`
		Port 	string `yaml:"port"`
	} `yaml:"grpc,omitempty"`
	TLS struct {
		Enabled bool `yaml:"enabled"`
		Certfile string `yaml:"certfile"`
	} `yaml:"tls,omitempty"`
}


func main() {
	//setEnvVariables() // setting environment variables for test purposes
	command := flag.String("command","stop", "In default stops scheduled events." )
	flag.Parse()
	log.Debug().Msgf("Flag has been set %s",*command)

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	loggerFile, err := os.OpenFile("schedulerlog", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal().Msgf("error opening file: %v", err)
	}

	defer loggerFile.Close()
	logger.Output(loggerFile)

	var userCredentials UserCredentials

	f, err := ioutil.ReadFile("/app/conf.yml") // mount the directory when running docker
	if err != nil {
		logger.Debug().Msgf("Error while reading user credentials file")
	}
	err = yaml.Unmarshal(f, &userCredentials)
	if err != nil {
		logger.Error().Msgf("Error while reading user credentials")
	}
	opts := []grpc.DialOption{

	}
	if userCredentials.TLS.Enabled {
		// no need to use certkey file which can be used in server side but not in client side
		creds, err := credentials.NewClientTLSFromFile(userCredentials.TLS.Certfile, "")
		if err != nil {
			log.Error().Msgf("Error while enabling certificates for GRPC, error : %s",err)
		}
		opts = append(opts,grpc.WithTransportCredentials(creds))

	} else {
		opts = append(opts,grpc.WithInsecure())
	}
	conn, err := grpc.Dial(userCredentials.Grpc.Address+":"+userCredentials.Grpc.Port, opts...)
	if err != nil {
		logger.Debug().Msgf("Error : %s", err)
	}
	c := pb.NewDaemonClient(conn)

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	loginResp, err := c.LoginUser(ctx, &pb.LoginUserRequest{
		Username: userCredentials.User.Username,
		Password: userCredentials.User.Password,
	})
	if err != nil {
		logger.Debug().Msgf("Login user error ! %s", err)
	}
	md := metadata.New(map[string]string{"token": loginResp.GetToken()})
	contextWithMetaData := metadata.NewOutgoingContext(context.Background(), md)

	events, err := c.ListEvents(contextWithMetaData, &pb.ListEventsRequest{})
	if err != nil {
		logger.Error().Msgf("Error on listing events %s", err)
	}

	if  events != nil {
		for _, e := range events.Events {
			finishTime, err := time.Parse(timeLayout, e.FinishTime)
			if err != nil {
				logger.Error().Msgf("Finish time parsing error! %s ", err)
			}
			if finishTime.After(time.Now()){
				log.Info().Msgf("Expected finish time for event: %s, is %s, skipping ... ",e.Name,e.FinishTime)
			}
			if finishTime.Before(time.Now()) {
				log.Debug().Msgf("Checking event %s", e.Name)
				if *command == "stop" {
					stopStream, err := c.StopEvent(contextWithMetaData, &pb.StopEventRequest{
						Tag: e.Tag,
					})
					if err != nil {
						log.Error().Msgf("Error on stopping event: %s ", err)
						return
					}
					for {
						_, err := stopStream.Recv()
						if err == io.EOF {
							break
						}

						if err != nil {
							log.Error().Msgf("Error: %s ", err)
							return
						}
					}
					logger.Debug().Msgf("Event %s is stopped ! ", e.Name)
				}
			}
			if *command == "start" {
			  _, err := c.StartEvent(contextWithMetaData,&pb.Empty{});
			  if err !=nil {
			  	 log.Error().Msgf("Error on starting event %s", err)
				}
			}
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

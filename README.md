## Scheduler 

This is a basic client command which will be used to automatically close all events in given time. 
It will run once in a day then checks events which has expired, then close those events no matter what is happening. 
There are some important parts to run this auto-cleaner for events: 
1. Specify volume of configuration file 
2. Bound certificates volume to docker image when `docker run ` command takes place

#### Example Run 

  *Building the image:*
 
- `docker build -t autokiller . `
 
 *Running the image* 

__In order to run stop command following command should be run :__
 
- `docker run -it --rm -v (path_of_conf_file):/app/conf.yml -v (path_of_certs):/certs autokiller -command stop`

__In order to run start command following command should be run :__

- `docker run -it --rm -v (path_of_conf_file):/app/conf.yml -v (path_of_certs):/certs autokiller -command start`


*Note*: In order to test it on your local environment ` --network=host` should be clarified when `grpc_endpoint` is `localhost`. 


#### Cronjob configuration on host side 

> Will gonna run at 02:00 midnight and stop the events which are expired. 
 
` 0 2 * * * docker run -it --rm -v (path_of_conf_file):/app/conf.yml -v (path_of_certs):/certs autokiller >> autokiller_logs `

__Docker container will gonna remove itself when it is done.__

#### Test Functionality 

- In order to test the functionality of this approach, `setEnvVariables` function should be set correctly
- According to environment variables which are specified in that function, it will make request to given grpc endpoint. 
- For testing functionality, no need to use TLS for GRPC endpoint, however it can be used as well.


#### Todo
- [x] Starting events automatically could be added (useful for booking functionality ...)
- [ ] Github actions for autodeploy and build should be integrated
- [ ] The functionality of this program should be tested on test environment properly


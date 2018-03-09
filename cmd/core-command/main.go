/*******************************************************************************
 * Copyright 2018 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * @microservice: core-command-go service
 * @author: Spencer Bull, Dell
 * @version: 0.5.0
 *******************************************************************************/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tsconn23/edgex-go"
	"github.com/tsconn23/edgex-go/core/command"
	"github.com/tsconn23/edgex-go/pkg/config"
	"github.com/tsconn23/edgex-go/pkg/heartbeat"
	logger "github.com/tsconn23/edgex-go/support/logging-client"
)

const (
	configFile string = "res/configuration.json"
)

var loggingClient logger.LoggingClient

func main() {
	start := time.Now()
	var (
		useConsul = flag.String("consul", "", "Should the service use consul?")
		useProfile = flag.String("profile", "default", "Specify a profile other than default.")
	)
	flag.Parse()

	//Read Configuration
	configuration := &command.ConfigurationStruct{}
	err := config.LoadFromFile(*useProfile, configuration)
	if err != nil {
		logBeforeTermination(err)
		return
	}

	//Determine if configuration should be overridden from Consul
	var consulMsg string
	if *useConsul == "y" {
		consulMsg = "Loading configuration from Consul..."
		err := command.ConnectToConsul(*configuration)
		if err != nil {
			logBeforeTermination(err)
			return //end program since user explicitly told us to use Consul.
		}
	} else {
		consulMsg = "Bypassing Consul configuration..."
	}

	// Setup Logging
	logTarget := setLoggingTarget(*configuration)
	var loggingClient = logger.NewClient(configuration.ApplicationName, configuration.EnableRemoteLogging, logTarget)

	loggingClient.Info(consulMsg)
	loggingClient.Info(fmt.Sprintf("Starting %s %s ", command.COMMANDSERVICENAME, edgex.Version))

	command.Init(*configuration, loggingClient)

	r := command.LoadRestRoutes()
	http.TimeoutHandler(nil, time.Millisecond*time.Duration(configuration.ServiceTimeout), "Request timed out")
	loggingClient.Info(configuration.AppOpenMsg, "")

	heartbeat.Start(configuration.HeartBeatMsg, configuration.HeartBeatTime, loggingClient)

	// Time it took to start service
	loggingClient.Info("Service started in: "+time.Since(start).String(), "")
	loggingClient.Info("Listening on port: "+strconv.Itoa(configuration.ServerPort), "")
	loggingClient.Error(http.ListenAndServe(":"+strconv.Itoa(configuration.ServerPort), r).Error())
}

func logBeforeTermination(err error) {
	loggingClient = logger.NewClient(command.COMMANDSERVICENAME, false, "")
	loggingClient.Error(err.Error())
}

func setLoggingTarget(conf command.ConfigurationStruct) string {
	logTarget := conf.LoggingRemoteURL
	if !conf.EnableRemoteLogging {
		return conf.LogFile
	}
	return logTarget
}

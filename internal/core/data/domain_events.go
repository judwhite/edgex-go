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
 *******************************************************************************/
package data

import "github.com/edgexfoundry/go-mod-core-contracts/clients/logger"

// An event indicating that a given device has just reported some data
type DeviceLastReported struct {
	DeviceName string
}

//An event indicating that the service associated with the device that just reported data is alive.
type DeviceServiceLastReported struct {
	DeviceName string
}

func initEventHandlers(loggingClient logger.LoggingClient) {
	go func() {
		for {
			select {
			case e, ok := <-chEvents:
				if ok {
					switch e.(type) {
					case DeviceLastReported:
						dlr := e.(DeviceLastReported)
						updateDeviceLastReportedConnected(dlr.DeviceName, loggingClient)
						break
					case DeviceServiceLastReported:
						dslr := e.(DeviceServiceLastReported)
						updateDeviceServiceLastReportedConnected(dslr.DeviceName, loggingClient)
						break
					}
				} else {
					return
				}
			}
		}
	}()
}

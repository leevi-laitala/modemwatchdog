package main

const ( // Time vars are in seconds
	pingInterval = 10

	powercycleDuration  = 15 // Duration of power being cut off
	powercycleRetries   = 10 // Times failed connection is retried
	powercycleRetryWait = 5  // Time between retries

	logSyslogAddr     = "" // Blank uses stdout
	logSyslogPort     = 514
	logSyslogProtocol = "tcp"
)

var pingUrls = []string{
	// "http://127.0.0.1:5000/gen_204", // Virtual python plug can respond to the ping, it will fail every third req
	"https://www.google.com/generate_204",
	"https://1.1.1.1/cdn-cgi/trace",
}

// Initialize smart plug
// "PythonPlug" is virtual plug that is "faulty", only every third request succeed
// "MockPlug" is virtual, and each request succeed
// "ShellyPlug" is interface for my Shelly Plug S Gen3
func initPlug() SmartPlug {
	//return ShellyPlug{
	//	apiUrl:     "http://192.168.1.215/rpc",
	//	apiTurnOff: "/Switch.Set?id=0&on=false",
	//	apiTurnOn:  "/Switch.Set?id=0&on=true",
	//	apiStatus:  "/Switch.GetStatus?id=0",
	//}
	return PythonPlug{
		apiUrl:     "http://127.0.0.1:5000/api",
		apiTurnOff: "/turnoff",
		apiTurnOn:  "/turnon",
		apiStatus:  "/status",
	}
	//
	// return MockPlug{}
}

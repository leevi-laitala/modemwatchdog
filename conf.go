package main

const ( // Time vars are in seconds
	pingInterval = 10

	powercycleDuration = 15 // Duration of power being cut off
	powercycleRetries = 10 // Times failed connection is retried
	powercycleRetryWait = 5 // Time between retries

	logSyslogAddr = "" // Blank uses stdout
	logSyslogPort = 514
	logSyslogProtocol = "tcp"
)

var pingUrls = []string{
//	"http://127.0.0.1:5000/gen_204", // Virtual python plug can respond to the ping, it will fail every third req
	"https://www.google.com/generate_204",
	"https://1.1.1.1/cdn-cgi/trace",
}

// Initialize smart plug
// "PythonPlug" is virtual plug that is "faulty", only every third request succeed
// "MockPlug" is virtual, and each request succeed
func initPlug() SmartPlug {
	return PythonPlug{
		apiUrl: "http://127.0.0.1:5000/api",
		apiTurnOff: "/turnoff",
		apiTurnOn: "/turnon",
		apiStatus: "/status",
	}
//	return MockPlug{}
}


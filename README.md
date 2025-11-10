# Modemwatchdog

Check internet connectivity and restart modem via a smart plug when connection is lost.

Creted for personal use as my ISP has shipped me modem with elusive fault, where it randomly drops internet connection and cannot recover from it without powercycle.

My band-aid solution is to use IoT smartplug that has exposed some form of API so it can be automated with a script. Ping some internet address, and when connection is lost, powercycle the modem by using the smartplug.

## Usage

Just execute built binary. No flags, args or CLI.

```
$ ./modemwatchdog
```

## Building

Program uses only golang standard library.

```
$ go build
```

## Configuration

The app does not have any CLI or flags. All configuration is done in the source code.

All configurable variables are located in `conf.go`

#### Logging

By default the app uses stdout for logs, but there is support for syslog. The implementation does not support TLS.

To enable syslog, fill server details into `conf.go`

#### Smart plug API

There is generic interface called `SmartPlug`, which has three functions: set power on, set power off and check power status.

To make this work with any smart plug, you need to make a new plug struct and implement these functions.

In `conf.go` there is a function `initPlug` that returns the common interface. Inside the function the struct of choice is be constructed.

##### Virtual plugs

There are two virtual "smart plugs": `PythonPlug` and `MockPlug`. PythonPlug is a separate program with exposed REST API. However only every third request is set to succeed. MockPlug is implemented in the `modemwatchdog` and is set to work every time.

Python plug can be executed with:

```
$ python -m venv runner
$ ./runner/bin/pip install flask
$ ./runner/bin/python virt-faulty-plug.py
```

#### Ping URLs

Ping is a wrong term here, as the app does not send ICMP packets. The internet connection check is done by sending HEAD requests.

There can be one or more URLs. If connection fails to each URL, the internet is deemed down, and powercycle starts.

#### Intervals and durations

Each time related duration is configured in seconds.

How often URLs are pinged is set in `pingInterval`, which is by default 15 seconds.

How many times smart plug API calls are retried upon failure is set in `powercycleRetries`, which is by default 10 times.

How much time is waited between retries is set with `powercycleRetryWait`, which is by default 5 seconds


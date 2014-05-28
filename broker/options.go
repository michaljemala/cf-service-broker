// Copyright 2014, The cf-service-broker Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that
// can be found in the LICENSE file.

package broker

import (
	"flag"
)

var Opts Options = Options{}

func init() {
	Opts.configure(flag.CommandLine)
}

type Options struct {
	Host     string
	Port     int
	Username string
	Password string
	Debug    bool
	LogFile  string
	Trace    bool
	PidFile  string
}

func (o *Options) configure(fs *flag.FlagSet) {
	fs.StringVar(&o.Host, "bh", "", "")
	fs.StringVar(&o.Host, "broker-host", "", "")

	fs.IntVar(&o.Port, "br", 9999, "")
	fs.IntVar(&o.Port, "broker-port", 9999, "")

	fs.StringVar(&o.Username, "bu", "admin", "")
	fs.StringVar(&o.Username, "broker-user", "admin", "")

	fs.StringVar(&o.Password, "bp", "secret", "")
	fs.StringVar(&o.Password, "broker-password", "secret", "")

	fs.BoolVar(&o.Debug, "D", false, "")

	fs.StringVar(&o.LogFile, "L", "", "")

	fs.BoolVar(&o.Trace, "V", false, "")

	fs.StringVar(&o.PidFile, "P", "", "")
}

var UsageStr = `
Broker Options:
    -bh, --host HOST                   Bind to HOST address (default: 0.0.0.0)
    -br, --port PORT                   Use PORT (default: 9999)
    -bu, --user USERNAME               User required to authenticate requests (default: admin)
    -bp, --pass PASSWORD               Password for the USERNAME user (default: secret)
    -D                                 Enable debugging output
    -L FILE                            File to redirect log output to
    -V                                 Trace the incoming service broker's HTTP requests
    -P FILE                            File to store broker's PID to
`

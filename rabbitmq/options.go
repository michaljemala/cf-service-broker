// Copyright 2014, The cf-service-broker Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that
// can be found in the LICENSE file.

package rabbitmq

import (
	"flag"
)

var Opts Options = Options{}

func init() {
	Opts.configure(flag.CommandLine)
}

type Options struct {
	Catalog  string
	Host     string
	Port     int
	MgmtHost string
	MgmtPort int
	MgmtUser string
	MgmtPass string
	Trace    bool // TODO: Create Rabbit-Hole PR to enable such tracing
}

func (o *Options) configure(fs *flag.FlagSet) {
	fs.StringVar(&o.Catalog, "c", "", "")
	fs.StringVar(&o.Catalog, "catalog", "", "")

	fs.StringVar(&o.Host, "rh", "127.0.0.1", "")
	fs.StringVar(&o.Host, "rabbit-host", "127.0.0.1", "")

	fs.IntVar(&o.Port, "rr", 5672, "")
	fs.IntVar(&o.Port, "rabbit-port", 5672, "")

	fs.StringVar(&o.MgmtHost, "rmh", "127.0.0.1", "")
	fs.StringVar(&o.MgmtHost, "rabbit-mgmt-host", "127.0.0.1", "")

	fs.IntVar(&o.MgmtPort, "rmr", 15672, "")
	fs.IntVar(&o.MgmtPort, "rabbit-mgmt-port", 15672, "")

	fs.StringVar(&o.MgmtUser, "rmu", "guest", "")
	fs.StringVar(&o.MgmtUser, "rabbit-mgmt-user", "guest", "")

	fs.StringVar(&o.MgmtPass, "rmp", "guest", "")
	fs.StringVar(&o.MgmtPass, "rabbit-mgmt-pass", "guest", "")

	fs.BoolVar(&o.Trace, "R", false, "")
}

var UsageStr = `
RabbitMQ Service Options:
    -c,   --catalog CATALOG            A file to load the RabbitMQ broker's catalog from
    -rh,  --rabbit-host HOST           Hostname of RabbitMQ server (default: 127.0.0.1)
    -rr,  --rabbit-port PORT           Port on which RabbitMQ server listens for messages  (default: 5672)
    -rmh, --rabbit-mgmt-host HOST      Hostname of RabbitMQ server (default: 127.0.0.1)
    -rmr, --rabbit-mgmt-port PORT      Port on which RabbitMQ server listens for management requests (default: 15672)
    -rmu, --rabbit-mgmt-user USERNAME  Username of the RabbitMQ server user with 'administrator' tag assigned (default: guest)
    -rmp, --rabbit-mgmt-pass PASSWORD  Password for the USERNAME user (default: guest)
    -R                                 Trace the outgoing RabbitMQ server management requests
`

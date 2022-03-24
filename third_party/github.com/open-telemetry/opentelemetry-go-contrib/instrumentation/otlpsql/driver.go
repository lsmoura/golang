// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"sync"

	"go.opentelemetry.io/otel/attribute"
)

var (
	regMu              sync.Mutex
	attrMissingContext = attribute.String("sql.warning", "missing upstream context")
	attrDeprecated     = attribute.String("sql.warning", "database driver uses deprecated features")

	// Compile time assertions
	_ driver.Driver = &otlpDriver{}
)

// Register initializes and registers our otelsql wrapped database driver
// identified by its driverName and using provided TraceOptions. On success it
// returns the generated driverName to use when calling sql.Open.
// It is possible to register multiple wrappers for the same database driver if
// needing different TraceOptions for different connections.
func Register(driverName string, options ...WrapperOption) (string, error) {
	return RegisterWithSource(driverName, "", options...)
}

// RegisterWithSource initializes and registers our otelsql wrapped database driver
// identified by its driverName, using provided TraceOptions.
// source is useful if some drivers do not accept the empty string when opening the DB.
// On success it returns the generated driverName to use when calling sql.Open.
// It is possible to register multiple wrappers for the same database driver if
// needing different TraceOptions for different connections.
func RegisterWithSource(driverName string, source string, options ...WrapperOption) (string, error) {
	// retrieve the driver implementation we need to wrap with instrumentation
	db, err := sql.Open(driverName, source)
	if err != nil {
		return "", err
	}
	dri := db.Driver()
	if err = db.Close(); err != nil {
		return "", err
	}

	regMu.Lock()
	defer regMu.Unlock()

	// Since we might want to register multiple otelsql drivers to have different
	// TraceOptions, but potentially the same underlying database driver, we
	// cycle through to find available driver names.
	driverName = driverName + "-otelsql-"
	for i := int64(0); i < 100; i++ {
		var (
			found   = false
			regName = driverName + strconv.FormatInt(i, 10)
		)
		for _, name := range sql.Drivers() {
			if name == regName {
				found = true
			}
		}
		if !found {
			sql.Register(regName, Wrap(dri, options...))
			return regName, nil
		}
	}
	return "", errors.New("unable to register driver, all slots have been taken")
}

// Wrap takes an SQL driver and wraps it with OpenCensus instrumentation.
func Wrap(d driver.Driver, options ...WrapperOption) driver.Driver {
	var o wrapper
	o.SetDefaults()
	o.ApplyOptions(options...)
	if o.InstanceName == "" {
		o.InstanceName = defaultInstanceName
	} else {
		o.DefaultAttributes = append(o.DefaultAttributes, attribute.String("sql.instance", o.InstanceName))
	}
	if o.QueryParams && !o.Query {
		o.QueryParams = false
	}
	return wrapDriver(d, o)
}

// Open implements driver.Driver
func (d otlpDriver) Open(name string) (driver.Conn, error) {
	c, err := d.parent.Open(name)
	if err != nil {
		return nil, err
	}
	return wrapConn(c, d.options), nil
}

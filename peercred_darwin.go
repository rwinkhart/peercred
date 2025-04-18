// Copyright (c) 2021-2025 AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package peercred

import (
	"fmt"
	"net"
	"strconv"

	"golang.org/x/sys/unix"
)

func get(c *net.UnixConn) (*Creds, error) {
	raw, err := c.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("SyscallConn: %w", err)
	}

	var cred *unix.Xucred
	var pid int
	cerr := raw.Control(func(fd uintptr) {
		cred, err = unix.GetsockoptXucred(int(fd),
			unix.SOL_LOCAL,
			unix.LOCAL_PEERCRED)
		if err != nil {
			err = fmt.Errorf("unix.GetsockoptXucred: %w", err)
			return
		}
		pid, err = unix.GetsockoptInt(int(fd),
			unix.SOL_LOCAL,
			unix.LOCAL_PEERPID)
		if err != nil {
			err = fmt.Errorf("unix.GetsockoptInt: %w", err)
		}
	})
	if cerr != nil {
		return nil, fmt.Errorf("raw.Control: %w", cerr)
	}
	if err != nil {
		return nil, err
	}
	return &Creds{
		PID: pid,
		UID: strconv.FormatUint(uint64(cred.Uid), 10),
	}, nil
}

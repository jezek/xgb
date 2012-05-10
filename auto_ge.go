package xgb

/*
	This file was generated by ge.xml on May 10 2012 12:39:33pm EDT.
	This file is automatically generated. Edit at your peril!
*/

// GeInit must be called before using the Generic Event Extension extension.
func (c *Conn) GeInit() error {
	reply, err := c.QueryExtension(23, "Generic Event Extension").Reply()
	switch {
	case err != nil:
		return err
	case !reply.Present:
		return errorf("No extension named Generic Event Extension could be found on on the server.")
	}

	c.extLock.Lock()
	c.extensions["Generic Event Extension"] = reply.MajorOpcode
	for evNum, fun := range newExtEventFuncs["Generic Event Extension"] {
		newEventFuncs[int(reply.FirstEvent)+evNum] = fun
	}
	for errNum, fun := range newExtErrorFuncs["Generic Event Extension"] {
		newErrorFuncs[int(reply.FirstError)+errNum] = fun
	}
	c.extLock.Unlock()

	return nil
}

func init() {
	newExtEventFuncs["Generic Event Extension"] = make(map[int]newEventFun)
	newExtErrorFuncs["Generic Event Extension"] = make(map[int]newErrorFun)
}

// Skipping definition for base type 'Int8'

// Skipping definition for base type 'Card16'

// Skipping definition for base type 'Char'

// Skipping definition for base type 'Card32'

// Skipping definition for base type 'Double'

// Skipping definition for base type 'Bool'

// Skipping definition for base type 'Float'

// Skipping definition for base type 'Card8'

// Skipping definition for base type 'Int16'

// Skipping definition for base type 'Int32'

// Skipping definition for base type 'Void'

// Skipping definition for base type 'Byte'

// Request GeQueryVersion
// size: 8
type GeQueryVersionCookie struct {
	*cookie
}

func (c *Conn) GeQueryVersion(ClientMajorVersion uint16, ClientMinorVersion uint16) GeQueryVersionCookie {
	cookie := c.newCookie(true, true)
	c.newRequest(c.geQueryVersionRequest(ClientMajorVersion, ClientMinorVersion), cookie)
	return GeQueryVersionCookie{cookie}
}

func (c *Conn) GeQueryVersionUnchecked(ClientMajorVersion uint16, ClientMinorVersion uint16) GeQueryVersionCookie {
	cookie := c.newCookie(false, true)
	c.newRequest(c.geQueryVersionRequest(ClientMajorVersion, ClientMinorVersion), cookie)
	return GeQueryVersionCookie{cookie}
}

// Request reply for GeQueryVersion
// size: 32
type GeQueryVersionReply struct {
	Sequence uint16
	Length   uint32
	// padding: 1 bytes
	MajorVersion uint16
	MinorVersion uint16
	// padding: 20 bytes
}

// Waits and reads reply data from request GeQueryVersion
func (cook GeQueryVersionCookie) Reply() (*GeQueryVersionReply, error) {
	buf, err := cook.reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return geQueryVersionReply(buf), nil
}

// Read reply into structure from buffer for GeQueryVersion
func geQueryVersionReply(buf []byte) *GeQueryVersionReply {
	v := new(GeQueryVersionReply)
	b := 1 // skip reply determinant

	b += 1 // padding

	v.Sequence = Get16(buf[b:])
	b += 2

	v.Length = Get32(buf[b:]) // 4-byte units
	b += 4

	v.MajorVersion = Get16(buf[b:])
	b += 2

	v.MinorVersion = Get16(buf[b:])
	b += 2

	b += 20 // padding

	return v
}

func (cook GeQueryVersionCookie) Check() error {
	return cook.check()
}

// Write request to wire for GeQueryVersion
func (c *Conn) geQueryVersionRequest(ClientMajorVersion uint16, ClientMinorVersion uint16) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	buf[b] = c.extensions["GENERIC EVENT EXTENSION"]
	b += 1

	buf[b] = 0 // request opcode
	b += 1

	Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	Put16(buf[b:], ClientMajorVersion)
	b += 2

	Put16(buf[b:], ClientMinorVersion)
	b += 2

	return buf
}

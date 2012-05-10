package xgb

/*
	This file was generated by res.xml on May 10 2012 12:39:33pm EDT.
	This file is automatically generated. Edit at your peril!
*/

// Imports are not necessary for XGB because everything is 
// in one package. They are still listed here for reference.
// import "xproto"

// ResInit must be called before using the X-Resource extension.
func (c *Conn) ResInit() error {
	reply, err := c.QueryExtension(10, "X-Resource").Reply()
	switch {
	case err != nil:
		return err
	case !reply.Present:
		return errorf("No extension named X-Resource could be found on on the server.")
	}

	c.extLock.Lock()
	c.extensions["X-Resource"] = reply.MajorOpcode
	for evNum, fun := range newExtEventFuncs["X-Resource"] {
		newEventFuncs[int(reply.FirstEvent)+evNum] = fun
	}
	for errNum, fun := range newExtErrorFuncs["X-Resource"] {
		newErrorFuncs[int(reply.FirstError)+errNum] = fun
	}
	c.extLock.Unlock()

	return nil
}

func init() {
	newExtEventFuncs["X-Resource"] = make(map[int]newEventFun)
	newExtErrorFuncs["X-Resource"] = make(map[int]newErrorFun)
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

// 'ResClient' struct definition
// Size: 8
type ResClient struct {
	ResourceBase uint32
	ResourceMask uint32
}

// Struct read ResClient
func ReadResClient(buf []byte, v *ResClient) int {
	b := 0

	v.ResourceBase = Get32(buf[b:])
	b += 4

	v.ResourceMask = Get32(buf[b:])
	b += 4

	return b
}

// Struct list read ResClient
func ReadResClientList(buf []byte, dest []ResClient) int {
	b := 0
	for i := 0; i < len(dest); i++ {
		dest[i] = ResClient{}
		b += ReadResClient(buf[b:], &dest[i])
	}
	return pad(b)
}

// Struct write ResClient
func (v ResClient) Bytes() []byte {
	buf := make([]byte, 8)
	b := 0

	Put32(buf[b:], v.ResourceBase)
	b += 4

	Put32(buf[b:], v.ResourceMask)
	b += 4

	return buf
}

// Write struct list ResClient
func ResClientListBytes(buf []byte, list []ResClient) int {
	b := 0
	var structBytes []byte
	for _, item := range list {
		structBytes = item.Bytes()
		copy(buf[b:], structBytes)
		b += pad(len(structBytes))
	}
	return b
}

// 'ResType' struct definition
// Size: 8
type ResType struct {
	ResourceType Atom
	Count        uint32
}

// Struct read ResType
func ReadResType(buf []byte, v *ResType) int {
	b := 0

	v.ResourceType = Atom(Get32(buf[b:]))
	b += 4

	v.Count = Get32(buf[b:])
	b += 4

	return b
}

// Struct list read ResType
func ReadResTypeList(buf []byte, dest []ResType) int {
	b := 0
	for i := 0; i < len(dest); i++ {
		dest[i] = ResType{}
		b += ReadResType(buf[b:], &dest[i])
	}
	return pad(b)
}

// Struct write ResType
func (v ResType) Bytes() []byte {
	buf := make([]byte, 8)
	b := 0

	Put32(buf[b:], uint32(v.ResourceType))
	b += 4

	Put32(buf[b:], v.Count)
	b += 4

	return buf
}

// Write struct list ResType
func ResTypeListBytes(buf []byte, list []ResType) int {
	b := 0
	var structBytes []byte
	for _, item := range list {
		structBytes = item.Bytes()
		copy(buf[b:], structBytes)
		b += pad(len(structBytes))
	}
	return b
}

// Request ResQueryVersion
// size: 8
type ResQueryVersionCookie struct {
	*cookie
}

func (c *Conn) ResQueryVersion(ClientMajor byte, ClientMinor byte) ResQueryVersionCookie {
	cookie := c.newCookie(true, true)
	c.newRequest(c.resQueryVersionRequest(ClientMajor, ClientMinor), cookie)
	return ResQueryVersionCookie{cookie}
}

func (c *Conn) ResQueryVersionUnchecked(ClientMajor byte, ClientMinor byte) ResQueryVersionCookie {
	cookie := c.newCookie(false, true)
	c.newRequest(c.resQueryVersionRequest(ClientMajor, ClientMinor), cookie)
	return ResQueryVersionCookie{cookie}
}

// Request reply for ResQueryVersion
// size: 12
type ResQueryVersionReply struct {
	Sequence uint16
	Length   uint32
	// padding: 1 bytes
	ServerMajor uint16
	ServerMinor uint16
}

// Waits and reads reply data from request ResQueryVersion
func (cook ResQueryVersionCookie) Reply() (*ResQueryVersionReply, error) {
	buf, err := cook.reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return resQueryVersionReply(buf), nil
}

// Read reply into structure from buffer for ResQueryVersion
func resQueryVersionReply(buf []byte) *ResQueryVersionReply {
	v := new(ResQueryVersionReply)
	b := 1 // skip reply determinant

	b += 1 // padding

	v.Sequence = Get16(buf[b:])
	b += 2

	v.Length = Get32(buf[b:]) // 4-byte units
	b += 4

	v.ServerMajor = Get16(buf[b:])
	b += 2

	v.ServerMinor = Get16(buf[b:])
	b += 2

	return v
}

func (cook ResQueryVersionCookie) Check() error {
	return cook.check()
}

// Write request to wire for ResQueryVersion
func (c *Conn) resQueryVersionRequest(ClientMajor byte, ClientMinor byte) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	buf[b] = c.extensions["X-RESOURCE"]
	b += 1

	buf[b] = 0 // request opcode
	b += 1

	Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	buf[b] = ClientMajor
	b += 1

	buf[b] = ClientMinor
	b += 1

	return buf
}

// Request ResQueryClients
// size: 4
type ResQueryClientsCookie struct {
	*cookie
}

func (c *Conn) ResQueryClients() ResQueryClientsCookie {
	cookie := c.newCookie(true, true)
	c.newRequest(c.resQueryClientsRequest(), cookie)
	return ResQueryClientsCookie{cookie}
}

func (c *Conn) ResQueryClientsUnchecked() ResQueryClientsCookie {
	cookie := c.newCookie(false, true)
	c.newRequest(c.resQueryClientsRequest(), cookie)
	return ResQueryClientsCookie{cookie}
}

// Request reply for ResQueryClients
// size: (32 + pad((int(NumClients) * 8)))
type ResQueryClientsReply struct {
	Sequence uint16
	Length   uint32
	// padding: 1 bytes
	NumClients uint32
	// padding: 20 bytes
	Clients []ResClient // size: pad((int(NumClients) * 8))
}

// Waits and reads reply data from request ResQueryClients
func (cook ResQueryClientsCookie) Reply() (*ResQueryClientsReply, error) {
	buf, err := cook.reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return resQueryClientsReply(buf), nil
}

// Read reply into structure from buffer for ResQueryClients
func resQueryClientsReply(buf []byte) *ResQueryClientsReply {
	v := new(ResQueryClientsReply)
	b := 1 // skip reply determinant

	b += 1 // padding

	v.Sequence = Get16(buf[b:])
	b += 2

	v.Length = Get32(buf[b:]) // 4-byte units
	b += 4

	v.NumClients = Get32(buf[b:])
	b += 4

	b += 20 // padding

	v.Clients = make([]ResClient, v.NumClients)
	b += ReadResClientList(buf[b:], v.Clients)

	return v
}

func (cook ResQueryClientsCookie) Check() error {
	return cook.check()
}

// Write request to wire for ResQueryClients
func (c *Conn) resQueryClientsRequest() []byte {
	size := 4
	b := 0
	buf := make([]byte, size)

	buf[b] = c.extensions["X-RESOURCE"]
	b += 1

	buf[b] = 1 // request opcode
	b += 1

	Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	return buf
}

// Request ResQueryClientResources
// size: 8
type ResQueryClientResourcesCookie struct {
	*cookie
}

func (c *Conn) ResQueryClientResources(Xid uint32) ResQueryClientResourcesCookie {
	cookie := c.newCookie(true, true)
	c.newRequest(c.resQueryClientResourcesRequest(Xid), cookie)
	return ResQueryClientResourcesCookie{cookie}
}

func (c *Conn) ResQueryClientResourcesUnchecked(Xid uint32) ResQueryClientResourcesCookie {
	cookie := c.newCookie(false, true)
	c.newRequest(c.resQueryClientResourcesRequest(Xid), cookie)
	return ResQueryClientResourcesCookie{cookie}
}

// Request reply for ResQueryClientResources
// size: (32 + pad((int(NumTypes) * 8)))
type ResQueryClientResourcesReply struct {
	Sequence uint16
	Length   uint32
	// padding: 1 bytes
	NumTypes uint32
	// padding: 20 bytes
	Types []ResType // size: pad((int(NumTypes) * 8))
}

// Waits and reads reply data from request ResQueryClientResources
func (cook ResQueryClientResourcesCookie) Reply() (*ResQueryClientResourcesReply, error) {
	buf, err := cook.reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return resQueryClientResourcesReply(buf), nil
}

// Read reply into structure from buffer for ResQueryClientResources
func resQueryClientResourcesReply(buf []byte) *ResQueryClientResourcesReply {
	v := new(ResQueryClientResourcesReply)
	b := 1 // skip reply determinant

	b += 1 // padding

	v.Sequence = Get16(buf[b:])
	b += 2

	v.Length = Get32(buf[b:]) // 4-byte units
	b += 4

	v.NumTypes = Get32(buf[b:])
	b += 4

	b += 20 // padding

	v.Types = make([]ResType, v.NumTypes)
	b += ReadResTypeList(buf[b:], v.Types)

	return v
}

func (cook ResQueryClientResourcesCookie) Check() error {
	return cook.check()
}

// Write request to wire for ResQueryClientResources
func (c *Conn) resQueryClientResourcesRequest(Xid uint32) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	buf[b] = c.extensions["X-RESOURCE"]
	b += 1

	buf[b] = 2 // request opcode
	b += 1

	Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	Put32(buf[b:], Xid)
	b += 4

	return buf
}

// Request ResQueryClientPixmapBytes
// size: 8
type ResQueryClientPixmapBytesCookie struct {
	*cookie
}

func (c *Conn) ResQueryClientPixmapBytes(Xid uint32) ResQueryClientPixmapBytesCookie {
	cookie := c.newCookie(true, true)
	c.newRequest(c.resQueryClientPixmapBytesRequest(Xid), cookie)
	return ResQueryClientPixmapBytesCookie{cookie}
}

func (c *Conn) ResQueryClientPixmapBytesUnchecked(Xid uint32) ResQueryClientPixmapBytesCookie {
	cookie := c.newCookie(false, true)
	c.newRequest(c.resQueryClientPixmapBytesRequest(Xid), cookie)
	return ResQueryClientPixmapBytesCookie{cookie}
}

// Request reply for ResQueryClientPixmapBytes
// size: 16
type ResQueryClientPixmapBytesReply struct {
	Sequence uint16
	Length   uint32
	// padding: 1 bytes
	Bytes         uint32
	BytesOverflow uint32
}

// Waits and reads reply data from request ResQueryClientPixmapBytes
func (cook ResQueryClientPixmapBytesCookie) Reply() (*ResQueryClientPixmapBytesReply, error) {
	buf, err := cook.reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return resQueryClientPixmapBytesReply(buf), nil
}

// Read reply into structure from buffer for ResQueryClientPixmapBytes
func resQueryClientPixmapBytesReply(buf []byte) *ResQueryClientPixmapBytesReply {
	v := new(ResQueryClientPixmapBytesReply)
	b := 1 // skip reply determinant

	b += 1 // padding

	v.Sequence = Get16(buf[b:])
	b += 2

	v.Length = Get32(buf[b:]) // 4-byte units
	b += 4

	v.Bytes = Get32(buf[b:])
	b += 4

	v.BytesOverflow = Get32(buf[b:])
	b += 4

	return v
}

func (cook ResQueryClientPixmapBytesCookie) Check() error {
	return cook.check()
}

// Write request to wire for ResQueryClientPixmapBytes
func (c *Conn) resQueryClientPixmapBytesRequest(Xid uint32) []byte {
	size := 8
	b := 0
	buf := make([]byte, size)

	buf[b] = c.extensions["X-RESOURCE"]
	b += 1

	buf[b] = 3 // request opcode
	b += 1

	Put16(buf[b:], uint16(size/4)) // write request size in 4-byte units
	b += 2

	Put32(buf[b:], Xid)
	b += 4

	return buf
}

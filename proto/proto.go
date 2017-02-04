package proto

// PacketHead the protocol header
type PacketHead struct {
	Len    int32
	Cmd    int32
	ReqSid string
	RspSid string
}

// PacketBody represents the protocol body
type PacketBody interface{}

// Proto protocol definition
type Proto struct {
	Head PacketHead
	Body PacketBody
}

// PacketHeadLength the head length
const PacketHeadLength int32 = 4

func New(cmd int32) {

}

// Parse parse the protocol
func Parse(msg []byte) *Proto {
	return &Proto{}
}

// Serialize serialize proto to send
func Serialize(p *Proto) []byte {
	return []byte{}
}

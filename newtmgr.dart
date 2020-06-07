//  Converting from Go to Dart: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/nmp.go

const NMP_HDR_SIZE = 8

type NmpHdr struct {
	Op    uint8 /* 3 bits of opcode */
	Flags uint8
	Len   uint16
	Group uint16
	Seq   uint8
	Id    uint8
}

void main() {
  print("Hello")
}
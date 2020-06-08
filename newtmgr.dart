//  Converting from Go to Dart: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/nmp.go

const NMP_HDR_SIZE = 8;

class NmpHdr {
	int Op;    //  uint8: 3 bits of opcode
	int Flags; //  uint8
	int Len;   //  uint16
	int Group; //  uint16
	int Seq;   //  uint8
	int Id;    //  uint8
}

class NmpMsg {
	NmpHdr Hdr;
	dynamic Body;  //  interface{}
  NmpMsg(this.Hdr, this.Body);
}

// Combine req + rsp.
mixin NmpReq {
	NmpHdr Hdr();
	void SetHdr(NmpHdr hdr);

	NmpMsg Msg();
}

mixin NmpRsp {
	NmpHdr Hdr();
	void SetHdr(NmpHdr msg);

	NmpMsg Msg();
}

class NmpBase {
	NmpHdr hdr;  //  `codec:"-"`
  
  NmpHdr Hdr() {
	  return hdr;
  }
  
  void SetHdr(NmpHdr h) {
	  hdr = h;
  }
}

NmpMsg MsgFromReq(NmpReq r) {
	return NmpMsg(
		r.Hdr(),
		r
	);
}

NmpMsg NewNmpMsg() {
	return NmpMsg(
    NmpHdr(),
    null
  );
}

func DecodeNmpHdr(data []byte) (*NmpHdr, error) {
	if len(data) < NMP_HDR_SIZE {
		return nil, fmt.Errorf(
			"Newtmgr request buffer too small %d bytes", len(data))
	}

	hdr := &NmpHdr{}

	hdr.Op = uint8(data[0])
	hdr.Flags = uint8(data[1])
	hdr.Len = binary.BigEndian.Uint16(data[2:4])
	hdr.Group = binary.BigEndian.Uint16(data[4:6])
	hdr.Seq = uint8(data[6])
	hdr.Id = uint8(data[7])

	return hdr, nil
}

func (hdr *NmpHdr) Bytes() []byte {
	buf := make([]byte, 0, NMP_HDR_SIZE)

	buf = append(buf, byte(hdr.Op))
	buf = append(buf, byte(hdr.Flags))

	u16b := make([]byte, 2)
	binary.BigEndian.PutUint16(u16b, hdr.Len)
	buf = append(buf, u16b...)

	binary.BigEndian.PutUint16(u16b, hdr.Group)
	buf = append(buf, u16b...)

	buf = append(buf, byte(hdr.Seq))
	buf = append(buf, byte(hdr.Id))

	return buf
}

func BodyBytes(body interface{}) ([]byte, error) {
	data := make([]byte, 0)

	enc := codec.NewEncoderBytes(&data, new(codec.CborHandle))
	if err := enc.Encode(body); err != nil {
		return nil, fmt.Errorf("Failed to encode message %s", err.Error())
	}

	log.Debugf("Encoded %+v to:\n%s", body, hex.Dump(data))

	return data, nil
}

func EncodeNmpPlain(nmr *NmpMsg) ([]byte, error) {
	bb, err := BodyBytes(nmr.Body)
	if err != nil {
		return nil, err
	}

	nmr.Hdr.Len = uint16(len(bb))

	hb := nmr.Hdr.Bytes()
	data := append(hb, bb...)

	log.Debugf("Encoded:\n%s", hex.Dump(data))

	return data, nil
}

func fillNmpReqWithSeq(req NmpReq, op uint8, group uint16, id uint8, seq uint8) {
	hdr := NmpHdr{
		Op:    op,
		Flags: 0,
		Len:   0,
		Group: group,
		Seq:   seq,
		Id:    id,
	}

	req.SetHdr(&hdr)
}

func fillNmpReq(req NmpReq, op uint8, group uint16, id uint8) {
	fillNmpReqWithSeq(req, op, group, id, nmxutil.NextNmpSeq())
}

/*
void main() {
  print("Hello")
}
*/

////////////////////////////////////////

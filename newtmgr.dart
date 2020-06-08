//  Converting from Go to Dart: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/nmp.go

const NMP_HDR_SIZE = 8;

/// SMP Header
class NmpHdr {
	int Op;    //  uint8: 3 bits of opcode
	int Flags; //  uint8
	int Len;   //  uint16
	int Group; //  uint16
	int Seq;   //  uint8
	int Id;    //  uint8
  
  /// Return this SMP Header as a list of bytes
  List<int> Bytes() /* []byte */ {
	  List<int> buf = [];  //  make([]byte, 0, NMP_HDR_SIZE);
    
    buf.add(this.Op);
	  buf.add(this.Flags);

	  List<int> u16b = binaryBigEndianPutUint16(this.Len);
	  buf.addAll(u16b);

	  u16b = binaryBigEndianPutUint16(this.Group);
	  buf.addAll(u16b);

	  buf.add(this.Seq);
	  buf.add(this.Id);
    assert(buf.length == NMP_HDR_SIZE);

	  return buf;
  }  
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

NmpHdr DecodeNmpHdr(List<int> data /* []byte */) {
	if (data.length < NMP_HDR_SIZE) {
    throw Exception(
      "Newtmgr request buffer too small ${data.length} bytes"
    );
	}

	var hdr = NmpHdr();

	hdr.Op    = data[0];  //  uint8
	hdr.Flags = data[1];  //  uint8
	hdr.Len   = binaryBigEndianUint16(data[2], data[3]);  //  binary.BigEndian.Uint16
	hdr.Group = binaryBigEndianUint16(data[4], data[5]);  //  binary.BigEndian.Uint16
	hdr.Seq   = data[6];  //  uint8
	hdr.Id    = data[7];  //  uint8

	return hdr;
}

List<int> BodyBytes(dynamic body /* interface{} */) /* []byte */ {
	var data = make([]byte, 0);

	var enc = codec.NewEncoderBytes(data, codec.CborHandle);
  try {
    enc.Encode(body);    
  } catch (err) {
    throw Exception("Failed to encode message ${err.Error()}");   
  }
  
	print("Encoded ${body} to:\n${ hexDump(data) }");

	return data;
}

/// Return byte array [a,b] as unsigned 16-bit int
int binaryBigEndianUint16(int a, int b) {
  return (a << 8) + b;
}

/// Return unsigned int u as big endian byte array
List<int> binaryBigEndianPutUint16(int u) {
  return [
    u >> 8,
    u & 0xff
  ];
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

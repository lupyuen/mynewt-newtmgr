//  Converting from Go to Dart: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/nmp.go
import 'package:cbor/cbor.dart' as cbor;  //  From https://pub.dev/packages/cbor

const NMP_HDR_SIZE = 8;

/// SMP Header
class NmpHdr {
	int Op;    //  uint8: 3 bits of opcode
	int Flags; //  uint8
	int Len;   //  uint16
	int Group; //  uint16
	int Seq;   //  uint8
	int Id;    //  uint8
  
  /// Construct an SMP Header
  NmpHdr(
    this.Op,    //  uint8: 3 bits of opcode
	  this.Flags, //  uint8
	  this.Len,   //  uint16
	  this.Group, //  uint16
	  this.Seq,   //  uint8
	  this.Id     //  uint8
  );
  
  /// Return this SMP Header as a list of bytes
  List<int> Bytes() {  //  Returns []byte
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
    NmpHdr(0, 0, 0, 0, 0, 0),
    null
  );
}

NmpHdr DecodeNmpHdr(List<int> data /* []byte */) {
	if (data.length < NMP_HDR_SIZE) {
    throw Exception(
      "Newtmgr request buffer too small ${data.length} bytes"
    );
	}

	final hdr = NmpHdr(
    data[0],  //  Op:    uint8
    data[1],  //  Flags: uint8
    binaryBigEndianUint16(data[2], data[3]),  //  Len: binary.BigEndian.Uint16
    binaryBigEndianUint16(data[4], data[5]),  //  Group: binary.BigEndian.Uint16
    data[6],  //  Seq:   uint8
    data[7],  //  Id:    uint8       
  );

	return hdr;
}

/*
/// Encode body with CBOR and return the byte array
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

/// Encode the SMP Message with CBOR and return the byte array
List<int> EncodeNmpPlain(NmpMsg nmr) /* []byte */ {
	final bb = BodyBytes(nmr.Body);

	nmr.Hdr.Len = bb.length;  //  uint16

	final hb = nmr.Hdr.Bytes();
	final data = [...hb, ...bb];

	print("Encoded:\n${ hexDump(data) }");

	return data;
}
*/

/// Init the SMP Request and set the sequence number
void fillNmpReqWithSeq(
  NmpReq req,
  int op,     //  uint8
  int group,  //  uint16
  int id,     //  uint8
  int seq     //  uint8
) {
	final hdr = NmpHdr(
		op,     //  Op
		0,      //  Flags
		0,      //  Len
		group,  //  Group
		seq,    //  Seq
		id      //  Id
	);

	req.SetHdr(hdr);
}

/*
/// Init the SMP Request and set the next sequence number
void fillNmpReq(
  NmpReq req, 
  int op,     //  uint8
  int group,  //  uint16
  int id      //  uint8
) {
	fillNmpReqWithSeq(
    req, 
    op, 
    group, 
    id, 
    nmxutil.NextNmpSeq()
  );
}
*/

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

/*
void main() {
  print("Hello");
}
*/

////////////////////////////////////////

/// An example of using the Map Builder class.
/// Map builder is used to build maps with complex values such as tag values, indefinite sequences
/// and the output of other list or map builders.
int main() {
  // Get our cbor instance, always do this,it correctly
  // initialises the decoder.
  final inst = cbor.Cbor();

  // Get our encoder
  final encoder = inst.encoder;

  // Encode some values
  encoder.writeArray(<int>[1, 2, 3]);
  encoder.writeFloat(67.89);
  encoder.writeInt(10);

  // Get our map builder
  final mapBuilder = cbor.MapBuilder.builder();

  // Add some map entries to the list.
  // Entries are added as a key followed by a value, this ordering is enforced.
  // Map keys can be integers or strings only, this is also enforced.
  mapBuilder.writeString('a'); // key
  mapBuilder.writeURI('a/ur1');
  mapBuilder.writeString('b'); // key
  mapBuilder.writeEpoch(1234567899);
  mapBuilder.writeString('c'); // key
  mapBuilder.writeDateTime('19/04/2020');

  // Get our built map output and add it to the encoding stream.
  // The key/value pairs must be balanced, i.e. you must end the map building with
  // a value else the getData method will throw an exception.
  // Use the addBuilderOutput method to add built output to the encoder.
  // You can use the addBuilderOutput method on the map builder to add
  // the output of other list or map builders to its encoding stream.
  final mapBuilderOutput = mapBuilder.getData();
  encoder.addBuilderOutput(mapBuilderOutput);

  // Add another value
  encoder.writeRegEx('^[12]g');

  // Decode ourselves and pretty print it.
  inst.decodeFromInput();
  print(inst.decodedPrettyPrint(false));

  // Finally to JSON
  print(inst.decodedToJSON());

  // JSON output is :-
  // [1,2,3],67.89,10,{"a":"a/ur1","b":1234567899,"c":"19/04/2020"},"^[12]g"

  return 0;
}

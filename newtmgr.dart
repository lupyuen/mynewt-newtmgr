import 'dart:math';
import 'package:cbor/cbor.dart' as cbor;  //  From https://pub.dev/packages/cbor

////////////////////////////////////////
//  nmxact/nmp/nmp.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/nmp.go

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

/// SMP Message
class NmpMsg {
	NmpHdr Hdr;
	NmpReq Body;  //  Previously interface{}

  /// Construct an SMP Message
  NmpMsg(this.Hdr, this.Body);
}

/// SMP Request Message
abstract class NmpReq {
	NmpHdr Hdr();
	void SetHdr(NmpHdr hdr);

	NmpMsg Msg();
}

/// SMP Response Message
abstract class NmpRsp {
	NmpHdr Hdr();
	void SetHdr(NmpHdr msg);

	NmpMsg Msg();
}

/// SMP Base Message
mixin NmpBase {
	NmpHdr hdr;  //  Will not be encoded: `codec:"-"`
  
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
    NextNmpSeq()  //  From nmxutil
  );
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

////////////////////////////////////////
//  nmxact/nmp/image.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/image.go

//////////////////////////////////////////////////////////////////////////////
// $state                                                                   //
//////////////////////////////////////////////////////////////////////////////

/* TODO
  type SplitStatus int

  const (
    NOT_APPLICABLE SplitStatus = iota
    NOT_MATCHING
    MATCHING
  )

  //  returns the enum as a string
  func (sm SplitStatus) String() string {
    names := map[SplitStatus]string{
      NOT_APPLICABLE: "N/A",
      NOT_MATCHING:   "non-matching",
      MATCHING:       "matching",
    }

    str := names[sm]
    if str == "" {
      return "Unknown!"
    }
    return str
  }

  type ImageStateEntry struct {
    NmpBase
    Image     int    `codec:"image"`
    Slot      int    `codec:"slot"`
    Version   string `codec:"version"`
    Hash      []byte `codec:"hash"`
    Bootable  bool   `codec:"bootable"`
    Pending   bool   `codec:"pending"`
    Confirmed bool   `codec:"confirmed"`
    Active    bool   `codec:"active"`
    Permanent bool   `codec:"permanent"`
  }
*/

class ImageStateReadReq 
  with NmpBase       //  Get and set SMP Message Header
  implements NmpReq  //  SMP Request Message  
{
	NmpBase base;  //  Will not be encoded: `codec:"-"`

  NmpMsg Msg() { return MsgFromReq(this); }
}

/* TODO
  type ImageStateWriteReq struct {
    NmpBase `codec:"-"`
    Hash    []byte `codec:"hash"`
    Confirm bool   `codec:"confirm"`
  }

  type ImageStateRsp struct {
    NmpBase
    Rc          int               `codec:"rc"`
    Images      []ImageStateEntry `codec:"images"`
    SplitStatus SplitStatus       `codec:"splitStatus"`
  }
*/

ImageStateReadReq NewImageStateReadReq() {
	var r = ImageStateReadReq();
	fillNmpReq(r, NMP_OP_READ, NMP_GROUP_IMAGE, NMP_ID_IMAGE_STATE);
	return r;
}

/* TODO
  func NewImageStateWriteReq() *ImageStateWriteReq {
    r := &ImageStateWriteReq{}
    fillNmpReq(r, NMP_OP_WRITE, NMP_GROUP_IMAGE, NMP_ID_IMAGE_STATE)
    return r
  }

  func (r *ImageStateWriteReq) Msg() *NmpMsg { return MsgFromReq(r) }

  func NewImageStateRsp() *ImageStateRsp {
    return &ImageStateRsp{}
  }

  func (r *ImageStateRsp) Msg() *NmpMsg { return MsgFromReq(r) }
*/

////////////////////////////////////////
//  nmxact/nmxutil/nmxutil.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmxutil/nmxutil.go

int nextNmpSeq = 0;  //  Previously uint8
bool nmpSeqBeenRead = false;

/// Return the next SMP Message Sequence Number, 0 to 255. The first number is random.
int NextNmpSeq() {  //  Returns uint8
	//  TODO: seqMutex.Lock()
	//  TODO: defer seqMutex.Unlock()

	if (!nmpSeqBeenRead) {
    //  First number is random
    var rng = new Random();
		nextNmpSeq = rng.nextInt(256);  //  Returns 0 to 255
		nmpSeqBeenRead = true;
	}

	final val = nextNmpSeq;
	nextNmpSeq = (nextNmpSeq + 1) % 256;
  assert(val >= 0 && val <= 255);
	return val;
}

////////////////////////////////////////
//  nmxact/nmp/defs.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/nmp/defs.go

const
	NMP_OP_READ      = 0,
	NMP_OP_READ_RSP  = 1,
	NMP_OP_WRITE     = 2,
	NMP_OP_WRITE_RSP = 3;

const
	NMP_ERR_OK       = 0,
	NMP_ERR_EUNKNOWN = 1,
	NMP_ERR_ENOMEM   = 2,
	NMP_ERR_EINVAL   = 3,
	NMP_ERR_ETIMEOUT = 4,
	NMP_ERR_ENOENT   = 5;

// First 64 groups are reserved for system level newtmgr commands.
// Per-user commands are then defined after group 64.

const
	NMP_GROUP_DEFAULT = 0,
	NMP_GROUP_IMAGE   = 1,
	NMP_GROUP_STAT    = 2,
	NMP_GROUP_CONFIG  = 3,
	NMP_GROUP_LOG     = 4,
	NMP_GROUP_CRASH   = 5,
	NMP_GROUP_SPLIT   = 6,
	NMP_GROUP_RUN     = 7,
	NMP_GROUP_FS      = 8,
	NMP_GROUP_SHELL   = 9,
	NMP_GROUP_PERUSER = 64;

// Default group (0).
const
	NMP_ID_DEF_ECHO           = 0,
	NMP_ID_DEF_CONS_ECHO_CTRL = 1,
	NMP_ID_DEF_TASKSTAT       = 2,
	NMP_ID_DEF_MPSTAT         = 3,
	NMP_ID_DEF_DATETIME_STR   = 4,
	NMP_ID_DEF_RESET          = 5;

// Image group (1).
const
	NMP_ID_IMAGE_STATE    = 0,
	NMP_ID_IMAGE_UPLOAD   = 1,
	NMP_ID_IMAGE_CORELIST = 3,
	NMP_ID_IMAGE_CORELOAD = 4,
	NMP_ID_IMAGE_ERASE    = 5;

// Stat group (2).
const
	NMP_ID_STAT_READ = 0,
	NMP_ID_STAT_LIST = 1;

// Config group (3).
const
	NMP_ID_CONFIG_VAL = 0;

// Log group (4).
const
	NMP_ID_LOG_SHOW        = 0,
	NMP_ID_LOG_CLEAR       = 1,
	NMP_ID_LOG_APPEND      = 2,
	NMP_ID_LOG_MODULE_LIST = 3,
	NMP_ID_LOG_LEVEL_LIST  = 4,
	NMP_ID_LOG_LIST        = 5;

// Crash group (5).
const
	NMP_ID_CRASH_TRIGGER = 0;

// Run group (7).
const
	NMP_ID_RUN_TEST = 0,
	NMP_ID_RUN_LIST = 1;

// File system group (8).
const
	NMP_ID_FS_FILE = 0;

// Shell group (8).
const
	NMP_ID_SHELL_EXEC = 0;

////////////////////////////////////////

/*
void main() {
  print("Hello");
}
*/

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
  final mapBuilder = cbor.MapBuilder
    .builder();

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
  final mapBuilderOutput = mapBuilder
    .getData();
  encoder
    .addBuilderOutput(mapBuilderOutput);

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

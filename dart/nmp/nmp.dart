/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
import 'package:cbor/cbor.dart' as cbor;               //  CBOR Encoder and Decoder. From https://pub.dev/packages/cbor
import 'package:typed_data/typed_data.dart' as typed;  //  Helpers for Byte Buffers. From https://pub.dev/packages/typed_data
import '../nmxutil/nmxutil.dart';

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
  typed.Uint8Buffer Bytes() {  //  Returns []byte
    var buf = typed.Uint8Buffer();  //  make([]byte, 0, NMP_HDR_SIZE);
    
    buf.add(this.Op);
    buf.add(this.Flags);

    typed.Uint8Buffer u16b = binaryBigEndianPutUint16(this.Len);
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
  void Encode(cbor.MapBuilder builder);
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

NmpHdr DecodeNmpHdr(typed.Uint8Buffer data /* []byte */) {
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

/// Encode SMP Request Body with CBOR and return the byte array
typed.Uint8Buffer BodyBytes(  //  Returns []byte
  NmpReq body  //  Previously interface{}
) {
  // Get our cbor instance, always do this, it correctly initialises the decoder.
  final inst = cbor.Cbor();

  // Get our encoder and map builder
  final encoder = inst.encoder;
  final mapBuilder = cbor.MapBuilder.builder();

  //  Encode the body as a CBOR map
  body.Encode(mapBuilder);
  final mapData = mapBuilder.getData();
  encoder.addBuilderOutput(mapData);

  //  Get the encoded body
  final data = inst.output.getData();

  //  Decode the encoded body and pretty print it
  inst.decodeFromInput();  //  print(inst.decodedPrettyPrint(false));
  final hdr = body.Hdr();
  print(
    "Encoded {NmpBase:{hdr:{"
    "Op:${ hdr.Op } "
    "Flags:${ hdr.Flags } "
    "Len:${ hdr.Len } "
    "Group:${ hdr.Group } "
    "Seq:${ hdr.Seq } "
    "Id:${ hdr.Id }}}} "
    "${ inst.decodedToJSON() } "
    "to:\n${ hexDump(data) }"
  );
  return data;
}

/// Encode the SMP Message with CBOR and return the byte array
typed.Uint8Buffer EncodeNmpPlain(NmpMsg nmr) {  //  Returns []byte
  final bb = BodyBytes(nmr.Body);

  nmr.Hdr.Len = bb.length;  //  uint16

  final hb = nmr.Hdr.Bytes();
  var data = typed.Uint8Buffer();
  data.addAll(hb);
  data.addAll(bb);

  print("Encoded:\n${ hexDump(data) }");
  return data;
}

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
typed.Uint8Buffer binaryBigEndianPutUint16(int u) {
  var data = typed.Uint8Buffer();
  data.add(u >> 8);
  data.add(u & 0xff);
  return data;
}

/// Return the buffer buf dumped as hex numbers
String hexDump(typed.Uint8Buffer buf) {
  return buf.map(
    (b) {
      return b.toRadixString(16).padLeft(2, "0");
    }
  ).join(" ");
}

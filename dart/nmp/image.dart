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
import 'defs.dart';
import 'nmp.dart';

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

  /// Encode the SMP Request fields to CBOR
  void Encode(cbor.MapBuilder builder) {
    // Add some map entries to the list.
    // Entries are added as a key followed by a value, this ordering is enforced.
    // Map keys can be integers or strings only, this is also enforced.
    // mapBuilder.writeString('a');   // key
    // mapBuilder.writeURI('a/ur1');  // value
    // mapBuilder.writeString('b');      // key
    // mapBuilder.writeEpoch(1234567899);// value
    // mapBuilder.writeString('c');           // key
    // mapBuilder.writeDateTime('19/04/2020');// value

    // Get our built map output and add it to the encoding stream.
    // The key/value pairs must be balanced, i.e. you must end the map building with
    // a value else the getData method will throw an exception.
    // Use the addBuilderOutput method to add built output to the encoder.
    // You can use the addBuilderOutput method on the map builder to add
    // the output of other list or map builders to its encoding stream.

    //  encoder.writeArray(<int>[1, 2, 3]);
    //  encoder.writeFloat(67.89);
    //  encoder.writeInt(10);
  }
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
//  TODO: Check response from PineTime

class ImageStateRsp 
  implements NmpRsp
{
  //  TODO
  int Rc;
  //  TODO
  NmpHdr Hdr() { return NmpHdr(0, 0, 0, 0, 0, 0); }
  //  TODO
  NmpMsg Msg() { return NmpMsg(null, null); }
  //  TODO
  void SetHdr(NmpHdr msg) {}
}

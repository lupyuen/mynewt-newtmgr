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
import 'nmp/image.dart';
import 'nmp/nmp.dart';
import 'sesn/sesn.dart';
import 'xact/image.dart';

////////////////////////////////////////
//  Send Simple Mgmt Protocol Command to PineTime over Bluetooth LE

void main() {
  composeRequest();
}

/// Compose a request to query firmware images on PineTime
typed.Uint8Buffer composeRequest() {
  //  Create the SMP Request
  final req = NewImageStateReadReq();

  //  Encode the SMP Message with CBOR
  final msg = req.Msg();
  final data = EncodeNmpPlain(msg);
  return data;
}

/// Query firmware images on PineTime
void testCommand() {
  //  Fetch the Bluetooth LE Session
  final s = GetSesn();

  //  Create the SMP Command
  final c = NewImageStateReadCmd();  //  Previously xact.NewImageStateReadCmd()

  //  TODO: Set the Bluetooth LE transmission options
  //  c.SetTxOptions(nmutil.TxOptions());

  //  Transmit the SMP Command
  final res = c.Run(s);

  //  TODO: Handle SMP Response
  //  final ires = res.ImageStateReadResult;  //  Previously xact.ImageStateReadResult
  //  imageStatePrintRsp(ires.Rsp);
}

/// Test the CBOR library for Encoding. Based on https://github.com/shamblett/cbor/blob/master/example/cbor_map_builder.dart
void testCborEncoding() {
  /// An example of using the Map Builder class.
  /// Map builder is used to build maps with complex values such as tag values, indefinite sequences
  /// and the output of other list or map builders.

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
  final mapData = mapBuilder
    .getData();
  encoder
    .addBuilderOutput(mapData);

  // Add another value
  encoder.writeRegEx('^[12]g');

  // Decode ourselves and pretty print it.
  inst.decodeFromInput();
  print(inst.decodedPrettyPrint(false));

  // Finally to JSON
  print(inst.decodedToJSON());

  // JSON output is :-
  // [1,2,3],67.89,10,{"a":"a/ur1","b":1234567899,"c":"19/04/2020"},"^[12]g"

  //  Get the encoded body
  final data = inst.output.getData();
  print("Encoded ${ inst.decodedToJSON() } to:\n${ hexDump(data) }");
}

/// Test the CBOR library for Decoding. Based on https://github.com/shamblett/cbor/blob/master/example/cbor_payload_decode.dart
void testCborDecoding() {
  /// An example of using the Map Builder class.
  /// Map builder is used to build maps with complex values such as tag values, indefinite sequences
  /// and the output of other list or map builders.

  final payload = <int>[
    0xbf, 0x66, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, // |.fimages|
    0x9f, 0xbf, 0x64, 0x73, 0x6c, 0x6f, 0x74, 0x00, 0x67, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, // |..dslot.gversion|
    0x65, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x64, 0x68, 0x61, 0x73, 0x68, 0x58, 0x20, 0x70, 0x3e, 0xbb, // |e1.0.0dhashX, 0xp>.|
    0xf8, 0x11, 0x45, 0x8b, 0x1f, 0xad, 0x18, 0x9e, 0x64, 0xe3, 0xa5, 0xe0, 0xf8, 0x09, 0xcb, 0xe6, // |..E.....d.......|
    0xba, 0xd8, 0x83, 0xc7, 0x6b, 0x3d, 0xd7, 0x12, 0x79, 0x1c, 0x82, 0x2f, 0xb5, 0x68, 0x62, 0x6f, // |....k=..y../.hbo|
    0x6f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0xf5, 0x67, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0xf4, // |otable.gpending.|
    0x69, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x65, 0x64, 0xf5, 0x66, 0x61, 0x63, 0x74, 0x69, // |iconfirmed.facti|
    0x76, 0x65, 0xf5, 0x69, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x65, 0x6e, 0x74, 0xf4, 0xff, 0xff, // |ve.ipermanent...|
    0x6b, 0x73, 0x70, 0x6c, 0x69, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x00, 0xff              // |ksplitStatus..| 
  ];

  // Get our cbor instance, always do this,it correctly
  // initialises the decoder.
  final inst = cbor.Cbor();

  final payloadBuffer = typed.Uint8Buffer();
  payloadBuffer.addAll(payload);

  // Decode from the buffer, you can also decode from the
  // int list if you prefer.
  inst.decodeFromBuffer(payloadBuffer);

  // Pretty print, note that these methods use [GetDecodedData] and will
  // thus build the payload buffer.
  // If you do not want to pretty print or use Json just get the list of
  // decoded data directly by calling [GetDecodedData()]
  print(inst.decodedPrettyPrint());

  // JSON, maps can only have string keys to decode to JSON
  print(inst.decodedToJSON());

  // Print CBOR and JSON sizes
  print(
    "CBOR size: ${payload.length} bytes\n"
    "JSON size: ${inst.decodedToJSON().length} bytes\n"
  );
}
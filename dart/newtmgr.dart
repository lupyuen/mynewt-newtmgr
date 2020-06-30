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

/// Compose a PineTime request
typed.Uint8Buffer composeRequest() {
  //  Create the SMP Request
  //  final req = NewImageStateReadReq();  //  Query firmware images on PineTime
  final req = NewImageUploadReq();  //  Upload firmware image to PineTime

  //  Set the request parameters
  req.ImageNum = 0;  //  ImageNum:0
  req.Off = 0;  //  Off:0 
  req.Len = 210196; //  Len:210196 
  //  DataSha:[145 42 176 2 222 234 181 71 24 73 54 150 175 150 35 35 57 129 146 149 39 236 233 86 236 10 222 79 48 211 23 39] 
  req.DataSha = typed.Uint8Buffer();
  req.DataSha.addAll([145, 42, 176, 2, 222, 234, 181, 71, 24, 73, 54, 150, 175, 150, 35, 35, 57, 129, 146, 149, 39, 236, 233, 86, 236, 10, 222, 79, 48, 211, 23, 39]);
  //  Upgrade:false 
  req.Upgrade = false;
  //  Data:[61 184 243 150 0 0 0 0 32 0 0 0 204 52 3 0 0 0 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 1 32 249 128 0 0 85 129 0 0 87 129 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 71 161 1 0 0 0 0 0 0 0 0 0 141 161 1 0 211 161 1 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 0 0 0 0 0 0 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 101 129 0 0 79 240 0 0 12 74 13 75 154 66 188 191 66 248 4 11 250 231 11 73 11 74 12 75 155 26 3 221 4 59 200 88 208 80 251 220 9 72 10 73 0 240 34 248 9 72 128 71 0 240 68 248 8 72 0 71 232 4 0 32 140 219 0 32 0 177 3 0 216 0 0 32 196 4 0 32 144 219 0 32 80 254 0 32 105 132 0 0 61 134 0 0 254 231 254 231 254 231 254 231 254 231 254 231 254 231 254 231 254 231 0 191 2 75 24 96 89 96 152 96 112 71 0 191 216 0 0 32 3 30 9 219 12 74 144 104 82 104 18 26 154 66 15 219 3 68 9 74 147 96 112 71 7 74 144 104 3 68 18 104 147 66 2 211 4 74 147 96 112 71 79 240 255 48 112 71 79 240 255 48 112 71 0 191 216 0 0 32 79 240 128 67 1 34 195 248 120 37 112 71 2 75 24 104 0 240 1 0 112 71 0 191 240 237 0 224 8 181 255 247 245 255 0 177 0 190 191 243 79 143 5 73 202 104 2 244 224 98 4 75 19 67 203 96 191 243 79 143 0 191 253 231 0 237 0 224 4 0 250 5 79 240 128 67]} 
  req.Data = typed.Uint8Buffer();
  req.Data.addAll([61, 184, 243, 150, 0, 0, 0, 0, 32, 0, 0, 0, 204, 52, 3, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 249, 128, 0, 0, 85, 129, 0, 0, 87, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 71, 161, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 141, 161, 1, 0, 211, 161, 1, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 101, 129, 0, 0, 79, 240, 0, 0, 12, 74, 13, 75, 154, 66, 188, 191, 66, 248, 4, 11, 250, 231, 11, 73, 11, 74, 12, 75, 155, 26, 3, 221, 4, 59, 200, 88, 208, 80, 251, 220, 9, 72, 10, 73, 0, 240, 34, 248, 9, 72, 128, 71, 0, 240, 68, 248, 8, 72, 0, 71, 232, 4, 0, 32, 140, 219, 0, 32, 0, 177, 3, 0, 216, 0, 0, 32, 196, 4, 0, 32, 144, 219, 0, 32, 80, 254, 0, 32, 105, 132, 0, 0, 61, 134, 0, 0, 254, 231, 254, 231, 254, 231, 254, 231, 254, 231, 254, 231, 254, 231, 254, 231, 254, 231, 0, 191, 2, 75, 24, 96, 89, 96, 152, 96, 112, 71, 0, 191, 216, 0, 0, 32, 3, 30, 9, 219, 12, 74, 144, 104, 82, 104, 18, 26, 154, 66, 15, 219, 3, 68, 9, 74, 147, 96, 112, 71, 7, 74, 144, 104, 3, 68, 18, 104, 147, 66, 2, 211, 4, 74, 147, 96, 112, 71, 79, 240, 255, 48, 112, 71, 79, 240, 255, 48, 112, 71, 0, 191, 216, 0, 0, 32, 79, 240, 128, 67, 1, 34, 195, 248, 120, 37, 112, 71, 2, 75, 24, 104, 0, 240, 1, 0, 112, 71, 0, 191, 240, 237, 0, 224, 8, 181, 255, 247, 245, 255, 0, 177, 0, 190, 191, 243, 79, 143, 5, 73, 202, 104, 2, 244, 224, 98, 4, 75, 19, 67, 203, 96, 191, 243, 79, 143, 0, 191, 253, 231, 0, 237, 0, 224, 4, 0, 250, 5, 79, 240, 128, 67]);

  //  Encode the SMP Message with CBOR and display the encoded message
  final msg = req.Msg();
  final data = EncodeNmpPlain(msg);
  return data;
}

/* Output:
> Executing task: cd dart && pub get && dart newtmgr.dart <
Resolving dependencies... 
Got dependencies!
Encoded {NmpBase:{hdr:{Op:2 Flags:0 Len:0 Group:1 Seq:152 Id:1}}} {"image":0,"off":0,"len":210196,"sha":[145,42,176,2,222,234,181,71,24,73,54,150,175,150,35,35,57,129,146,149,39,236,233,86,236,10,222,79,48,211,23,39],"upgrade":false,"data":[61,184,243,150,0,0,0,0,32,0,0,0,204,52,3,0,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,0,0,0,1,32,249,128,0,0,85,129,0,0,87,129,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,71,161,1,0,0,0,0,0,0,0,0,0,141,161,1,0,211,161,1,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,0,0,0,0,0,0,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,101,129,0,0,79,240,0,0,12,74,13,75,154,66,188,191,66,248,4,11,250,231,11,73,11,74,12,75,155,26,3,221,4,59,200,88,208,80,251,220,9,72,10,73,0,240,34,248,9,72,128,71,0,240,68,248,8,72,0,71,232,4,0,32,140,219,0,32,0,177,3,0,216,0,0,32,196,4,0,32,144,219,0,32,80,254,0,32,105,132,0,0,61,134,0,0,254,231,254,231,254,231,254,231,254,231,254,231,254,231,254,231,254,231,0,191,2,75,24,96,89,96,152,96,112,71,0,191,216,0,0,32,3,30,9,219,12,74,144,104,82,104,18,26,154,66,15,219,3,68,9,74,147,96,112,71,7,74,144,104,3,68,18,104,147,66,2,211,4,74,147,96,112,71,79,240,255,48,112,71,79,240,255,48,112,71,0,191,216,0,0,32,79,240,128,67,1,34,195,248,120,37,112,71,2,75,24,104,0,240,1,0,112,71,0,191,240,237,0,224,8,181,255,247,245,255,0,177,0,190,191,243,79,143,5,73,202,104,2,244,224,98,4,75,19,67,203,96,191,243,79,143,0,191,253,231,0,237,0,224,4,0,250,5,79,240,128,67]} to:
a6 65 69 6d 61 67 65 00 63 6f 66 66 00 63 6c 65 6e 1a 00 03 35 14 63 73 68 61 98 20 18 91 18 2a 18 b0 02 18 de 18 ea 18 b5 18 47 18 18 18 49 18 36 18 96 18 af 18 96 18 23 18 23 18 39 18 81 18 92 18 95 18 27 18 ec 18 e9 18 56 18 ec 0a 18 de 18 4f 18 30 18 d3 17 18 27 67 75 70 67 72 61 64 65 f4 64 64 61 74 61 99 02 00 18 3d 18 b8 18 f3 18 96 00 00 00 00 18 20 00 00 00 18 cc 18 34 03 00 00 00 00 00 01 01 00 00 00 00 00 00 00 00 00 00 00 00 01 18 20 18 f9 18 80 00 00 18 55 18 81 00 00 18 57 18 81 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 18 47 18 a1 01 00 00 00 00 00 00 00 00 00 18 8d 18 a1 01 00 18 d3 18 a1 01 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 00 00 00 00 00 00 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 4f 18 f0 00 00 0c 18 4a 0d 18 4b 18 9a 18 42 18 bc 18 bf 18 42 18 f8 04 0b 18 fa 18 e7 0b 18 49 0b 18 4a 0c 18 4b 18 9b 18 1a 03 18 dd 04 18 3b 18 c8 18 58 18 d0 18 50 18 fb 18 dc 09 18 48 0a 18 49 00 18 f0 18 22 18 f8 09 18 48 18 80 18 47 00 18 f0 18 44 18 f8 08 18 48 00 18 47 18 e8 04 00 18 20 18 8c 18 db 00 18 20 00 18 b1 03 00 18 d8 00 00 18 20 18 c4 04 00 18 20 18 90 18 db 00 18 20 18 50 18 fe 00 18 20 18 69 18 84 00 00 18 3d 18 86 00 00 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 00 18 bf 02 18 4b 18 18 18 60 18 59 18 60 18 98 18 60 18 70 18 47 00 18 bf 18 d8 00 00 18 20 03 18 1e 09 18 db 0c 18 4a 18 90 18 68 18 52 18 68 12 18 1a 18 9a 18 42 0f 18 db 03 18 44 09 18 4a 18 93 18 60 18 70 18 47 07 18 4a 18 90 18 68 03 18 44 12 18 68 18 93 18 42 02 18 d3 04 18 4a 18 93 18 60 18 70 18 47 18 4f 18 f0 18 ff 18 30 18 70 18 47 18 4f 18 f0 18 ff 18 30 18 70 18 47 00 18 bf 18 d8 00 00 18 20 18 4f 18 f0 18 80 18 43 01 18 22 18 c3 18 f8 18 78 18 25 18 70 18 47 02 18 4b 18 18 18 68 00 18 f0 01 00 18 70 18 47 00 18 bf 18 f0 18 ed 00 18 e0 08 18 b5 18 ff 18 f7 18 f5 18 ff 00 18 b1 00 18 be 18 bf 18 f3 18 4f 18 8f 05 18 49 18 ca 18 68 02 18 f4 18 e0 18 62 04 18 4b 13 18 43 18 cb 18 60 18 bf 18 f3 18 4f 18 8f 00 18 bf 18 fd 18 e7 00 18 ed 00 18 e0 04 00 18 fa 05 18 4f 18 f0 18 80 18 43
Encoded:
02 00 03 84 00 01 98 01 a6 65 69 6d 61 67 65 00 63 6f 66 66 00 63 6c 65 6e 1a 00 03 35 14 63 73 68 61 98 20 18 91 18 2a 18 b0 02 18 de 18 ea 18 b5 18 47 18 18 18 49 18 36 18 96 18 af 18 96 18 23 18 23 18 39 18 81 18 92 18 95 18 27 18 ec 18 e9 18 56 18 ec 0a 18 de 18 4f 18 30 18 d3 17 18 27 67 75 70 67 72 61 64 65 f4 64 64 61 74 61 99 02 00 18 3d 18 b8 18 f3 18 96 00 00 00 00 18 20 00 00 00 18 cc 18 34 03 00 00 00 00 00 01 01 00 00 00 00 00 00 00 00 00 00 00 00 01 18 20 18 f9 18 80 00 00 18 55 18 81 00 00 18 57 18 81 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 18 47 18 a1 01 00 00 00 00 00 00 00 00 00 18 8d 18 a1 01 00 18 d3 18 a1 01 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 00 00 00 00 00 00 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 65 18 81 00 00 18 4f 18 f0 00 00 0c 18 4a 0d 18 4b 18 9a 18 42 18 bc 18 bf 18 42 18 f8 04 0b 18 fa 18 e7 0b 18 49 0b 18 4a 0c 18 4b 18 9b 18 1a 03 18 dd 04 18 3b 18 c8 18 58 18 d0 18 50 18 fb 18 dc 09 18 48 0a 18 49 00 18 f0 18 22 18 f8 09 18 48 18 80 18 47 00 18 f0 18 44 18 f8 08 18 48 00 18 47 18 e8 04 00 18 20 18 8c 18 db 00 18 20 00 18 b1 03 00 18 d8 00 00 18 20 18 c4 04 00 18 20 18 90 18 db 00 18 20 18 50 18 fe 00 18 20 18 69 18 84 00 00 18 3d 18 86 00 00 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 18 fe 18 e7 00 18 bf 02 18 4b 18 18 18 60 18 59 18 60 18 98 18 60 18 70 18 47 00 18 bf 18 d8 00 00 18 20 03 18 1e 09 18 db 0c 18 4a 18 90 18 68 18 52 18 68 12 18 1a 18 9a 18 42 0f 18 db 03 18 44 09 18 4a 18 93 18 60 18 70 18 47 07 18 4a 18 90 18 68 03 18 44 12 18 68 18 93 18 42 02 18 d3 04 18 4a 18 93 18 60 18 70 18 47 18 4f 18 f0 18 ff 18 30 18 70 18 47 18 4f 18 f0 18 ff 18 30 18 70 18 47 00 18 bf 18 d8 00 00 18 20 18 4f 18 f0 18 80 18 43 01 18 22 18 c3 18 f8 18 78 18 25 18 70 18 47 02 18 4b 18 18 18 68 00 18 f0 01 00 18 70 18 47 00 18 bf 18 f0 18 ed 00 18 e0 08 18 b5 18 ff 18 f7 18 f5 18 ff 00 18 b1 00 18 be 18 bf 18 f3 18 4f 18 8f 05 18 49 18 ca 18 68 02 18 f4 18 e0 18 62 04 18 4b 13 18 43 18 cb 18 60 18 bf 18 f3 18 4f 18 8f 00 18 bf 18 fd 18 e7 00 18 ed 00 18 e0 04 00 18 fa 05 18 4f 18 f0 18 80 18 43
*/

//  Encoding according to newtmgr with NmpBase:{hdr:{Op:2 Flags:0 Len:0 Group:1 Seq:67 Id:1}}...`
//  00000000  a5 64 64 61 74 61 59 02  00 3d b8 f3 96 00 00 00  |.ddataY..=......|
//  00000010  00 20 00 00 00 cc 34 03  00 00 00 00 00 01 01 00  |. ....4.........|
//  00000020  00 00 00 00 00 00 00 00  00 00 00 01 20 f9 80 00  |............ ...|
//  00000030  00 55 81 00 00 57 81 00  00 00 00 00 00 00 00 00  |.U...W..........|
//  00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
//  00000050  00 00 00 00 00 47 a1 01  00 00 00 00 00 00 00 00  |.....G..........|
//  00000060  00 8d a1 01 00 d3 a1 01  00 65 81 00 00 65 81 00  |.........e...e..|
//  00000070  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  00000080  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  00000090  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  000000a0  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  000000b0  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  000000c0  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  000000d0  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  000000e0  00 00 00 00 00 00 00 00  00 65 81 00 00 65 81 00  |.........e...e..|
//  000000f0  00 65 81 00 00 65 81 00  00 65 81 00 00 65 81 00  |.e...e...e...e..|
//  00000100  00 4f f0 00 00 0c 4a 0d  4b 9a 42 bc bf 42 f8 04  |.O....J.K.B..B..|
//  00000110  0b fa e7 0b 49 0b 4a 0c  4b 9b 1a 03 dd 04 3b c8  |....I.J.K.....;.|
//  00000120  58 d0 50 fb dc 09 48 0a  49 00 f0 22 f8 09 48 80  |X.P...H.I..\"..H.|
//  00000130  47 00 f0 44 f8 08 48 00  47 e8 04 00 20 8c db 00  |G..D..H.G... ...|
//  00000140  20 00 b1 03 00 d8 00 00  20 c4 04 00 20 90 db 00  | ....... ... ...|
//  00000150  20 50 fe 00 20 69 84 00  00 3d 86 00 00 fe e7 fe  | P.. i...=......|
//  00000160  e7 fe e7 fe e7 fe e7 fe  e7 fe e7 fe e7 fe e7 00  |................|
//  00000170  bf 02 4b 18 60 59 60 98  60 70 47 00 bf d8 00 00  |..K.`Y`.`pG.....|
//  00000180  20 03 1e 09 db 0c 4a 90  68 52 68 12 1a 9a 42 0f  | .....J.hRh...B.|
//  00000190  db 03 44 09 4a 93 60 70  47 07 4a 90 68 03 44 12  |..D.J.`pG.J.h.D.|
//  000001a0  68 93 42 02 d3 04 4a 93  60 70 47 4f f0 ff 30 70  |h.B...J.`pGO..0p|
//  000001b0  47 4f f0 ff 30 70 47 00  bf d8 00 00 20 4f f0 80  |GO..0pG..... O..|
//  000001c0  43 01 22 c3 f8 78 25 70  47 02 4b 18 68 00 f0 01  |C.\"..x%pG.K.h...|
//  000001d0  00 70 47 00 bf f0 ed 00  e0 08 b5 ff f7 f5 ff 00  |.pG.............|
//  000001e0  b1 00 be bf f3 4f 8f 05  49 ca 68 02 f4 e0 62 04  |.....O..I.h...b.|
//  000001f0  4b 13 43 cb 60 bf f3 4f  8f 00 bf fd e7 00 ed 00  |K.C.`..O........|
//  00000200  e0 04 00 fa 05 4f f0 80  43 65 69 6d 61 67 65 00  |.....O..Ceimage.|
//  00000210  63 6c 65 6e 1a 00 03 35  14 63 6f 66 66 00 63 73  |clen...5.coff.cs|
//  00000220  68 61 58 20 91 2a b0 02  de ea b5 47 18 49 36 96  |haX .*.....G.I6.|
//  00000230  af 96 23 23 39 81 92 95  27 ec e9 56 ec 0a de 4f  |..##9...'..V...O|
//  00000240  30 d3 17 27                                       |0..'|

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
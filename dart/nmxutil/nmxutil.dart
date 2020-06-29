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
import 'dart:math';
import 'package:cbor/cbor.dart' as cbor;               //  CBOR Encoder and Decoder. From https://pub.dev/packages/cbor
import 'package:typed_data/typed_data.dart' as typed;  //  Helpers for Byte Buffers. From https://pub.dev/packages/typed_data

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

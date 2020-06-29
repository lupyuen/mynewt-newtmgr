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
import '../sesn/sesn.dart';
import '../nmp/nmp.dart';
import '../nmp/image.dart';
import './cmd.dart';

////////////////////////////////////////
//  nmxact/xact/xact.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/xact/xact.go

/// Transmit an SMP Request and get the SMP Response
NmpRsp txReq(  //  Returns nmp.NmpRsp
  Sesn s,      //  Previously sesn.Sesn
  NmpMsg m,    //  Previously nmp.NmpMsg
  CmdBase c
) {
  //  TODO: assert(c != null);
  if (c != null) {  //  TODO: Should not be null
    if (c.abortErr != null) {
      throw c.abortErr;
    }
    c.curNmpSeq = m.Hdr.Seq;
    c.curSesn = s;
  }

	//  TODO: final rsp = sesn.TxRxMgmt(s, m, c.TxOptions());
  final rsp = ImageStateRsp();
  final data = EncodeNmpPlain(m);

  if (c != null) {  //  TODO: Should not be null
    c.curNmpSeq = 0;
    c.curSesn = null;
  }
	return rsp;
}

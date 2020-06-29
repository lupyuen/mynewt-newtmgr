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
import '../nmp/image.dart';
import '../xact/xact.dart';
import 'cmd.dart';

////////////////////////////////////////
//  nmxact/xact/image.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/xact/image.go

//////////////////////////////////////////////////////////////////////////////
// $state read                                                              //
//////////////////////////////////////////////////////////////////////////////

class ImageStateReadCmd 
  with CmdBase 
  implements Cmd 
{
  CmdBase base;

  //  TODO: ImageStateReadCmd(this.base);

  Result Run(
    Sesn s  //  Previously sesn.Sesn
  ) {
    final r = NewImageStateReadReq();  //  Previously nmp.NewImageStateReadReq()

    final rsp = txReq(s, r.Msg(), this.base);
    
    //  TODO: final srsp = rsp.ImageStateRsp;  //  Previously nmp.ImageStateRsp

    var res = newImageStateReadResult();
    //  TODO: res.Rsp = srsp;
    return res;
  }
}

class ImageStateReadResult implements Result {
  ImageStateRsp Rsp;  //  Previously nmp.ImageStateRsp

  int Status() {
    return this.Rsp.Rc;
  }
}

ImageStateReadCmd NewImageStateReadCmd() {
  return ImageStateReadCmd(
    //  TODO: NewCmdBase()
  );
}

ImageStateReadResult newImageStateReadResult() {
  return ImageStateReadResult();
}

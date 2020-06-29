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

////////////////////////////////////////
//  nmxact/xact/cmd.go
//  Converted from Go: https://github.com/lupyuen/mynewt-newtmgr/blob/master/nmxact/xact/cmd.go

/// Result of an SMP operation
abstract class Result {
	int Status();
}

/// SMP Command
abstract class Cmd {
	/// Transmits request and listens for response; blocking.
	Result Run(Sesn s);                //  Previously sesn.Sesn
	void Abort();

	//  TxOptions TxOptions();             //  Previously sesn.TxOptions
	void SetTxOptions(TxOptions opt);  //  Previously sesn.TxOptions
}

/// Base Class for SMP Command
mixin CmdBase {
	TxOptions txOptions;  //  Previously sesn.TxOptions
	int curNmpSeq;        //  Previously uint8
	Sesn curSesn;         //  Previously sesn.Sesn
	Exception abortErr;   // Previously error

  /// Constructor
  //  CmdBase(this.txOptions);

  /*
  TxOptions TxOptions() {  //  Previously sesn.TxOptions
    return this.txOptions;
  }
  */

  void SetTxOptions(
    TxOptions opt  //  Previously sesn.TxOptions
  ) {
    this.txOptions = opt;
  }

  void Abort() {
    if (this.curSesn != null) {
        //  TODO: this.curSesn.AbortRx(this.curNmpSeq);
    }
    this.abortErr = Exception("Command aborted");
  }
}

/*
CmdBase NewCmdBase() {
	return CmdBase(
		NewTxOptions()  //  Previously sesn.NewTxOptions
  );
}
*/

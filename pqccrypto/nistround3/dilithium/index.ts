import { getSignSysIF } from "../../../crypto/signSysIF";
const Dilithium2Addon = require("./build/Release/Dilithium2RefNistRound3.node");
const Dilithium3Addon = require("./build/Release/Dilithium3RefNistRound3.node");
const Dilithium5Addon = require("./build/Release/Dilithium5RefNistRound3.node");


const Dilithium2IF = getSignSysIF();
Dilithium2IF.signCore = Dilithium2Addon;
Dilithium2IF.signSysName = 'Nist_round3_Dilithium2';

class Dilithium2 extends Dilithium2IF {
	constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
		super(privateKey, publicKey, encryptedPrivateKeyMode);
	}
}


const Dilithium3IF = getSignSysIF();
Dilithium3IF.signCore = Dilithium3Addon;
Dilithium3IF.signSysName = 'Nist_round3_Dilithium3';

class Dilithium3 extends Dilithium3IF {
	constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
		super(privateKey, publicKey, encryptedPrivateKeyMode);
	}
}


const Dilithium5IF = getSignSysIF();
Dilithium5IF.signCore = Dilithium5Addon;
Dilithium5IF.signSysName = 'Nist_round3_Dilithium5';

class Dilithium5 extends Dilithium5IF {
	constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
		super(privateKey, publicKey, encryptedPrivateKeyMode);
	}
}

export { Dilithium2, Dilithium3, Dilithium5 };

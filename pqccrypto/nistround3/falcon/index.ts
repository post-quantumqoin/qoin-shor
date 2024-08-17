import { getSignSysIF } from "../../../crypto/signSysIF";

function getFalconIF() {
	let signSysIF = getSignSysIF();
	return class FalconIF extends signSysIF {
		constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
			super(privateKey, publicKey, encryptedPrivateKeyMode);
		}

		sign(data: Buffer, aesKey?: Buffer): false | Buffer {
			let signMsg = super.sign(data, aesKey);
			if(!signMsg) {
				return false;
			}

			let signMsgTrueLength = signMsg.readUInt16BE() + 42; // 2 (sign len) + 40 (random solt)

			return signMsg.subarray(0, signMsgTrueLength);
		}

		static sign(data: Buffer, privateKey: Buffer): false | Buffer {
			let signMsg = signSysIF.sign(data, privateKey);
			if(!signMsg) {
				return false;
			}

			let signMsgTrueLength = signMsg.readUInt16BE() + 42; // 2 (sign len) + 40 (random solt)

			return signMsg.subarray(0, signMsgTrueLength);
		}

		static verify(signature: Buffer, data: Buffer, publicKey: Buffer): boolean {
			let signMsgTrueLength = signature.readUInt16BE() + 42; // 2 (sign len) + 40 (random solt)
			if (signMsgTrueLength !== signature.length) {
				return false;
			}
			return signSysIF.verify(signature, data, publicKey);
		}
	};
}

const Falcon512Addon = require("./build/Release/Falcon512NistRound3.node");
const signSysIF512 = getFalconIF();
signSysIF512.signCore = Falcon512Addon;
signSysIF512.signSysName = 'Nist_round3_Falcon512';

class Falcon512 extends signSysIF512 {
	constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
		super(privateKey, publicKey, encryptedPrivateKeyMode);
	}
}


const Falcon1024Addon = require("./build/Release/Falcon1024NistRound3.node");
const signSysIF1024 = getFalconIF();
signSysIF1024.signCore = Falcon1024Addon;
signSysIF1024.signSysName = 'Nist_round3_Falcon1024';

class Falcon1024 extends signSysIF1024 {
	constructor(privateKey: Buffer, publicKey: Buffer, encryptedPrivateKeyMode: boolean = false) {
		super(privateKey, publicKey, encryptedPrivateKeyMode);
	}
}

export { Falcon512, Falcon1024 };



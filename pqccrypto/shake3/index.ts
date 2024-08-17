const napiSha3 = require("./build/Release/napiSha3.node");
type BinaryToTextEncoding = 'base64' | 'hex';

function napiShake256(data: Buffer): Buffer;
function napiShake256(data: Buffer, outputLength: number): Buffer;
function napiShake256(data: Buffer, outputLength: number, dig: BinaryToTextEncoding): string;
function napiShake256(data, outputLength = 32, dig?): Buffer | string {
	let temp = Buffer.alloc(outputLength);
	napiSha3.shake256(data, data.length, temp, outputLength);
	return (dig) ? temp.toString(dig) : temp;
}

export { napiShake256 }
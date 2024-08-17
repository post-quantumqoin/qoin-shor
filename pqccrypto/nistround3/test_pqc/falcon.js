const Falcon512NR3 = require("../falcon/build/Release/Falcon512NistRound3.node");
const Falcon1024NR3 = require("../falcon/build/Release/Falcon1024NistRound3.node");
const count = 1000;
var err = 0;
var suc = 0;

function testF512() {
    let pk = Buffer.alloc(Falcon512NR3.getPkLength());
    let sk = Buffer.alloc(Falcon512NR3.getSkLength());
    let seed = Buffer.alloc(Falcon512NR3.getSeedLength());
    let sm = Buffer.alloc(Falcon512NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextFalconSign");
    let messageErr = Buffer.from("AEGISMessageErrrFalconSign");
    Falcon512NR3.genkey(pk, sk, seed);
    // console.log(pk);
    // console.log(sk);
    // console.log(seed);
    Falcon512NR3.sign(sm, message, sk);
    let r = Falcon512NR3.verifySign(sm, message, pk);
    let r_err = Falcon512NR3.verifySign(sm, messageErr, pk);
    if (r === false || r_err === true) err++;
    else suc++;
    // console.log(sm);
    // console.log(r);
}

function testF1024() {
    let pk = Buffer.alloc(Falcon1024NR3.getPkLength());
    let sk = Buffer.alloc(Falcon1024NR3.getSkLength());
    let seed = Buffer.alloc(Falcon1024NR3.getSeedLength());
    let sm = Buffer.alloc(Falcon1024NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextFalconSign");
    let messageErr = Buffer.from("AEGISMessageErrrFalconSign");
    Falcon1024NR3.genkey(pk, sk, seed);
    // console.log(pk);
    // console.log(sk);
    // console.log(seed);
    Falcon1024NR3.sign(sm, message, sk);
    let r = Falcon1024NR3.verifySign(sm, message, pk);
    let r_err = Falcon1024NR3.verifySign(sm, messageErr, pk);
    if (r === false || r_err === true) err++;
    else suc++;
    // console.log(sm);
    // console.log(r);
}

function testF512GenBySeed() {
    console.log("test falcon 512 gen by seed");
    let pk = Buffer.alloc(Falcon512NR3.getPkLength());
    let pk2 = Buffer.alloc(Falcon512NR3.getPkLength());
    let sk = Buffer.alloc(Falcon512NR3.getSkLength());
    let sk2 = Buffer.alloc(Falcon512NR3.getSkLength());
    let seed = Buffer.alloc(Falcon512NR3.getSeedLength());
    Falcon512NR3.genkey(pk, sk, seed);
    Falcon512NR3.genPkBySeed(pk2, seed);
    Falcon512NR3.genSkBySeed(sk2, seed);
    console.log('verify pk, ', pk.equals(pk2));
    console.log('verify sk, ', sk.equals(sk2));
    let sm = Buffer.alloc(Falcon512NR3.getSignLength());
    let sm2 = Buffer.alloc(Falcon512NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextFalconSign");
    Falcon512NR3.sign(sm, message, sk);
    Falcon512NR3.signBySeed(sm2, message, seed);
    console.log('verify message, ', Falcon512NR3.verifySign(sm2, message, pk));
    console.log('verify message, ', Falcon512NR3.verifySign(sm, message, pk2));
}
function testF1024GenBySeed() {
    console.log("test falcon 1024 gen by seed");
    let pk = Buffer.alloc(Falcon1024NR3.getPkLength());
    let pk2 = Buffer.alloc(Falcon1024NR3.getPkLength());
    let sk = Buffer.alloc(Falcon1024NR3.getSkLength());
    let sk2 = Buffer.alloc(Falcon1024NR3.getSkLength());
    let seed = Buffer.alloc(Falcon1024NR3.getSeedLength());
    Falcon1024NR3.genkey(pk, sk, seed);
    Falcon1024NR3.genPkBySeed(pk2, seed);
    Falcon1024NR3.genSkBySeed(sk2, seed);
    console.log('verify pk, ', pk.equals(pk2));
    console.log('verify sk, ', sk.equals(sk2));
    let sm = Buffer.alloc(Falcon1024NR3.getSignLength());
    let sm2 = Buffer.alloc(Falcon1024NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextFalconSign");
    Falcon1024NR3.sign(sm, message, sk);
    Falcon1024NR3.signBySeed(sm2, message, seed);
    console.log('verify message, ', Falcon1024NR3.verifySign(sm2, message, pk));
    console.log('verify message, ', Falcon1024NR3.verifySign(sm, message, pk2));
}

console.log('-------------test falcon 512-------------');
for(let i=0;i < count;i++){
    testF512();
}
console.log('-------------test falcon 1024-------------');
for (let i = 0; i < count; i++) {
    testF1024();
}
testF512GenBySeed();
testF1024GenBySeed();
console.log('end test falcon');
console.log('err:' + err);
console.log('suc:' + suc);
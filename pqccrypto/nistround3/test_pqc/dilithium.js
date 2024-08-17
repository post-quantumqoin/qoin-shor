const Dilithium2NR3 = require("../dilithium/build/Release/Dilithium2RefNistRound3.node");
const Dilithium3NR3 = require("../dilithium/build/Release/Dilithium3RefNistRound3.node");
const Dilithium5NR3 = require("../dilithium/build/Release/Dilithium5RefNistRound3.node");
const count = 10000;
var err = 0;
var suc = 0;

function testD2() {
    let pk = Buffer.alloc(Dilithium2NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium2NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium2NR3.getSeedLength());
    let sm = Buffer.alloc(Dilithium2NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    let messageErr = Buffer.from("AEGISMessageErrrDilithiumSign");
    Dilithium2NR3.genkey(pk, sk, seed);
    // console.log(pk);
    // console.log(sk);
    // console.log(seed);
    Dilithium2NR3.sign(sm, message, sk);
    let r = Dilithium2NR3.verifySign(sm, message, pk);
    let r_err = Dilithium2NR3.verifySign(sm, messageErr, pk);
    if (r === false || r_err === true) err++;
    else suc++;
    // console.log(sm);
    // console.log(r);
}

function testD3() {
    let pk = Buffer.alloc(Dilithium3NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium3NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium3NR3.getSeedLength());
    let sm = Buffer.alloc(Dilithium3NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    let messageErr = Buffer.from("AEGISMessageErrrDilithiumSign");
    Dilithium3NR3.genkey(pk, sk, seed);
    // console.log(pk);
    // console.log(sk);
    // console.log(seed);
    Dilithium3NR3.sign(sm, message, sk);
    let r = Dilithium3NR3.verifySign(sm, message, pk);
    let r_err = Dilithium3NR3.verifySign(sm, messageErr, pk);
    if (r === false || r_err === true) err++;
    else suc++;
    // console.log(sm);
    // console.log(r);
}

function testD5() {
    let pk = Buffer.alloc(Dilithium5NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium5NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium5NR3.getSeedLength());
    let sm = Buffer.alloc(Dilithium5NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    let messageErr = Buffer.from("AEGISMessageErrrDilithiumSign");
    Dilithium5NR3.genkey(pk, sk, seed);
    // console.log(pk);
    // console.log(sk);
    // console.log(seed);
    Dilithium5NR3.sign(sm, message, sk);
    let r = Dilithium5NR3.verifySign(sm, message, pk);
    let r_err = Dilithium5NR3.verifySign(sm, messageErr, pk);
    if (r === false || r_err === true) err++;
    else suc++;
    // console.log(sm);
    // console.log(r);
}

function testD2GenBySeed() {
    console.log("test dilithium2 gen by seed");
    let pk = Buffer.alloc(Dilithium2NR3.getPkLength());
    let pk2 = Buffer.alloc(Dilithium2NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium2NR3.getSkLength());
    let sk2 = Buffer.alloc(Dilithium2NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium2NR3.getSeedLength());
    Dilithium2NR3.genkey(pk, sk, seed);
    Dilithium2NR3.genPkBySeed(pk2, seed);
    Dilithium2NR3.genSkBySeed(sk2, seed);
    console.log('verify pk, ', pk.equals(pk2));
    console.log('verify sk, ', sk.equals(sk2));
    let sm = Buffer.alloc(Dilithium2NR3.getSignLength());
    let sm2 = Buffer.alloc(Dilithium2NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    Dilithium2NR3.sign(sm, message, sk);
    Dilithium2NR3.signBySeed(sm2, message, seed);
    console.log('verify message, ', Dilithium2NR3.verifySign(sm2, message, pk));
    console.log('verify message, ', Dilithium2NR3.verifySign(sm, message, pk2));
}
function testD3GenBySeed() {
    console.log("test dilithium3 gen by seed");
    let pk = Buffer.alloc(Dilithium3NR3.getPkLength());
    let pk2 = Buffer.alloc(Dilithium3NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium3NR3.getSkLength());
    let sk2 = Buffer.alloc(Dilithium3NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium3NR3.getSeedLength());
    Dilithium3NR3.genkey(pk, sk, seed);
    Dilithium3NR3.genPkBySeed(pk2, seed);
    Dilithium3NR3.genSkBySeed(sk2, seed);
    console.log('verify pk, ', pk.equals(pk2));
    console.log('verify sk, ', sk.equals(sk2));
    let sm = Buffer.alloc(Dilithium3NR3.getSignLength());
    let sm2 = Buffer.alloc(Dilithium3NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    Dilithium3NR3.sign(sm, message, sk);
    Dilithium3NR3.signBySeed(sm2, message, seed);
    console.log('verify message, ', Dilithium3NR3.verifySign(sm2, message, pk));
    console.log('verify message, ', Dilithium3NR3.verifySign(sm, message, pk2));
}

function testD5GenBySeed() {
    console.log("test dilithium5 gen by seed");
    let pk = Buffer.alloc(Dilithium5NR3.getPkLength());
    let pk2 = Buffer.alloc(Dilithium5NR3.getPkLength());
    let sk = Buffer.alloc(Dilithium5NR3.getSkLength());
    let sk2 = Buffer.alloc(Dilithium5NR3.getSkLength());
    let seed = Buffer.alloc(Dilithium5NR3.getSeedLength());
    Dilithium5NR3.genkey(pk, sk, seed);
    Dilithium5NR3.genPkBySeed(pk2, seed);
    Dilithium5NR3.genSkBySeed(sk2, seed);
    console.log('verify pk, ', pk.equals(pk2));
    console.log('verify sk, ', sk.equals(sk2));
    let sm = Buffer.alloc(Dilithium5NR3.getSignLength());
    let sm2 = Buffer.alloc(Dilithium5NR3.getSignLength());
    let message = Buffer.from("AEGISMessageTextDilithiumSign");
    Dilithium5NR3.sign(sm, message, sk);
    Dilithium5NR3.signBySeed(sm2, message, seed);
    console.log('verify message, ', Dilithium5NR3.verifySign(sm2, message, pk));
    console.log('verify message, ', Dilithium5NR3.verifySign(sm, message, pk2));
}

console.log('-------------test dilithium2-------------');
for(let i=0;i < count;i++){
    testD2();
}
console.log('-------------test dilithium3-------------');
for (let i = 0; i < count; i++) {
    testD3();
}
console.log('-------------test dilithium5-------------');
for (let i = 0; i < count; i++) {
    testD5();
}
testD2GenBySeed();
testD3GenBySeed();
testD5GenBySeed();
console.log('end test dilithium');
console.log('err:' + err);
console.log('suc:' + suc);
import {sha256} from "js-sha256";
// @ts-ignore
import hexer from "browser-string-hexer";
import {Buffer} from 'buffer';

const Web3 = require('web3');

function hash256(message: string): string {
    return sha256(message)
}

async function sign(message: string): Promise<string> {
    if (typeof (window as any).ethereum === 'undefined') {
        console.log("no eth found")
        return "";
    }
    const ethereum = (window as any).ethereum;

    const accounts = await ethereum.request({method: 'eth_requestAccounts'})
    let msgParam = hexer(message);
    let param = [msgParam, accounts[0]];

    return ethereum.request({
        method: 'personal_sign',
        params: param,
        id: 0
    })
}

// base64(publicKey, " ", message)
async function encodePubKeyAndSig(message: string) {
    if (typeof (window as any).ethereum === 'undefined') {
        console.log("no eth found")
        return "";
    }
    const ethereum = (window as any).ethereum;
    const web3 = new Web3(ethereum);

    let accounts = await web3.eth.getAccounts();

    await ethereum.request({method: 'eth_requestAccounts'})
    let msgParam = hexer(message);

    let param = [msgParam, accounts[0]];

    const sig: string = await ethereum.request({
        method: 'personal_sign',
        params: param,
        id: 0
    })

    console.log(accounts)
    return Buffer.from(`${accounts[0]} ${sig}` as string).toString('base64');
}

export {hash256, sign, encodePubKeyAndSig};
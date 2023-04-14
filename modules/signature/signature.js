require('dotenv').config();
require('../components/enc-base64-min.js')();
const CryptoJS = require('../rollups/hmac-sha256.js');

let accessKey = process.env['NCR_ACCESS_KEY'];

let makeSignature = (method, path, timestamp)=>{
	let secretKey = process.env['NCR_SECRET_KEY'];
	let hmac = CryptoJS.algo.HMAC.create(CryptoJS.algo.SHA256, secretKey);
  let sig = `${method} ${path}\n${timestamp}\n${accessKey}`;
  hmac.update(sig);
	let hash = hmac.finalize();
	return hash.toString(CryptoJS.enc.Base64);
}

let getHeader = (method,path)=>{
  let timestamp = String(new Date().valueOf());
  let sig = makeSignature(method,path,timestamp);
  return {
    'x-ncp-apigw-timestamp':timestamp,
    'x-ncp-iam-access-key':accessKey,
    'x-ncp-apigw-signature-v2':sig
  }
}

module.exports = getHeader;
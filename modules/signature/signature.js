require('dotenv').config();
require('../components/enc-base64-min.js');
const CryptoJS = require('../rollups/hmac-sha256.js');

let makeSignature = (method, url)=>{
	let space = ' ';
	let newLine = '\n';
	// let method = "GET";
	// let url = "/photos/puppy.jpg?query1=&query2";
	let timestamp = '1681439143336';			// current timestamp (epoch)
	let accessKey = process.env['NCR_ACCESS_KEY'];
	let secretKey = process.env['NCR_SECRET_KEY'];

	let hmac = CryptoJS.algo.HMAC.create(CryptoJS.algo.SHA256, secretKey);
	hmac.update(method);
	hmac.update(space);
	hmac.update(url);
	hmac.update(newLine);
	hmac.update(timestamp);
	hmac.update(newLine);
	hmac.update(accessKey);

	let hash = hmac.finalize();

	return hash.toString(CryptoJS.enc.Base64);
}
module.exports = makeSignature;
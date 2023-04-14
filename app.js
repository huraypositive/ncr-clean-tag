require('dotenv').config();
const apigw = process.env['NCR_API_GATEWAY_URL'];
const fetch = require('request-promise-native');
const getHeader = require('./modules/signature/signature.js');
const timeout = process.env['TIMEOUT'] || 5000;

console.log(`apigw: ${apigw}`)
let method = 'GET';
let path = `/ncr/api/v2/repositories`;
let headers = getHeader(method,path);
headers['Content-Type'] = 'application/json; charset=utf-8';
console.log(headers);

fetch({
  uri: `${apigw}${path}`,
  method,
  headers,
  timeout
})
.then(res=>{
  console.log(res);
  console.log("arstasrteiranstasjrtie")
})
.catch(err=>{
  // console.error(err);
})
require('dotenv').config();
const apigw = process.env['NCR_API_GATEWAY_URL'] || 'https://ncr.apigw.ntruss.com';
const fetch = require('request-promise-native');
const getHeader = require('./modules/signature/signature.js');
const timeout = process.env['TIMEOUT'] || 5000;

let method = 'GET';
let path = '/ncr/api/v2/repositories';
let headers = getHeader(method,path);
headers['Content-Type'] = 'application/json; charset=utf-8';

fetch({
  uri: `${apigw}${path}`,
  method,
  headers,
  timeout
})
.then(res=>{
  console.log(JSON.parse(res));
})
.catch(err=>{
  console.error(err);
})
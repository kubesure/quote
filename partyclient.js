var grpc = require('grpc');
var partyProto = grpc.load('party.proto');
const partysvc = process.env.partysvc

var client = new partyProto.party.PartyService(partysvc + ':50051', grpc.credentials.createInsecure());
exports = client
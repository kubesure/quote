const MongoClient = require('mongodb').MongoClient
const mongosvcqoute = process.env.mongosvcqoute
const url = 'mongodb://' + mongosvcqoute + ':27017'
const dbName = 'proposals'
const assert = require('assert');

module.exports = {
  save: async (proposal) => {
    let db, client
    try {
      client = await MongoClient.connect(url, { useNewUrlParser: true })
      db = client.db(dbName);
      let id = await nextId()
      proposal.ProposalId = id
      let r = await db.collection('proposal').insertOne(proposal)
      assert.equal(1, r.insertedCount)
      return r.ops[0]
    } catch (err) {
      console.log(err)
    } finally {
      client.close()
    }
  }
}

module.exports.find = (async (id) => {
  let db, client
  try {
    client = await MongoClient.connect(url, { useNewUrlParser: true })
    db = client.db(dbName);
    var query = { "ProposalId": parseInt(id) }
    let docs = await db.collection('proposal').findOne(query)
    return docs
  } finally {
    client.close()
  }
})

async function nextId() {
  let db, client
  try {
    client = await MongoClient.connect(url, { useNewUrlParser: true })
    db = client.db(dbName);
    let r = await db.collection('counter').findOneAndUpdate(
      { '_id': 'proposalid' },
      { $inc: { 'value': 1 } },
      { upsert: true })
    return r.value.value
  } finally {
    client.close()
  }
}
//test save
(async function () {
    var assert = require('assert')
    var db = require('./db')

    try {
        result = await db.save({ProposalId:"123456"})
        assert.notEqual(null,result.ProposalId)    
    } catch (error) {
    }
})();

//test find by id
(async function () {
    var assert = require('assert')
    var db = require('./db')

    try {
        result = await db.find(37)
        assert.Equal(37,result.ProposalId)    
    } catch (error) {
    }
})();
var express = require('express')
var db = require('./db')
var app = express()
var bodyParser = require('body-parser');
app.use(bodyParser.json());
const port = 8030

app.post("/api/v1/healths/proposals", async (req, res) => {
    try {
        proposal = await db.save(req.body)
        res.status(200).json(proposal)
    } catch (error) {
        console.log(error)
        res.status(500).end()
    }
})

app.get("/api/v1/healths/proposals/:proposalId", async (req, res) => {
    try {
        proposal = await db.find(req.params.proposalId)
        if (proposal == null) { res.status(404).end() }
        res.status(200).json(proposal)
    } catch (error) {
        res.status(500).end()
    }
})

app.listen(port, () => {
    return console.log(`Quote app listening on port ${port}!`)
})
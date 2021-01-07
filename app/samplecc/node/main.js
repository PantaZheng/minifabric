/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const {Gateway, Wallets} = require('fabric-network');
const fs = require('fs');
const path = require('path');
const YAML = require('yaml');
const log4js = require('log4js');
const uuid = require('uuid');
const seedrandom = require('seedrandom');
const WaitGroup = require("sync-wait-group").WaitGroup;

log4js.configure({
    appenders: {node: {type: "file", filename: "node.log"}},
    categories: {default: {appenders: ["node"], level: "all"}}
});

const logger = log4js.getLogger('node');
logger.level = "all"

async function main() {
    try {
        // load the network configuration
        const ccp = YAML.parse(fs.readFileSync('./profiles/connection.yaml', 'utf8'))
        // Create a new file system based wallet for managing identities.
        const wallet = await Wallets.newFileSystemWallet('./profiles/wallets');
        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get('Admin');
        if (!identity) {
            logger.error('Admin identity can not be found in the wallet');
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, {wallet, identity: 'Admin', discovery: {enabled: true, asLocalhost: false}});
        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');
        // Get the contract from the network.
        const contract = network.getContract('samplecc');
        // Submit the specified transaction.

        const nowTime = Date.now()
        logger.info(`Time is ${nowTime.toString()}`)
        const seededRand = seedrandom(nowTime)

        let wg = new WaitGroup()
        let counter = []
        for (let i = 0; i < 10; i++) {
            wg.add(1)
            const uid = uuid.v4().toString() + i.toString();
            const sid = seededRand.int32();
            counter.push(
                contract.submitTransaction('invoke', 'put',
                    uid, sid).then((response)=>{
                    return {uid: uid,sid:sid, message: response.toJSON().data}
                })
            )
            wg.done()
        }

        let results = await Promise.all(counter)
        logger.info(results)
        await wg.wait()
        logger.warn(`promise.all time took is ${Date.now()-nowTime}`)

        // for (let i = 0; i < 10; i++) {
        //     const uid = uuid.v4().toString() + i.toString();
        //     const sid = seededRand.int32();
        //     let result = await contract.submitTransaction('invoke', 'put',
        //         uid, sid)
        //     logger.info({uid: uid, sid: sid, message: result.toJSON().data})
        // }
        // logger.warn(`loop async time took is ${Date.now() - nowTime}`)
        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        logger.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main().then();

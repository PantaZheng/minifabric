/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

const {WorkloadModuleBase} = require('@hyperledger/caliper-core');

/**
 * Workload module for initializing the SUT with various marbles.
 */
class AddPrivateRecordWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = -1;
        this.startTime = Date.now()
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
    }

    /**
     * Assemble TXs for creating new marbles.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        this.txIndex++;
        let args = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'addPvtRecord',
            contractArguments: [
                'Pvt'+(this.txIndex).toString(),
                Math.floor(Math.random()*10)
            ],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(args);
    }

    async cleanupWorkloadModule() {
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new AddPrivateRecordWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;

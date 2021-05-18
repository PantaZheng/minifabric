/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// Running TestApp:
// gradle runApp

package application.java;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Random;
import java.util.UUID;
import java.util.concurrent.TimeoutException;
import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.GatewayException;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.gateway.Wallets;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App {

  public static final String NETWORK_NAME = "mychannel";
  public static final String CONTRACT_ID = "samplecc";

  public static void main(String[] args) throws Exception {

    Logger logger = LoggerFactory.getLogger(App.class);
    // Load a file system based wallet for managing identities.
    Path walletPath = Paths.get("profiles","wallets");
    Wallet wallet = Wallets.newFileSystemWallet(walletPath);
    // load a CCP
    Path networkConfigPath = Paths.get("profiles", "connection.yaml");

    Gateway.Builder builder = Gateway.createBuilder().
        identity(wallet, "Admin").networkConfig(networkConfigPath);

    try (Gateway gateway = builder.connect()) {

      Network network = gateway.getNetwork(NETWORK_NAME);
      Contract contract = network.getContract(CONTRACT_ID);

      Random random = new Random();

      long startTime = System.currentTimeMillis();

      for (int i = 0; i < 10; i++) {
        UUID uuid = UUID.randomUUID();
        byte[] results =
            contract.submitTransaction("Put", uuid.toString(), Integer.toString(random.nextInt()));
        logger.error("The results is " + Arrays.toString(results));
        System.out.println("The results is " + Arrays.toString(results));
      }
      logger.error("The time took is " + (System.currentTimeMillis() - startTime));
    }catch (GatewayException | TimeoutException | InterruptedException e) {
      e.printStackTrace();
    }
  }
}

#!/usr/bin/env bash
# reference https://zhuanlan.zhihu.com/p/36119670
# install make tool

#make cryptogen
#
#make configtxgen
#
#make orderer
#
#make peer

rm -rf testnet

mkdir -p testnet/config

cp orderer.yaml ./testnet/config

cp core.yaml ./testnet/config

export PATH=`pwd`/.build/bin:$PATH

cd testnet/config

# generate crypto-config.yaml

echo "OrdererOrgs:
  - Name: Orderer
    Domain: example.com
    Specs:
      - Hostname: orderer
PeerOrgs:
  - Name: Org
    Domain: example.com
    Template:
      Count: 1
      Hostname: peer
    Users:
      Count: 1" > crypto-config.yaml

cryptogen generate --config=./crypto-config.yaml --output=./crypto-config

# configtx.yaml

echo "Organizations:
    - &OrdererOrg
        Name: OrdererOrg
        ID: OrdererMSP
        MSPDir: crypto-config/ordererOrganizations/example.com/msp

    - &Org
        Name: OrgMSP
        ID: OrgMSP
        MSPDir: crypto-config/peerOrganizations/example.com/msp
        AnchorPeers:
            - Host: peer.example.com
              Port: 7051

Orderer: &OrdererDefaults
    OrdererType: solo
    Addresses:
        - orderer.example.com:7050
    BatchTimeout: 2s
    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB
    Organizations:

Application: &ApplicationDefaults
    Organizations:

Profiles:
    SingleSoloOrdererGenesis:
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
        Consortiums:
            SampleConsortium:
                Organizations:
                    - *Org
    SingleSoloChannel:
        Consortium: SampleConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
                - *Org" > configtx.yaml

mkdir channel-artifacts

export FABRIC_CFG_PATH=`pwd`

configtxgen -outputBlock ./channel-artifacts/genesis.block -profile SingleSoloOrdererGenesis

export CHANNEL_NAME=sxlchannel
configtxgen -outputCreateChannelTx ./channel-artifacts/channel.tx -profile SingleSoloChannel -channelID $CHANNEL_NAME

# orderer environment
mkdir data
export rootDir=`pwd`
export PATH=$rootDir/bin:$PATH
export ORDERER_GENERAL_LOGLEVEL=DEBUG
export ORDERER_GENERAL_TLS_ENABLED=false
export ORDERER_GENERAL_PROFILE_ENABLED=false
export ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
export ORDERER_GENERAL_LISTENPORT=7050
export ORDERER_GENERAL_GENESISMETHOD=file
export ORDERER_GENERAL_GENESISFILE=$rootDir/channel-artifacts/genesis.block
export ORDERER_GENERAL_LOCALMSPDIR=$rootDir/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp
export ORDERER_GENERAL_LOCALMSPID=OrdererMSP

export ORDERER_FILELEDGER_LOCATION=$rootDir/data/orderer

kill -9 `ps -ef | grep orderer | awk '{print $2}'`

#orderer >/tmp/orderer.log 2>&1 &
orderer

# if CORE_CHAINCODE_MODE=dev the chaincode can not publish
export rootDir=`pwd`
export CORE_PEER_ID=example_org
export CORE_CHAINCODE_MODE=net
export CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
export CORE_PEER_NETWORKID=net
export CORE_LOGGING_LEVEL=INFO
export CORE_PEER_TLS_ENABLED=false
export CORE_PEER_PROFILE_ENABLED=false
export CORE_PEER_ADDRESS=0.0.0.0:7051
export CORE_PEER_LISTENADDRESS=0.0.0.0:7051
export CORE_PEER_GOSSIP_ENDPOINT=0.0.0.0:7051
export CORE_PEER_EVENTS_ADDRESS=0.0.0.0:7053
export CORE_PEER_LOCALMSPID=OrgMSP
export CORE_LEDGER_STATE_STATEDATABASE=goleveldb
export CORE_PEER_MSPCONFIGPATH=$rootDir/crypto-config/peerOrganizations/example.com/peers/peer.example.com/msp
export CORE_PEER_FILESYSTEMPATH=$rootDir/data/peer

kill -9 `ps -ef | grep peer | awk '{print $2}'`

peer node start


# channel
export rootDir=`pwd`
export PATH=`pwd`/.build/bin:$PATH
export CHANNEL_NAME=sxlchannel
export CORE_CHAINCODE_MODE=net
export CORE_PEER_ID=peer-cli
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=OrgMSP
export CORE_PEER_MSPCONFIGPATH=$rootDir/crypto-config/peerOrganizations/example.com/users/Admin@example.com/msp

peer channel create -o 127.0.0.1:7050 -c $CHANNEL_NAME -f $rootDir/channel-artifacts/channel.tx

peer channel join -b $CHANNEL_NAME.block


# chaincode

export CHANNEL_NAME=sxlchannel
export CORE_CHAINCODE_MODE=net
export CORE_PEER_ID=peer-cli
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=OrgMSP
export CORE_PEER_MSPCONFIGPATH=$rootDir/crypto-config/peerOrganizations/example.com/users/Admin@example.com/msp

# install

peer chaincode install -n ctk -v 1.0 -p ctkcontract

# instantiate
peer chaincode instantiate -o 127.0.0.1:7050 -C $CHANNEL_NAME -n ctk -v 1.0 -c "{\"Args\":[\"init\"]}" -P "OR('OrgMSP.member')"
peer chaincode upgrade -o 127.0.0.1:7050 -C $CHANNEL_NAME -n ctk -v 1.0 -c "{\"Args\":[\"init\"]}" -P "OR('OrgMSP.member')"

# invoke query

peer chaincode invoke -o 127.0.0.1:7050 -C $CHANNEL_NAME -n ctk -c  "{\"Args\":[\"invoke\",\"info\",\"ctk\"]}"








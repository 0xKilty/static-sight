---
title: Patricoin
date: 12/4/2022
tags:
  - Solidity
---
## Motivation
Patricoin is an [ERC-777 token](https://eips.ethereum.org/EIPS/eip-777) that I made to further understand crypto currency and smart contracts.
I named the project/token, "Patricoin" because one of my friends name is Patrick and to poke at the fact that almost all crypto currencies are named after something with the word `coin` mixed in somewhere.
What sparked my interest in learning more about crypto currencies was a presentation by [Waylon Jepsen](https://github.com/0xJepsen) on crypto [arbitrage](https://en.wikipedia.org/wiki/Arbitrage).
## Creation
### ERC-777 Standard
The ordinary standard tokens are upheld to are ERC-20 which include a list of [functions](https://ethereum.org/en/developers/docs/standards/tokens/erc-20/#methods) that the token must implement. An ERC-777 token is backwards compatible with an ERC-20 token but adds a few new features such as [hooks](https://ethereum.org/en/developers/docs/standards/tokens/erc-777/#hooks). 
### Tools
To actually develop the token, I used an ethereum development tool chain called [hardhat](https://hardhat.org/). Hardhat is an amazing tool for the creation of the smart contract and the deployment. 
To protect against token vulnerabilities, such as a [reentrancy attack](https://hackernoon.com/hack-solidity-reentrancy-attack), and to speed up the development process, I used [OpenZepplin](https://www.openzeppelin.com/) which is a library for developing smart contracts in solidity. This can protect against re-entrancy attacks because of the wide use of OpenZepplin and the thousands of contributers to the project. 
## Deployment
### Goelri Testnet
Using hardhat, I deployed the token on the [Goelri testnet](https://goerli.net/) which is a testnet on the ethereum blockchain. A testnet is network separate from the real network of a crypto currency (in this case ethereum), so that developers of these smart contracts can test their deployment. I choose to not deploy Patricoin on the real network because it was made just for fun and it costs money to deploy it. Since the blockchain is visible to anyone, here is the [final deployment](https://goerli.etherscan.io/token/0xb5e499390d8cef2a6a76158963205d8d18e71df7#code) on the blockchain. 
### Uniswap Pool
[Uniswap](https://uniswap.org/) is a DEX or a Decentralized Exchange meaning you can buy and sell crypto currency directly from other people. To build a liquidity pool for patricoin, I used uniswap to put in some Goelri Testnet ETH into a liquidity pool for patricoin so that it could be worth some fake money and I could derive a price. To view this liquidity pool, make sure you have Goelri Testnet selected on the uniswap website, [Pool](https://app.uniswap.org/#/pool/45776).
## Creating the website
I thought it would be a fun idea to sharpen my react skills and learn about web3 to build a [website](https://0xkilty.github.io/patricoin/) surrounding patricoin. I used [ethers.js](https://docs.ethers.io/v5/|ethers.js) and [uniswap v3 sdk](https://docs.uniswap.org/sdk/v3/overview) for the development and connection to the blockchain along with Metamask compatibility.
## Conclusion
I enjoyed the process of creating patricoin and the [website](https://0xkilty.github.io/patricoin/) to go along with it. I also learned about the fine details of blockchain and about web3. Even though there wasn't much room for creativity within the token (because my economic knowledge is limited), I still felt that I learned a lot and about how other tokens on the blockchain function.
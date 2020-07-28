# upgradePreCheck

This is a pre-check program called by autoUpgrade and upgrade.

## How to complile

### 1. Build with docker container

go to build folder and run buildCLL.sh to build with docker container

`bash build/buildCLI.sh`

### 2. Build with make file

***Notice***

*This makefile is just make it simple to build upgradePreCheck for test. Since golang version depends on your env, it is recommended to build with docker container for production use.*

`cd upgradePreCheck ; make`

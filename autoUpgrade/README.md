# autoUpgrade

Automatically upgrade CDF. 

Support on both windows and linux.

## How to complile

### 1. Build with docker container

go to build folder and run buildCLL.sh to build with docker container

`bash build/buildCLI.sh`

### 2. Build with make file

***Notice***

*This makefile is just make it simple to build upgradePreCheck for test. Since golang version depends on your env, it is recommended to build with docker container for production use.*

Build autoUpgrade binary
`make all`

Check current env
`make check`

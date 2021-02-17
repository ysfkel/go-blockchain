
#Data architecture

Two groups of data


1. Block data: which is stored with metadata which describes all of the blocks on the chain
2. Chain state: state of the chain and all of the current unspent transactions ouputs 

- with bitcoin each block has its own separate file on the disk to increase performace by splitting each block into its own 
  file so that only that blockfile is opened

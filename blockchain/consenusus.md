# implements proof of work

Proof-of-work, POW — is the way to protect the system from DoS attacks and some other types of attacks which is founded on solving some tasks by client with defined difficulty (POW-task)where the solution could be checked easily by server

 1. Get the data from the block
 2. Create a counter (nonce) which starts at 0
 3. Create a hash of the data plus the counter
 4. Check the resulting hash to see if it meets a set of requirements 
    (this is the idae of difficulty)
 5. if the hash meets the set of requirements then we use that hash and say it signs the block
 6. otherwise we go back and create another hash and we repeat the process until we get a 
    hash that meets the set of requirements

	Requirements:
	1. First few bytes of the hash must contain 0s
	   - In the original bitcoin proof of work specification (hash cash) 
	     the original difficulty was to get 20 consecutive bit of the hash as 0s 
		 this requirement gets adjusted over time and that is essentially the difficulty.
		 the difficult goes up means we must have more preceeding zeros infront of the hash
		 to be valid
	   - In this implementation , the difficulty is a constant , however in a real blockchain
	     you would have an algorithm that slowly imcrements this difficulty over a large period of time 
		 The main reason to do this is to account for the increasing number of miners on the network and also
		 account for the increasing computing power of computers in general inorder to make the time to mine a block stay thesame
		 and also to make the block rate stay thesame. This means you would need to have a certain amount of computational
		 power to produce blocks at that rate , but also keep the time to sign a block down

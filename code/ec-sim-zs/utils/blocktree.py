import json

# Parse a block tree from json
class BlockTree:
    def __init__(self, filename):
        self.readAndParse(filename)

    def readAndParse(self, filename):
        # Setup output datastructures.
        self.BlocksByNonce = dict()
        self.BlocksByHeight = dict()
        self.Miners = []
        
        with open(filename, 'r') as jfile:
            data = json.load(jfile)
            
            # Get all blocks.  Add all into indexes.
            blocks = data["blocks"]
            for block in blocks:
                self.BlocksByNonce[block["nonce"]] = block
                h = block["height"]
                if h not in self.BlocksByHeight:
                    self.BlocksByHeight[h] = []
                self.BlocksByHeight[h].append(block)
                
            # Get all miners
            self.Miners = data["miners"]

    ### Extract various metrics from the block tree.

    # HeaviestChain returns a set of nonces of blocks that form the heaviest
    # chain in the block tree.
    def HeaviestChain(self):
        length = len(self.BlocksByHeight)
        cur = self.heaviestAtHeight(length)
        curHeight = cur["height"]
        chain = set()
        chain.add(cur["nonce"])
        while curHeight > 0:
            next = parentNonces(cur["tipset"]["name"])
            for nonce in next:
                chain.add(nonce)
            cur = self.BlocksByNonce[next[0]] 

    # RatioUsefulBlocks returns the ratio of blocks making it into the 
    # heaviest chain to the total blocks mined
    def RatioUsefulBlocks(self):
        mainChain = self.HeaviestChain()
        return float(len(mainChain)) / float(len(self.BlocksByHeight))

    # AvgHeadsPerRound returns the mean number of possible mining heads per
    # round.
    def AvgHeadsPerRound(self):
        acc = 0.0
        numTrials = len(self.BlocksByHeight)
        for round in range(0, numTrials):
            if round in self.BlocksByHeight:
                acc += len(self.BlocksByHeight[round])
                
        return acc / float(numTrials) 

    # NumReorgs returns the number of times that the heaviest tipset in round
    # n was on a different fork from the heaviest tipset at round n - 1.
    def NumReorgs(self):
        raise "lalala"

    def heaviestAtHeight(self, h):
        if h not in self.BlocksByHeight:
            raise "no BlocksByHeight entry at height " + string(n)
        hBlock = None
        hWeight = -1
        for block in self.BlocksByHeight[h]:
            if block["weight"] > hWeight:
                hBlock = block
                hWeight = block["weight"]

        if hBlock is None:
            raise "bad BlocksByHeight entry with 0 blocks"

        return hBlock


def parentNonces(name):
    nonces = name.split('-')
    return [int(nonce) for nonce in nonces]
        

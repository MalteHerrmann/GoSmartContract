# Remarks

This file contains some remarks, that occured to me during the work on this task.

----

## Initialization of a new node
After cloning the latest commit (\#41a7b4b), the `init.sh` fails to run properly. 
It claims, that the genesis file cannot be found, even though it was created in 
the correct location at `$HOME/.evmosd/config/genesis.json`. This is, because only 
the relative path (`.config/genesis.json`) is searched, when executing `evmosd validate-genesis`.

Full error:

```
Error: couldn't read GenesisDoc file: open config/genesis.json: no such file or 
directory. Make sure that you have correctly migrated all Tendermint consensus 
params, please see the chain migration guide at 
https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html for more 
info
```

#TODO: Problem seems to be in usage of Tendermint?

----

## Starting the node

The problem above also leads to `evmosd start` only working from within `$HOME/.evmosd/config`,
which was not the case in the previous version (v3.0.0), that I had installed on my laptop.

----

## Updating the config and data storage directory
The Evmos docs (https://docs.evmos.org/validators/quickstart/binary.html#config-and-data-directory) 
describe a way to update the location, where configuration and data storage of the 
node resides. 

However, neither `evmosd --home [directory]` nor `evmosd config --home [directory]`
are recognized as valid commands. This should either be updated in the docs if the 
command is indeed futile or made clearer in case I misunderstood the documentation.

package coremandos

// DharitriProtectedKeyPrefix prefixes all Dharitri reserved storage. Only the protocol can write to keys starting with this.
const DharitriProtectedKeyPrefix = "DHARITRI"

// DharitriRewardKey is the storage key where the protocol writes when sending out rewards.
const DharitriRewardKey = DharitriProtectedKeyPrefix + "reward"

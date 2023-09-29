package arwenmandos

// Kalyan3104ProtectedKeyPrefix prefixes all Kalyan3104 reserved storage. Only the protocol can write to keys starting with this.
const Kalyan3104ProtectedKeyPrefix = "KALYAN3104"

// Kalyan3104RewardKey is the storage key where the protocol writes when sending out rewards.
const Kalyan3104RewardKey = Kalyan3104ProtectedKeyPrefix + "reward"

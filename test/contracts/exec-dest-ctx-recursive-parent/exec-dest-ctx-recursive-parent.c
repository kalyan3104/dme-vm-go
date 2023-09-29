#include "../kalyan3104/context.h"
#include "../kalyan3104/bigInt.h"
#include "../kalyan3104/types.h"
#include "../kalyan3104/test_utils.h"

u64 maxGasForCalls = 100000;
 
byte childAddress[32] = "childSC.........................";
byte bigIntCounterKey[] = "recursiveIterationBigCounter....";
bigInt bigIntCounterID = 88;
byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,5};
byte smallCounterKey[] = "recursiveIterationCounter.......";

byte arguments[1] = {0};
int argumentsLengths[1] = {1};

void parentCallsChild() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);

	storeIterationNumber(iteration, 'P');
	finishIterationNumber(iteration, 'P');

  incrementIterCounter(smallCounterKey);
	incrementBigIntCounter(bigIntCounterID);

  // Run next iteration.
	byte functionName[] = "childCallsParent";
	if (iteration > 0) {
		arguments[0] = iteration - 1;
    int result = executeOnDestContext(
        maxGasForCalls,
        childAddress,
        executeValue,
        functionName,
        16,
        1,
        (byte*)argumentsLengths,
        arguments
    );

    finishResult(result);
  } else {
    bigIntStorageStoreUnsigned(bigIntCounterKey, 32, bigIntCounterID);
  }
}

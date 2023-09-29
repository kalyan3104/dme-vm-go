#include "../kalyan3104/context.h"
#include "../kalyan3104/bigInt.h"
#include "../kalyan3104/test_utils.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte parentFinishA[] = "parentFinishA";
byte parentFinishB[] = "parentFinishB";

byte childKey[] =  "childKey........................";

byte parentTransferReceiver[] = "parentTransferReceiver..........";
byte parentTransferValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};
byte parentTransferData[] = "parentTransferData";

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};
u32 executeArgumentsLengths[] = {15, 16, 10};
byte executeArgumentsData[] = "First sentence.Second sentence.Some text.";

void parentFunctionPrepare() {
	storageStore(parentKeyA, 32, parentDataA, 11);
	storageStore(parentKeyB, 32, parentDataB, 11);
	finish(parentFinishA, 13);
	finish(parentFinishB, 13);
	int result = transferValue(
			parentTransferReceiver,
			parentTransferValue,
			parentTransferData,
			18
	);
	finishResult(result);
}

void parentFunctionWrongCall() {
	parentFunctionPrepare();
	byte childAddress[] = "wrongSC.........................";
	byte functionName[] = "childFunction";

	int result = executeOnDestContext(
			10000,
			childAddress,
			executeValue,
			functionName,
			13,
			3,
			(byte*)executeArgumentsLengths,
			executeArgumentsData
	);
	finishResult(result);
}

void parentFunctionChildCall() {
	parentFunctionPrepare();
	byte childAddress[] = "childSC.........................";
	byte functionName[] = "childFunction";
	int result = executeOnDestContext(
			200000,
			childAddress,
			executeValue,
			functionName,
			13,
			3,
			(byte*)executeArgumentsLengths,
			executeArgumentsData
	);

	finishResult(result);

	// The parent cannot access the storage of the child.
	int len = storageLoadLength(childKey, 32);
	if (len == 0) {
		finishResult(0);
	} else {
		finishResult(1);
	}
}

void parentFunctionChildCall_BigInts() {
	bigInt intA = bigIntNew(84); 
	bigInt intB = bigIntNew(96);
	bigInt intC = bigIntNew(1024);

	byte argumentSize = sizeof(bigInt);

	// All SmartContracts expect their integer arguments in Big Endian form, so
	// we need to reverse them (we're in Little Endian here) in order to pass
	// them to the childSC.
	bigInt arguments[] = {
		reverseU32(intA),
		reverseU32(intB),
		reverseU32(intC)
	};
	int argumentLengths[3] = {argumentSize, argumentSize, argumentSize};

	byte childAddress[] = "childSC.........................";
	byte functionName[] = "childFunction_BigInts";
	int result = executeOnDestContext(
			200000,
			childAddress,
			executeValue,
			functionName,
			21,
			3,
			(byte*)argumentLengths,
			(byte*)arguments
	);
	finishResult(result);

	// The parent cannot access the big integer created by the child.
	result = 0;
	long long x = bigIntGetInt64(4);
	if (x != 0) {
		result = 1;
		int64finish(x);
	}
	finishResult(result);
}

void parentFunctionChildCall_OutOfGas() {
	storageStore(parentKeyA, 32, parentDataA, 11);
	bigIntSetInt64(12, 42);
	finish(parentFinishA, 13);

	byte childAddress[] = "childSC.........................";
	byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};
	byte functionName[] = "childFunction_OutOfGas";
	int result = executeOnDestContext(
			3500,
			childAddress,
			executeValue,
			functionName,
			22,
			0,
			0,
			0
	);

	storageStore(parentKeyB, 32, parentDataB, 11);
	finishResult(result);
}

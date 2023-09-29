#include "../kalyan3104/context.h"
#include "../kalyan3104/bigInt.h"
#include "../kalyan3104/test_utils.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte parentFinishA[] = "parentFinishA";
byte parentFinishB[] = "parentFinishB";

byte childAddress[] = "childSC.........................";
byte vaultAddress[] = "vaultAddress....................";
byte thirdPartyAddress[] = "thirdPartyAddress...............";

byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void handleBehaviorArgument();
void handleTransferToVault();
int mustTransferToVault();
int isVault();

void parentPerformAsyncCall() {
	storageStore(parentKeyA, 32, parentDataA, 11);
	storageStore(parentKeyB, 32, parentDataB, 11);
	finish(parentFinishA, 13);
	finish(parentFinishB, 13);

	value[31] = 3;
	byte transferData[] = "hello";
	transferValue(thirdPartyAddress, value, transferData, 5);
	
	byte callData[] = "transferToThirdParty@03@207468657265@00";
	callData[38] = int64getArgument(0) + '0';

	value[31] = 7;
	asyncCall(childAddress, value, callData, 39);
}

void callBack() {
	int numArgs = getNumArguments();
	if (numArgs < 2) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

	byte loadedData[11];
	storageLoad(parentKeyB, 32, loadedData);

	int status = 0;
	for (int i = 0; i < 11; i++) {
		if (loadedData[i] != parentDataB[i]) {
			status = 1;
			break;
		}
	}

	handleBehaviorArgument();
	handleTransferToVault();

	finishResult(status);
}

void handleTransferToVault() {
	if (mustTransferToVault()) {
		value[31] = 4;

		transferValue(vaultAddress, value, 0, 0);
	}
}

int mustTransferToVault() {
	int numArgs = getNumArguments();
	byte childArgument[10];

	if (numArgs == 3) {
		getArgument(2, childArgument);
		if (isVault(childArgument)) {
			return 0;
		}
	}

	if (numArgs == 4) {
		getArgument(3, childArgument);
		if (isVault(childArgument)) {
			return 0;
		}
	}

	return 1;
}

int isVault(byte *string) {
	byte vault[] = "vault";
	for (int i = 0; i < 5; i++) {
		if (vault[i] != string[i]) {
			return 0;
		}
	}

	return 1;
}

void handleBehaviorArgument() {
	int numArgs = getNumArguments();
	if (numArgs < 4) {
		return;
	}

	byte behavior = int64getArgument(1);

	if (behavior == 3) {
		byte msg[] = "callBack error";
		signalError(msg, 14);
	}
	if (behavior == 4) {
		byte msg[] = "loop";
		while (1) {
			finish(msg, 4);
		}
	}

	finish(&behavior, 1);
}


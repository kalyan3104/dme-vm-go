#include "../kalyan3104/context.h"

const int dataLen = 10000;
byte data[dataLen] = {};

void childFunction() {
	byte msg[] = "child";
	finish(msg, 5);

	for (int i = 0; i < dataLen; i++) {
		data[i] = i;
	}

	for (int i = 1; i <= dataLen; i++) {
		int64finish(data[i-1]);
	}
}

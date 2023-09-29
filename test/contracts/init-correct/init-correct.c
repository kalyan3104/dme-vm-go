#include "../kalyan3104/context.h"

void init() {
	int arg = int64getArgument(0);
	if (arg == 0) {
		unsigned char msg[] = "init successful";
		finish(msg, 15);
	}
	if (arg == 1) {
		byte msg[] = "don't do this";
		signalError(msg, 13);
	}
	if (arg == 2) {
		byte msg[] = "loop";
		while (1) {
			finish(msg, 4);
		}
	}
}

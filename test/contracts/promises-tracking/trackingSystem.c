#include "../kalyan3104/context.h"

byte isTrainBooked[] = "trainBooked";

void init() {
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 0);
}

void bookTrain() {
    int64storageStore(isTrainBooked, sizeof(isTrainBooked), 1);
}

package main

import "fmt"
import "unsafe"

/*
#cgo LDFLAGS: -laax
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

#include <aax/aax.h>


aaxConfig config=NULL;
aaxBuffer *buffers=NULL;
aaxEmitter *emitters=NULL;
int buffersLen=0;
int emittersLen=0;
int buffersCap=0;
int emittersCap=0;

void growArray(void *** array, int len, int * cap) {
if (len>(*cap)/4*3) {
(*cap)*=2;
*array=realloc(*array, *cap);
}
}

void init() {
config = aaxDriverOpenDefault(AAX_MODE_WRITE_HRTF);
}

void terminate() {
if (config!=NULL)
aaxDriverDestroy(config);
}

void createArray(void *** array, int * size) {
if (*array==NULL) {
*array = calloc(sizeof(void*), 100);
*size=100;
}
}

int arrayAdd(void *** array, int * len, int * cap, void * data) {
growArray(array, *len, cap);
*array[*len]=data;
return (*len)++;
}

int checkArrayElem(void ** array, int indx) {
if (array[indx]!=NULL)
return 1;
return 0;
}

int loadBuffer(char * file) {
createArray(&buffers, &buffersCap);
return arrayAdd(&buffers, &buffersLen, &buffersCap, aaxBufferReadFromStream(config, file));
}

int createEmitter() {
createArray(&emitters, &emittersCap);
return arrayAdd(&emitters, &emittersLen, &emittersCap, aaxEmitterCreate());
}

void setEmitterBuffer(int emitter, int buffer) {
if (!(checkArrayElem(buffers, buffer) ||checkArrayElem(emitters, emitter)))
return;
aaxEmitterAddBuffer(emitters[emitter], buffers[buffer]);
aaxEmitterSetState(emitters[emitter], AAX_PLAYING);
}
*/
import "C"

func LoadSound(file string) int {
	var ctext *C.char = C.CString(file)
	defer C.free(unsafe.Pointer(ctext))
	b := C.loadBuffer(ctext)
	e := C.createEmitter()
	C.setEmitterBuffer(e, b)
	return int(e)
}

func main() {
	C.init()
	i := LoadSound("/home/rkruger/audio/out.wav")
	fmt.Printf("Loaded sound %d\n", i)
	fmt.Scanf("\n")
	C.terminate()
}

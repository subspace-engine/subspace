package snd

// We have to use SDL2 Mixer through CGO, as veandco/sdl2 doesn't allow us to get to the chunk data

import "sync"
import "unsafe"

/*
#cgo LDFLAGS: -lSDL2 -lSDL2_mixer -lopenal
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

#include <SDL2/SDL_mixer.h>
#include <AL/al.h>
#include <AL/alc.h>

typedef struct {
int bufnum;
Mix_Chunk *chunk;
char * file;
} Buffer;

Buffer ** buffers=NULL;
int buffersLen=0;
int * sources=NULL;
int sourcesLen=0;
int buffersCap=10;
int sourcesCap=10;

ALCdevice *device;
ALCcontext *context;

void growArray(void *** array, int len, int * cap) {
if (len>(*cap)/4*3) {
(*cap)*=2;
*array=realloc(*array, *cap);
}
}

void setListenerPosition(float x, float y, float z) {
alListener3f(AL_POSITION, x, y, z);
}

void setListenerOrientation(float atx, float aty, float atz, float upx, float upy, float upz) {
float vec[6] = {atx, aty, atz, upx, upy, upz};
alListenerfv(AL_ORIENTATION, vec);
}

void init() {
createArray(&buffers, &buffersCap);
createArray(&sources, &sourcesCap);
ALCenum error;
Mix_Init(0);
Mix_OpenAudio(44100, 0x8010, 1, 8192);
int * attrlist[3];
attrlist[0]=ALC_FREQUENCY;
attrlist[1]=48000;
attrlist[2]=0;
alGetError();
device = alcOpenDevice(NULL);
error = alGetError();
if (error!=AL_NO_ERROR) {

}

   context = alcCreateContext(device, NULL);
if (error!=AL_NO_ERROR) {
}
if (!alcMakeContextCurrent(context)) {
error = alGetError();
if (error!=AL_NO_ERROR) {
}
}

setListenerPosition(0.0f, 0.0f, 0.0f);
setListenerOrientation(0.0f, 0.0f, -1.0f, 0.0f, 1.0f, 0.0f);
}

void terminate() {
for (int i = 0; i < sourcesLen; i++)
if (sources[i]!=NULL) {
alDeleteSources(1, &sources[i]);
}
for (int i = 0; i < buffersLen; i++)
if (buffers[i]!=NULL) {
alDeleteBuffers(1, &buffers[i]->bufnum);
Mix_FreeChunk(buffers[i]->chunk);
free(buffers[i]);
}
alcDestroyContext(context);
alcCloseDevice(device);
Mix_CloseAudio();
Mix_Quit();
}

void createArray(void *** array, int * size) {
if (*array==NULL) {
*array = calloc(sizeof(void*), 100);
*size=100;
}
}

int arrayAdd(void *** array, int * len, int * cap, void * data) {
growArray(array, *len, cap);
(*array)[*len]=data;
*len=*len+1;
return (*len)-1;
}

int checkArrayElem(void ** array, int indx) {
if (array[indx]!=NULL)
return 1;
return 0;
}

int findBuffer(char * file) {
for (int i =0; i < buffersLen; i++)
if (strcmp(file, buffers[i]->file) ==0)
return i;
return -1;
}

int loadBuffer(char * file) {
int ret;
if ((ret=findBuffer(file))>=0) {
return buffers[ret]->bufnum;
}
ALCenum error;
Buffer *buffer = malloc(sizeof(Buffer));
buffer->chunk=Mix_LoadWAV(file);
buffer->file=malloc(strlen(file)+1);
strcpy(buffer->file, file);
alGetError();
alGenBuffers(1, &buffer->bufnum);
error = alGetError();
if (error!=AL_NO_ERROR)
printf("Error generating buffer");
ret= arrayAdd(&buffers, &buffersLen, &buffersCap, (void*)buffer);
alBufferData(buffers[ret]->bufnum, AL_FORMAT_MONO16, buffers[ret]->chunk->abuf, buffers[ret]->chunk->alen, 44100);
return buffers[ret]->bufnum;
}

int createSource() {
int source = 0;
ALCenum error;
alGetError();
alGenSources(1, &source);
error = alGetError();
if (error!=AL_NO_ERROR)
printf("Error generating sources");
arrayAdd(&sources, &sourcesLen, &sourcesCap, source);
return source;
}

void setSourceBuffer(int source, int buffer) {
ALCenum error;
alGetError();
alSourcei(source, AL_BUFFER, buffer);
error=alGetError();
if (error!=AL_NO_ERROR)
printf("Error setting buffer\n");
alSourcef(source, AL_PITCH, 1.0f);
   alSourcef(source, AL_GAIN, 1.0f);
}

void playSource(int source) {
alSourcePlay(source);
}

void setPosition(int source, float x, float y, float z) {
ALCenum error;
alGetError();
alSource3f(source, AL_POSITION, x, y, z);
error = alGetError();
if (error!=AL_NO_ERROR)
printf("Error setting position");
}

int isPlaying(int source) {
int ret=0;
alGetSourcei(source, AL_SOURCE_STATE, &ret);
return ret==AL_PLAYING;
}

void setLooping(int sound, int looping) {
alSourcei(sound, AL_LOOPING, looping);
}
*/
import "C"

var lock sync.Mutex

type ALDriver struct {
}

func (self ALDriver) LoadSound(file string) int {
	lock.Lock()
	defer lock.Unlock()
	var ctext *C.char = C.CString(file)
	defer C.free(unsafe.Pointer(ctext))
	b := C.loadBuffer(ctext)
	s := C.createSource()
	C.setSourceBuffer(s, b)
	return int(s)
}

func (self ALDriver) SetPosition(source int, x float64, y float64, z float64) {
	C.setPosition(C.int(source), C.float(x), C.float(y), C.float(z))
}

func (self ALDriver) PlaySource(source int) {
	C.playSource(C.int(source))
}

func (self ALDriver) SetListenerPosition(x float64, y float64, z float64) {
	C.setListenerPosition(C.float(x), C.float(y), C.float(z))
}

func (self ALDriver) SetListenerOrientation(xAt float64, yAt float64, zAt float64, xUp float64, yUp float64, zUp float64) {
	C.setListenerOrientation(C.float(xAt), C.float(yAt), C.float(zAt), C.float(xUp), C.float(yUp), C.float(zUp))
}

func (self ALDriver) Destroy() {
	lock.Lock()
	defer lock.Unlock()
	C.terminate()
}

func MakeALDriver() ALDriver {
	lock.Lock()
	defer lock.Unlock()
	C.init()
	return ALDriver{}
}

func (self ALDriver) IsPlaying(source int) bool {
	if C.isPlaying(C.int(source)) > 0 {
		return true
	}
	return false
}

func (self ALDriver) SetLooping(sound int, looping bool) {
	var iloop int
	if looping {
		iloop = 1
	} else {
		iloop = 0
	}
	C.setLooping(C.int(sound), C.int(iloop))
}

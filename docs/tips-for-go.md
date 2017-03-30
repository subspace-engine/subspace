It is often easier to send a value over a channel rather than
synchronising access to it with a mutex. The latter will be faster if
done right though, but the first is more debugable.

---
Keep interfaces small. Rather nest / combine interfaces if needed.

### Example
```go
type Reader interface {
  Read()
}

type Writer interface {
  Write()
}

type ReadWriter interface {
  Reader
  Writer
}
```

---

When making types (classes), think in terms of interfaces, but try to
define only those that you really need.

---

Accept interfaces as parameters and return types (unless you don't
need interfaces).

---

`err` is almost always used as the lone name for errors.

---

Access to package internals is protected by making struct fields
lower-case and manipulating them through functions. The same holds for
types (lower-case named types can be returned, but not constructed
directly).

---

Everything in go is passed by value (func parameters and even the func
instance it is defined on). The only exceptions are slices and maps.

To that end, when wondering if your function should be defined on the
obj or *obj, go for *obj unless you have a good reason not to.

### Example
```go
func (o obj) edit(int data) {
  ...
}
```
Will modify the passed copy of the obj and will be gone once the func
returns, while
```go
func (o *obj) edit(int data) {
  ...
}
```
Will modify the obj being pointed to.

---

Be careful when passing slices as parameters. The underlying array
might be modified, but other copies of the slice will not have their
length and capasity updated. Something like apend makes a new slice.
Same goes for maps.

---

The range keyword is nice to iterate over maps and slices. Just
remember that if you want to modify a value inside the slice then you
have to modify it by index in the slice.

### Wrong example
```go
for i, val := range list {
  val=1 // will not reflect in the slice
  list[i] = 1 // will
}
```

---

### A select example
```go
func process(d []data, done chan struct{}) {
  pipe := make(chan data)
  // Process
  go func() { // we could start multiple workers here
  for {
    select {
      case <- done:
        return
      case item := <- pipe:
        work(item)
    }
  }
  // ... do stuff and put data into pipe.
  close(done) // a closed channel always returns a value if used as receiver
} ()
```



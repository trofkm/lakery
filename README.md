# lakery - very simple tag-based validator





## ROADMAP:

1. Simple tag validation (e.g. credential, email)
2. Validation Expressions (e.g. min=0,max=255)
3. Json-based initialization (e.g. size={min:0,max:255})


should support the following syntax:

This one I personally very like, because of 'each' keyword
```go
type Simple struct{
	vals []string `lakery:"each={min=0,max=23,credential}"`
}
```

I think it is good to use `min`, `max` tags for numbers, strings, arrays


also seems like a good approach to use a dive keyword:

```go

type Simple struct {
	vals []Value `lakery:"dive"`
}


I have following thoughts in mind:

1) should I recursively iterate fields using reflection or use `dive` keyword?
2) how can I support `lakery:"each={min=0,max=23,credential}"` format
3) I think min,max,dive,required, each should be builtin

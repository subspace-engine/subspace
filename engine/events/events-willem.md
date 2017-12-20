Subspace events

At the most basic level the purpose of events is to update state in the client or server. Clients need to know about certain changes on the server and the server needs to know when the client acts on the world.

requirements:
* must be able to convey small state changes as well as large portions of game state (performance must be scalable)
Allow for an arbitrary? / very large number of named peramiters to be transmitted in the same event. Motivation: avoid race conditions.
Support at least the following go types: string, float64 and int.

design:
I propose that we can use protobuf to send events. For both cases the user will still need to write some layer that knows what to do with events when they arive and what to send to clients. I am not sure what form the actual events should take.

I can think of two possible designs:

design 1:
static events
The user defines all valid forms of events by writing a protobuf spec for each. Both the client and server load the pb.go files and can then understand the selected events.

Advantages:
Predefined events will have better performance compared to events with more dynamic content. (design 2)
Easier to access fields of known protobuf where more dynamic events require interpretation.

disadvantages:
The user needs to understand engine internals. This can be counteracted by having pre-defined events.

design 2:
Dynamic events
All events superficially have the same structure.
message Dynamic event {
// each element in this array represents the type of the next parameter in the event.
repeated paramTypes ParamType; // param type is an enum with values representing the types messages can hold 
// possibly have a list of names for parameters?
repeated IntParams int;

//parameters should be extracted in the same order as they were inserted
repeated StringParams string;
repeated floatParams float64;
}

Three (or more) counters will need to be maintained during extraction, one for each type. The value list for each type in the paramType list will be consulted and the index for that type will be advanced.

Other similar schemes are also possible, but has similar advantages and disadvantages.

advantages:
No need for predefined events.

disadvantages:
Dynamic events will be less performant compared to predefined events with the same number of parameters.

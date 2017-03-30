# Thoughts on an Engine Design
#### Rynhardt Kruger, on the 29th day in the month of March, A.D. 2017

## World Layout

All simulation should be contained in a "world", of which several can be active at once. Each world should be able to have one or more active players associated with it, worlds with no active players should be unloaded (and stored if desired).

-- Mark: This type of "multiple worlds" thinking seems good for a MUD. How about for e.g. a strategy game, or other games in which "worlds" will be short-lived and with few players? Perhaps multiple worlds could be optional? 
-- Willem: I guess we can call each component what we like, but here is what I
think might work better after Mark and I spoke:
Define a concept called "area". It can function as a component that
restricts objects based on distance, as well as allow sections of the
game environment to be loaded (think of what you described as "world",
but also in terms of calculating what a player can see, either static
such as an entire room, or dynamically adjustable).
An area / world can be loaded and / or unloaded based on an automatic
calculation based on distance (but the recalculation rate should be
adjustable.
I also think partial loading should be optional. For some games it's
not needed and is just overhead. To my mind we need other core
features before this is added, but we should design in such a way that
it is possible.

The fixed structure of a world e.g. terrain, open areas, walls etc should be stored as an 3-dimensional array of byte values, in other words voxels. Each voxel's value should denote a specific type of object or element, i.e. stone, wood, grass, soil, sand, etc. The world should also maintain an array translating voxel values to their corresponding object representations.

-- Mark: Sounds good. I also recommend that the voxel size representation be specific to each game. (e.g. maybe 1x1x1m for Rynhardt's game, but 2x2x2m for my game or whatever)
-- Willem: This seems fine, but we should keep in mind when it might cause
problems if we want it to be a multi-purpose game engine.

Voxel arrays should have a fixed size, and a world should be able to maintain several arrays at once in "chunk" structures. Chunks should be loaded and unloaded as necessary, with only chunks with active players in memory at any point in time.

-- Mark: The idea of "chunks with active players in memor" seems focused on first-person type games. What about for strategy games? What about games with simulated NPCs that can go off and do their own thing, but where you want to keep simulating them?

Non-fixed objects, e.g. creatures, items, and the player itself, as well as complex objects like doors, should be represented in object form. An object should have one or more shapes associated with it, on which collision detection is performed. Collision detection should also be calculated between free-moving objects and the voxel structure.

-- Mark: Sounds good.
Willem: I  think we should also be able to ignore certain world objects when
doing collision detection (for performance and if detail is not
required or will get in the way of gameplay). Also, un-moved objects
can be ignored if nothing near them moved.

Objects should be able to have child-objects linked to it, which will move along with the parent object. This could be used to implement complex structures, i.e. a vehicle with several components. This could also be used to simulate the laws of physics, i.e. stepping onto a moving platform assigns the player as a child of the platform, which will move along with it.

-- Mark: Alright.

I'm not quite sure about the representation of objects themselves. At the moment I'm leaning towards hash maps, for their ease of serialization. If we follow this root, implementing a new kind of object would mearly require writing a "New*" function, i.e. NewPerson, which will populate the map with default values required for a person (possibly chaining to a base function like NewCreature).

-- Mark: Not sure I understand what you mean by representing objects with hashmaps? We should have a more in-depth discussion about this.
-- Willem: Something I've been taught: Implementation details should not be part
of feature design, but what features we decide on should obviously
fall within what the method of implementation can do. Weather we use
maps, serialized go types, protobufs or whatever is just an
implementation detail and all of them can work..

What I do know is that objects should be able to identify the kinds they represent, in superclass order, e.g. {"person", "creature", "mobile", "thing"}. This will enable officiant action processing (see the actions section below).

-- Mark: So, Go doesn't have inheritence. I assume this is kind of an 'internal' inheritance. Is it the kind of thing that can be updated 'on the fly'? (e.g. using scripting, implement a new kind).

### Printing Messages

Messages should be handled on a per-object basis, with each object implementing a notify method. Thus, calling notify on an object should print the message to all connected players within a specified radius of the object, and calling it on a world notifies everyone in the world. There should also be a global notify function/method, for things like "Server rebooting in 5 minutes!".

Calling notify on an object representing a player should format the message accordingly, so that it is 2nd person for the player itself, and 3rd person for other observers. For example, "You open the chest." vs "Jannie opens the chest.". This could be handled with tads-like format strings, e.g. "(You) open(s) the chest.".

-- Mark: This implementation seems applicable for RPG and MUD type games. Can we make it sort of optional and think of alternatives and/or see how it could be applied to strategy games?

### Actions

Events in the game should be managed by actions, which are represented as an ID, initiate, optional direct object, and optional indirect object. Later we might want to add several types of actions, e.g. direct object vs indirect object. Actions should be managed separately from the structure of the game. This will allow the implementation of actions that can act on more than one object, i.e. two things colliding should play a sound, or generate a textual description, regardless of the specific objects involved. Actions should also be prioritized, with actions registered later taking precedence over actions registered earlier. This can for instance be used to implement an action which prevents any creature (including a player) from performing any action while incapacitated.

-- Mark: Sounds good.
-- Willem: I think we need to have a whole design chapter just for game events,
but likewise also for game objects, world etc. For example, could we
tie the concept of area / world to events? Events can have a list of
targets, or maybe an area. For first / third person views two
different events can be generated by the world, one that has the
person as it's only target, the other one targeted in an area.
In my mind, events originating from the client / player as attempts at
acting on the world and events from the server are to report changes /
information to each client. I really think we need to have a
centralised server. A distributed world with state maintained by each
client would be crazy hard to keep in sync and bug-free.

The representation of game logic as actions also makes it easy to do command parsing, with single keystrokes translating to commands (which may be remapped), and commands translating to actions.

-- Mark: Cool

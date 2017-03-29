# Thoughts on an Engine Design
#### Rynhardt Kruger, on the 29th day in the month of March, A.D. 2017

This file contains Markdown directives. To view this file as HTML, process it with markdown. For instance, to view an HTML version in the console, do:

$ markdown design-rynkruger.md | elinks

## World Layout

All simulation should be contained in a "world", of which several can be active at once. Each world should be able to have one or more active players associated with it, worlds with no active players should be unloaded (and stored if desired).

The fixed structure of a world e.g. terrain, open areas, walls etc should be stored as an 3-dimensional array of byte values, in other words voxels. Each voxel's value should denote a specific type of object or element, i.e. stone, wood, grass, soil, sand, etc. The world should also maintain an array translating voxel values to their corresponding object representations.

Voxel arrays should have a fixed size, and a world should be able to maintain several arrays at once in "chunk" structures. Chunks should be loaded and unloaded as necessary, with only chunks with active players in memory at any point in time.

Non-fixed objects, e.g. creatures, items, and the player itself, as well as complex objects like doors, should be represented in object form. An object should have one or more shapes associated with it, on which collision detection is performed. Collision detection should also be calculated between free-moving objects and the voxel structure.

Objects should be able to have child-objects linked to it, which will move along with the parent object. This could be used to implement complex structures, i.e. a vehicle with several components. This could also be used to simulate the laws of physics, i.e. stepping onto a moving platform assigns the player as a child of the platform, which will move along with it.

I'm not quite sure about the representation of objects themselves. At the moment I'm leaning towards hash maps, for their ease of serialization. If we follow this root, implementing a new kind of object would mearly require writing a "New*" function, i.e. NewPerson, which will populate the map with default values required for a person (possibly chaining to a base function like NewCreature).

What I do know is that objects should be able to identify the kinds they represent, in superclass order, e.g. {"person", "creature", "mobile", "thing"}. This will enable officiant action processing (see the actions section below).

### Printing Messages

Messages should be handled on a per-object bases, with each object implementing a notify method. Thus, calling notify on an object should print the message to all connected players within a specified radius of the object, and calling it on a world notifies everyone in the world. There should also be a global notify function/method, for things like "Server rebooting in 5 minutes!".

Calling notify on an object representing a player should format the message accordingly, so that it is 2nd person for the player itself, and 3rd person for other observers. For example, "You open the chest." vs "Jannie opens the chest.". This could be handled with tads-like format strings, e.g. "(You) open(s) the chest.".

### Actions

Events in the game should be managed by actions, which are represented as an ID, initiate, optional direct object, and optional indirect object. Later we might want to add several types of actions, e.g. direct object vs indirect object. Actions should be managed separately from the structure of the game. This will allow the implementation of actions that can act on more than one object, i.e. two things colliding should play a sound, or generate a textual description, regardless of the specific objects involved. Actions should also be prioritized, with actions registered later taking precedence over actions registered earlier. This can for instance be used to implement an action which prevents any creature (including a player) from performing any action while incapacitated.

The representation of game logic as actions also makes it easy to do command parsing, with single keystrokes translating to commands (which may be remapped), and commands translating to actions.

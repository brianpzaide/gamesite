# Gamesite

A Go-based game server that hosts a collection of simple turn-based web games. It is designed with flexibility in mind, allowing for easy registration and addition of new games.

### Games Implemented (so far)

* #### Three Tic Tac Toe
  <img src="tttt.png" width="300" height="125">

* #### Nested Tic Tac Toe
  <img src="nttt.png" width="300" height="301">

* #### Reversi
  <img src="reversi.png" width="300" height="300">

* #### Pawns only Chess
  <img src="poc.png" width="300" height="300">

* #### [Maxit](https://play.google.com/store/apps/details?id=com.loonybot.maxitmonkey&gl=US)


### Idea

Players can select a game, create a room, and share the room's URL with a friend to play. Rooms are automatically closed under the following conditions:
- The game reaches an end state.
- The other player fails to join within the time limit.
- A player fails to make a move within the time limit (adjustable).


#### To build


```bash
make build
```

The above command generates a single binary using [stuffbin](https://github.com/knadh/stuffbin), bundling all HTML files and images.

#### To run

```bash
make run
```

The application listens on port ```8080```.

#### Architecture
The application consists of 3 components

* Web: Handles REST and WebSocket handlers.
* Room: An in-memory data structure that manages game states, including making moves, tracking players, and determining game outcomes.
* Hub: Responsible for managing rooms (creation and destruction).

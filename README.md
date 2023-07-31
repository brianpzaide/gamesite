# Gamesite

A collection of all simple games(more to add later), that I find interesting. All the games are turn based.

The following games are implemented

* #### Three Tic Tac Toe
  <img src="tttt.png" width="300" height="125">
* #### Nested Tic Tac Toe
  <img src="nttt.png" width="300" height="301">
* #### Reversi
  <img src="reversi.png" width="300" height="300">
* #### [Maxit](https://play.google.com/store/apps/details?id=com.loonybot.maxitmonkey&gl=US)


### Idea
A person selects a game, then creates a room and then shares the room's url with a friend they would like to play with.

The rooms are automatically destroyed when:

* the game reaches an end state.
* the other player fails to join the room with in the time limit.
* the player fails to make a move with in the time limit(this limit can be adjusted) on his/her turn.



#### To build

`make build`

Thanks to [stuffbin](https://github.com/knadh/stuffbin) the above command produces just one single binary, with all the html files and images stuffed into it.

#### To run

`make run`

The application listens on port 8080.

#### Architecture
This application consists of 3 components

* Web (REST and web socket handlers)
* Room (room is an in-memory data structure that stores the game and players. It manages the game's state, such as making the actual move, keeping track of which player to play, whether the game is draw, inprogress, who has won or lost the game)
* Hub (hub is responsible for managing the rooms (creating and destroying))






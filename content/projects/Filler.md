---
title: Filler
date: 2/4/2023
tags:
  - Web Dev
  - Rust
---
## Motivation 
Filler is one of my favorite games, the simplicity but also the small chance for very deep complexity. Whenever I looked online to play filler, I couldn't find any websites that provided a the game. So I decided to make one.
### How it works
At the start of the game, each player has a square in the bottom left or top right corner. Becoming a color of a square means that the player captures that square. The game ends when all the squares have been captured and the winner is the player with the most squares.
## Development
To build the website, I used [React JS](https://reactjs.org/), and to build the computer the player can play against, I used [webassembly](https://webassembly.org/). I wanted to use webassembly because I wanted the algorithm to be fast with the minimal cost of computation. I used Rust to build the algorithm, properly named *Phil*. 
### Algorithm
The algorithm is goes by a [brute force](https://en.wikipedia.org/wiki/Brute-force_search) method meaning it tries every color that it is currently touching at a depth of 8. 
This means that it looks 8 moves ahead and picks the move that would get it the most territory after 8 moves. This also includes capturing squares by blocking the opponent off from them.
This is checked by doing a [breath-first search](https://en.wikipedia.org/wiki/Breadth-first_search) whenever a square is captured that may be blocked off by the opponent. For instance, lets say a square is captured on the edge of the board and the squares adjacent to it are not already captured by the player. The algorithm then does a breadth-first search to see if it can path to an opponents square and if not, the number of squares that can't path to an opponents square are counted in the score. The blocked off squares are not shown in the actual score, rather, taken into account in the algorithm when choosing the next move.
### Algorithmic Flaws
There are many flaws in the algorithm. A better way to go about this would probably be to use a machine learning approach, but a straight forward simple algorithm was mainly what I was looking for. 
Brute force is not the best way to go about this, but either it was either brute force or a machine learning model. A dynamic programming approach would have not worked all that well because in order to implment it, the algorithm would have been more accurate, but the time complexity would have been very exponential. You see, with every move, the board changes, so to implement a dynaimc programming approach, you would have to calculate everyone of your move combinations and with everyone of those combinations, you would have to calculate everyone of your move combinations based on the opponents move.
I wanted the algorithm to be fast so I did not implment this dynamic programming approach even though it is interesting. In the future I am thinking of implementing a machine learning model which would be intersting because the colors on the board are generated randomly causing one player to almost always have an advantage over another. 
Finding this advantage is what makes filler fun.
## Conclusion
This was a very fun project to work on as I can finally play filler on a website! Creating Phil was also very fun, although there are still many flaws in the algorithm, it was intersting to use webassembly and I think I will be using webassembly in future projects. Overall a this was a very interesting project and I hope you can enjoy filler as much as I do.